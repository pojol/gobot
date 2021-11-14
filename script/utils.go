package script

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

type UtilsModule struct {
}

func (u *UtilsModule) Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"random": u.Random,
		"uuid":   u.UUID,
	})
	l.Push(mod)
	return 1
}

func (u *UtilsModule) doRandom(l *lua.LState, n int) (lua.LNumber, error) {

	if n <= 0 {
		return lua.LNumber(0), errors.New("")
	}

	return lua.LNumber(rand.Intn(n)), nil

}

func (u *UtilsModule) Random(l *lua.LState) int {

	v, err := u.doRandom(l, l.ToInt(1))

	l.Push(v)
	if err != nil {
		l.Push(lua.LString(err.Error()))
	} else {
		l.Push(lua.LNil)
	}

	return 2

}

func (u *UtilsModule) UUID(l *lua.LState) int {
	l.Push(lua.LString(uuid.NewString()))
	return 1
}
