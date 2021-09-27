package behavior

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type ProtoModule struct {
}

func (p *ProtoModule) Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"marshal": p.Marshal,
	})
	registerHttpResponseType(mod, l)
	l.Push(mod)
	return 1
}

func (p *ProtoModule) doMarshal(L *lua.LState, ty string, msg *lua.LTable) (lua.LString, error) {
	var body []byte
	var err error

	return lua.LString(body), err
}

func (p *ProtoModule) Marshal(L *lua.LState) int {
	res, err := p.doMarshal(L, L.ToString(1), L.ToTable(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(res)
	return 1
}
