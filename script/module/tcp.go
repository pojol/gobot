package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"sync"
	"syscall"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type headerPair struct {
	Desc string
	Len  int
}

type customHeader struct {
	header []headerPair
}

type TCPModule struct {
	conn *net.TCPConn
	fd   int

	buf   byteQueue
	bufMu sync.Mutex

	writeTime time.Time
	repolst   []Report

	done chan struct{} // 通知协程停止的通道
}

type byteQueue []byte

// 入队
func (q *byteQueue) push(b []byte) {
	*q = append(*q, b...)
}

// 出队
func (q *byteQueue) pop(maxLen int) ([]byte, bool) {
	if len(*q) == 0 {
		return nil, false
	}

	if maxLen > len(*q) {
		maxLen = len(*q)
	}

	data := (*q)[:maxLen]
	*q = (*q)[maxLen:]
	return data, true
}

// 队列当前长度
func (q *byteQueue) haveFull() bool {
	if len(*q) < 2 {
		return false
	}

	var header [2]byte
	copy(header[:], (*q)[:2])

	var msgleni int16
	binary.Read(bytes.NewBuffer(header[:]), binary.LittleEndian, &msgleni)
	if len(*q) < int(msgleni) {
		return false
	}

	return true
}

func NewTCPModule() *TCPModule {

	tcpm := &TCPModule{}

	return tcpm
}

func (t *TCPModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"dail":  t.dail,
		"close": t.Close,

		"write": t.write,
		"read":  t.read,

		"read_msg":  t.read_msg,
		"write_msg": t.write_msg,
	})
	L.Push(mod)
	return 1
}

func (t *TCPModule) dail(L *lua.LState) int {
	err := t._dail(L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func (t *TCPModule) _dail(host string, port string) error {
	tcpServer, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return fmt.Errorf("resolve tcp addr err:%s", err.Error())
	}

	t.conn, err = net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		return fmt.Errorf("dial tcp err:%s", err.Error())
	}

	f, _ := t.conn.File()
	t.fd = int(f.Fd())

	go t._read()
	return nil
}

func (t *TCPModule) Close(L *lua.LState) int {

	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}

	close(t.done)

	L.Push(lua.LString("succ"))
	return 1
}

func (t *TCPModule) write(L *lua.LState) int {
	if t.conn == nil {
		L.Push(lua.LString("not connected"))
		return 1
	}

	msg := L.ToString(1)
	_, err := t.conn.Write([]byte(msg))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func (t *TCPModule) _read() {

	if t.conn == nil {
		return
	}

	for {
		select {
		case <-t.done:
			fmt.Println("chan exit")
			return
		default:
			buf := make([]byte, 1024)
			n, err := t.conn.Read(buf)
			if err != nil {
				fmt.Printf("syscall.read fd %v size %v err %v\n", t.fd, n, err.Error())
			}

			if n != 0 {
				t.bufMu.Lock()
				t.buf.push(buf[:n])
				t.bufMu.Unlock()
			}
		}
	}

}

func readret(L *lua.LState, ty, custom, id int, body []byte, err string) int {

	L.Push(lua.LNumber(ty))
	L.Push(lua.LNumber(custom))
	L.Push(lua.LNumber(id))
	L.Push(lua.LString(body))
	L.Push(lua.LString(err))

	return 5
}

func (t *TCPModule) read_msg(L *lua.LState) int {

	msglen := L.ToInt(1)
	msgty := L.ToInt(2)
	msgcustom := L.ToInt(3)
	msgid := L.ToInt(4)

	msgleni := int16(0)
	msgtyi := int8(0)
	msgcustomi := int16(0)
	msgidi := int16(0)

	var msgbody []byte

	t.bufMu.Lock()
	if !t.buf.haveFull() {
		t.bufMu.Unlock()
		return readret(L, 0, 0, 0, []byte{}, "nodata")
	}

	msglenb, _ := t.buf.pop(msglen)
	binary.Read(bytes.NewBuffer(msglenb), binary.LittleEndian, &msgleni)

	msgtyb, _ := t.buf.pop(msgty)
	binary.Read(bytes.NewBuffer(msgtyb), binary.LittleEndian, &msgtyi)

	msgcustomb, _ := t.buf.pop(msgcustom)
	binary.Read(bytes.NewBuffer(msgcustomb), binary.LittleEndian, &msgcustomi)

	msgidb, _ := t.buf.pop(msgid)
	binary.Read(bytes.NewBuffer(msgidb), binary.LittleEndian, &msgidi)

	msgbody, _ = t.buf.pop(int(msgleni) - (msgty + msgcustom + msgid))

	t.bufMu.Unlock()

	info := Report{
		Api:     strconv.Itoa(int(msgidi)),
		ResBody: int(msgleni),
		Consume: int(time.Since(t.writeTime).Milliseconds()),
	}
	t.repolst = append(t.repolst, info)

	return readret(L, int(msgtyi), int(msgcustomi), int(msgidi), msgbody, "")
}

func (t *TCPModule) write_msg(L *lua.LState) int {

	msglen := L.ToInt(1)
	msgty := L.ToInt(2)
	msgcustom := L.ToInt(3)
	msgid := L.ToInt(4)
	msgbody := L.ToString(5)

	if t.conn == nil {
		L.Push(lua.LString("not connected"))
		return 1
	}

	buf := bytes.NewBuffer(make([]byte, 0, msglen))

	binary.Write(buf, binary.LittleEndian, uint16(msglen))
	binary.Write(buf, binary.LittleEndian, uint8(msgty))
	binary.Write(buf, binary.LittleEndian, uint16(msgcustom))
	binary.Write(buf, binary.LittleEndian, uint16(msgid))
	buf.WriteString(msgbody)

	t.writeTime = time.Now()

	_, err := t.conn.Write(buf.Bytes())
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func (t *TCPModule) read(L *lua.LState) int {

	if t.conn == nil {
		L.Push(lua.LString("fail"))
		L.Push(lua.LString("not connected"))
		return 2
	}

	buf := make([]byte, 1024)
	n, err := t.conn.Read(buf)
	// 处理读取结果
	if err == syscall.EWOULDBLOCK {
		L.Push(lua.LString("fail"))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if n == 0 {
		L.Push(lua.LString("fail"))
		L.Push(lua.LString("nodata"))
		return 2
	}

	content := string(buf[:n])
	L.Push(lua.LString("succ"))
	L.Push(lua.LString(content))

	// 立即返回,不阻塞
	return 2
}

func (t *TCPModule) GetReport() []Report {

	rep := []Report{}
	rep = append(rep, t.repolst...)

	t.repolst = t.repolst[:0]

	return rep
}
