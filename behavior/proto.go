package behavior

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/jsonpb"
	ggp "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
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

func (p *ProtoModule) doMarshal(L *lua.LState, msgty string, jstr string) (lua.LString, error) {
	var err error

	t := ggp.MessageType(msgty).Elem()
	tptr := reflect.New(t)
	tins := tptr.Interface().(proto.Message)

	err = jsonpb.Unmarshal(bytes.NewBufferString(jstr), tins)
	if err != nil {
		return lua.LString(""), fmt.Errorf("jsonpb.unmarshal %w", err)
	}

	byt, err := proto.Marshal(tins)
	if err != nil {
		return lua.LString(""), fmt.Errorf("proto.marshal %w", err)
	}

	return lua.LString(byt), nil
}

func (p *ProtoModule) Marshal(L *lua.LState) int {

	res, err := p.doMarshal(L, L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(res)
	return 1
}

func (p *ProtoModule) doUnmarshal(L *lua.LState, msgty string, buf string) (lua.LString, error) {
	var err error

	t := ggp.MessageType(msgty).Elem()
	tptr := reflect.New(t)
	tins, ok := tptr.Interface().(proto.Message)
	if !ok {
		return lua.LString(""), errors.New("msg type no proto.message")
	}

	err = proto.Unmarshal([]byte(buf), tins)
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
	t, err := p.doUnmarshal(L, L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(t)
	return 1
}