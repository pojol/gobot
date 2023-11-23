package script

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"syscall"

	lua "github.com/yuin/gopher-lua"
)

type TCPModule struct {
	conn *net.TCPConn
	fd   int
	buf  byteQueue
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
func (q *byteQueue) Len() int {
	return len(*q)
}

func NewTCPModule() *TCPModule {

	tcpm := &TCPModule{}

	return tcpm
}

func (t *TCPModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"dail":  t.dail,
		"close": t.close,

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
	syscall.SetNonblock(t.fd, true)

	return nil
}

func (t *TCPModule) close(L *lua.LState) int {

	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}

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

func (t *TCPModule) _read() error {

	if t.conn == nil {
		return errors.New("not connected")
	}

	buf := make([]byte, 1024)

	//n, err := t.conn.Read(buf)
	// 非阻塞读取
	n, err := syscall.Read(t.fd, buf)
	if err != nil {
		return fmt.Errorf("syscall.read fd %v size %v err %v", t.fd, n, err.Error())
	}

	if n != 0 {
		//dest := make([]byte, n)
		//copy(dest, buf)
		//fmt.Println(t.fd, "buf push", n, dest)

		t.buf.push(buf[:n])
	} else {
		fmt.Println("continue")
	}

	return nil
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

	err := t._read()
	if err != nil {
		return readret(L, 0, 0, 0, []byte{}, err.Error())
	}

	if t.buf.Len() < msglen+msgty+msgcustom+msgid {
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

	buf := make([]byte, 128) //test
	// 非阻塞读取
	n, err := syscall.Read(t.fd, buf)
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
