package script

import (
	"fmt"
	"net/url"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
	lua "github.com/yuin/gopher-lua"
)

type WebsocketModule struct {
	conn *websocket.Conn
	done chan struct{} // 通知协程停止的通道

	q   []queue
	qmu sync.Mutex

	repolst []Report
}

type queue struct {
	buff []byte
}

var (
	ErrNil = lua.LString("nil")
)

func NewWebsocketModule() *WebsocketModule {
	return &WebsocketModule{
		done: make(chan struct{}),
	}
}

func (ws *WebsocketModule) Loader(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"dail":  ws.dail,
		"close": ws.Close,

		"read":  ws.readmsg,
		"write": ws.writemsg,
	})

	L.Push(mod)

	return 1
}

func (ws *WebsocketModule) dail(L *lua.LState) int {

	err := ws._dail(L.ToString(1), L.ToString(2), L.ToString(3))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(ErrNil)

	return 1
}

func (ws *WebsocketModule) _read() {
	if ws.conn == nil {
		return
	}

	for {
		select {
		case <-ws.done:
			fmt.Println("chan exit")
			return
		default:

			defer func() {
				if r := recover(); r != nil {
					var buf [4096]byte
					n := runtime.Stack(buf[:], false)
					fmt.Println("panic:", string(buf[:n]))
				}
			}()

			_, msg, err := ws.conn.ReadMessage()
			if err != nil {
				fmt.Println("read msg err", err.Error())
				return
			}

			ws.qmu.Lock()
			ws.q = append(ws.q, queue{buff: msg})
			ws.qmu.Unlock()
		}
	}
}

func (ws *WebsocketModule) _dail(scheme string, host string, port string) error {

	u := url.URL{Scheme: scheme, Host: host + ":" + port, Path: "/ws"}

	fmt.Println("dail", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial:%s", err.Error())
	}

	ws.conn = c
	go ws._read()

	return nil
}

func (ws *WebsocketModule) _close() {
	if ws.conn != nil {
		ws.conn.Close()
		ws.conn = nil

		close(ws.done)
	}
}

func (ws *WebsocketModule) Close(L *lua.LState) int {
	ws._close()
	L.Push(ErrNil)
	return 1
}

func (ws *WebsocketModule) readmsg(L *lua.LState) int {
	if ws.conn == nil {
		L.Push(lua.LString(""))
		L.Push(lua.LString("not connected"))
		return 2
	}

	ws.qmu.Lock()
	if len(ws.q) == 0 {
		ws.qmu.Unlock()
		L.Push(lua.LString(""))
		L.Push(lua.LString("empty"))
		return 2
	}

	L.Push(lua.LString(ws.q[0].buff))
	L.Push(ErrNil)
	ws.q = ws.q[1:]
	ws.qmu.Unlock()

	return 2
}

func (ws *WebsocketModule) writemsg(L *lua.LState) int {

	if ws.conn == nil {
		L.Push(lua.LString("not connected"))
		return 1
	}

	buf := L.ToString(1)
	err := ws.conn.WriteMessage(websocket.BinaryMessage, []byte(buf))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(ErrNil)
	return 1
}

func (ws *WebsocketModule) GetReport() []Report {

	rep := []Report{}
	rep = append(rep, ws.repolst...)

	ws.repolst = ws.repolst[:0]

	return rep
}
