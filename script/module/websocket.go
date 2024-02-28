package script

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	lua "github.com/yuin/gopher-lua"
)

type WebsocketModule struct {
	conn *websocket.Conn
	done chan struct{} // 通知协程停止的通道

	q   []queue
	qmu sync.Mutex
}

type queue struct {
	buff []byte
}

func (ws *WebsocketModule) Loader(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"dail":  ws.dail,
		"close": ws.Close,

		"read_msg":  ws.read_msg,
		"write_msg": ws.write_msg,
	})

	L.Push(mod)

	return 1
}

func (ws *WebsocketModule) dail(L *lua.LState) int {

	err := ws._dail(L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LString("succ"))

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

			_, msg, err := ws.conn.ReadMessage()
			if err != nil {
				ws._close()
				fmt.Println("read msg err", err.Error())
				return
			}

			ws.qmu.Lock()
			ws.q = append(ws.q, queue{buff: msg})
			ws.qmu.Unlock()
		}
	}
}

func (ws *WebsocketModule) _dail(host string, port string) error {

	u := url.URL{Scheme: "ws", Host: host + ":" + port, Path: "/echo"}

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
	}
	close(ws.done)
}

func (ws *WebsocketModule) Close(L *lua.LState) int {
	ws._close()
	L.Push(lua.LString("succ"))
	return 1
}

func (ws *WebsocketModule) read_msg(L *lua.LState) int {
	if ws.conn == nil {
		L.Push(lua.LString("fail"))
		L.Push(lua.LString("not connected"))
		return 2
	}

	ws.qmu.Lock()
	if len(ws.q) == 0 {
		ws.qmu.Unlock()
		L.Push(lua.LString("succ"))
		L.Push(lua.LString("nodata"))
		return 2
	}

	L.Push(lua.LString("succ"))
	L.Push(lua.LString(ws.q[0].buff))
	ws.q = ws.q[1:]
	ws.qmu.Unlock()

	return 2
}

func (ws *WebsocketModule) write_msg(L *lua.LState) int {

	if ws.conn == nil {
		L.Push(lua.LString("not connected"))
		return 1
	}

	msgid := L.ToInt(1)
	msgbody := L.ToString(2)
	err := ws.conn.WriteMessage(msgid, []byte(msgbody))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}
