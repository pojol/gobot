package behavior

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
	"sync"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

func compileLuaByFile(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

func compileLuaByString(script string) (*lua.FunctionProto, error) {

	reader := strings.NewReader(script)

	chunk, err := parse.Parse(reader, script)
	if err != nil {
		return nil, err
	}

	proto, err := lua.Compile(chunk, script)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

func doCompiledFile(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}

func calcMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

var luamap = sync.Map{}

func DoString(L *lua.LState, script string) error {

	name := calcMD5(script)

	val, ok := luamap.Load(name)
	if !ok {
		p, err := compileLuaByString(script)
		if err != nil {
			return err
		}

		luamap.Store(name, p)
		val = p
	}

	return doCompiledFile(L, val.(*lua.FunctionProto))
}

func DoFile(L *lua.LState, filePath string) error {

	val, ok := luamap.Load(filePath)
	if !ok {
		p, err := compileLuaByFile(filePath)
		if err != nil {
			return err
		}

		luamap.Store(filePath, p)
		val = p
	}

	return doCompiledFile(L, val.(*lua.FunctionProto))
}
