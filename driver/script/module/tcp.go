package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type TCPModule struct {
	conn *net.TCPConn
	fd   int
	br   binary.ByteOrder

	buf   byteQueue
	bufMu sync.Mutex

	repolst []Report

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
func (q *byteQueue) haveFull(br binary.ByteOrder) bool {
	if len(*q) < 2 {
		return false
	}

	var header [2]byte
	copy(header[:], (*q)[:2])

	var msgleni int16
	binary.Read(bytes.NewBuffer(header[:]), br, &msgleni)

	return len(*q) >= int(msgleni)
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

		"report": t.report,
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

	// 需要解析 msglen，所以需要传个 byteorder 类型进来
	bs := L.GetGlobal("ByteOrder").String()
	if bs != Little && bs != Big {
		t.br = binary.LittleEndian
		fmt.Println("byteSort is not valid, use default LittleEndian")
	} else {
		if bs == Little {
			t.br = binary.LittleEndian
		} else {
			t.br = binary.BigEndian
		}
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

func (t *TCPModule) _close() {
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil

		close(t.done)
	}
}

func (t *TCPModule) Close(L *lua.LState) int {
	t._close()
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

func (t *TCPModule) read(L *lua.LState) int {

	if t.conn == nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(""))
		L.Push(lua.LString("not connected"))
		return 3
	}

	t.bufMu.Lock()
	if !t.buf.haveFull(t.br) {
		t.bufMu.Unlock()
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(""))
		L.Push(lua.LString("nodata"))
		return 3
	}

	msglen := L.ToInt(1)
	msgleni := int16(0) // msg length

	msglenb, _ := t.buf.pop(msglen)
	binary.Read(bytes.NewBuffer(msglenb), t.br, &msgleni)

	msgbody, _ := t.buf.pop(int(msgleni) - 2)

	t.bufMu.Unlock()

	L.Push(lua.LNumber(msgleni))
	L.Push(lua.LString(msgbody))
	L.Push(lua.LString(ErrNil))

	// 立即返回,不阻塞
	return 3
}

func (t *TCPModule) report(L *lua.LState) int {
	id := L.ToString(1)
	errmsg := L.ToString(2)
	t.repolst = append(t.repolst, Report{id, errmsg})
	return 0
}

func (t *TCPModule) GetReport() []Report {

	rep := []Report{}
	rep = append(rep, t.repolst...)

	t.repolst = t.repolst[:0]

	return rep
}
