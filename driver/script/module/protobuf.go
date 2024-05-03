package script

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/jsonpb"
	ggp "github.com/gogo/protobuf/proto"
	lua "github.com/yuin/gopher-lua"
)

type ProtoModule struct {
}

func (p *ProtoModule) Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"marshal":   p.Marshal,
		"unmarshal": p.Unmarshal,
	})
	registerHttpResponseType(mod, l)
	l.Push(mod)
	return 1
}

func (p *ProtoModule) doMarshal(msgty string, jstr string) (lua.LString, error) {
	var err error
	var byt []byte

	t := ggp.MessageType(msgty)
	if t == nil {
		return lua.LString(""), fmt.Errorf("unknow proto message type %v", msgty)
	}
	tptr := reflect.New(t.Elem())
	tins := tptr.Interface().(ggp.Message)

	if jstr != "[]" {
		err = jsonpb.Unmarshal(bytes.NewBufferString(jstr), tins)
		if err != nil {
			return lua.LString(""), fmt.Errorf(" jsonpb.unmarshal %v = %v %w", msgty, jstr, err)
		}

		byt, err = ggp.Marshal(tins)
		if err != nil {
			return lua.LString(""), fmt.Errorf("proto.marshal %w", err)
		}
	} else {
		byt = []byte("")
	}

	return lua.LString(byt), nil
}

func (p *ProtoModule) Marshal(L *lua.LState) int {

	res, err := p.doMarshal(L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(res)
	return 1
}

func (p *ProtoModule) doUnmarshal(msgty string, buf string) (lua.LString, error) {
	var err error

	t := ggp.MessageType(msgty)
	if t == nil {
		return lua.LString(""), fmt.Errorf("unknow proto message type %v", msgty)
	}

	tptr := reflect.New(t.Elem())
	tins, ok := tptr.Interface().(ggp.Message)
	if !ok {
		return lua.LString(""), errors.New("msg type no proto.message")
	}

	err = ggp.Unmarshal([]byte(buf), tins)
	if err != nil {
		return lua.LString(""), fmt.Errorf("proto.unmarshal %w", err)
	}

	byt, err := json.Marshal(tins)
	if err != nil {
		return lua.LString(""), fmt.Errorf("josn.marshal %w", err)
	}

	return lua.LString(byt), err
}

func (p *ProtoModule) Unmarshal(L *lua.LState) int {
	t, err := p.doUnmarshal(L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(t)
	return 1
}
