package script

import (
	"fmt"
	"testing"

	"github.com/pojol/gobot/driver/utils"
	lua "github.com/yuin/gopher-lua"
)

func TestScript(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	path := "../../script"

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		L.DoFile(path + "/" + v)
	}

	L.DoString(`
		return state.Error, "11"
	`)

	v2 := L.Get(-1)
	fmt.Println("v2", v2)
	if v2.Type() != lua.LTNil {
		L.Pop(1)
	}

	v1 := L.Get(-1)
	if v1.Type() != lua.LTNil {
		L.Pop(1)
	}

	fmt.Println("return value =", v1, v2)
}
