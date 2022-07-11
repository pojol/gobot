package script

import (
	"crypto/md5"

	lua "github.com/yuin/gopher-lua"
)

type MD5Module struct {
}

func (m *MD5Module) Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"sum": m.Sum,
	})
	l.Push(mod)
	return 1
}

func (m *MD5Module) doSum(l *lua.LState, dat []byte) (lua.LString, error) {

	md5h := md5.New()
	md5h.Write(dat)

	return lua.LString(md5h.Sum(nil)), nil

}

func (m *MD5Module) Sum(l *lua.LState) int {

	v, err := m.doSum(l, []byte(l.ToString(1)))
	l.Push(v)

	if err != nil {
		l.Push(lua.LString(err.Error()))
	} else {
		l.Push(lua.LNil)
	}

	return 2

}
