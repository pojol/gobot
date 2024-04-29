package script

import (
	"encoding/base64"

	lua "github.com/yuin/gopher-lua"
)

type Base64Module struct {
}

func (b *Base64Module) Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"encode": b.Encode,
		"decode": b.Decode,
	})
	l.Push(mod)
	return 1
}

func (b *Base64Module) doEncode(s string) lua.LString {

	ds := base64.StdEncoding.EncodeToString([]byte(s))
	return lua.LString(ds)
}

func (b *Base64Module) Encode(l *lua.LState) int {

	v := b.doEncode(l.ToString(1))
	l.Push(v)

	return 1
}

func (b *Base64Module) doDecode(s string) (lua.LString, error) {
	byt, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return lua.LString(""), err
	}

	return lua.LString(string(byt)), nil
}

func (b *Base64Module) Decode(l *lua.LState) int {

	v, err := b.doDecode(l.ToString(1))
	l.Push(v)

	if err != nil {
		l.Push(lua.LString(err.Error()))
	} else {
		l.Push(lua.LNil)
	}

	return 2
}
