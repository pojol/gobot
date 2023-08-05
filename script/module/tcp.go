package script

import (
	"fmt"
	"net"
	"syscall"

	lua "github.com/yuin/gopher-lua"
)

type TCPModule struct {
	conn *net.TCPConn
	fd   int
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
