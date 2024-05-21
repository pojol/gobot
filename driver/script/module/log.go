package script

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pojol/gobot/driver/utils"
	lua "github.com/yuin/gopher-lua"
)

type LogModule struct {
	logs []string
	sync.Mutex
}

func (lm *LogModule) Loader(L *lua.LState) int {
	log := L.NewTable()

	// 设置 log 的方法
	L.SetFuncs(log, map[string]lua.LGFunction{
		"info": lm.info,
	})

	L.SetGlobal("log", log)

	return 1
}

// info - 脚本端添加一行日志
func (lm *LogModule) info(L *lua.LState) int {
	lm.Lock()
	defer lm.Unlock()

	var info string

	// 获取参数数量
	n := L.GetTop()
	if n == 0 {
		return 0
	}

	for i := 1; i <= n; i++ {
		arg := L.Get(i)
		switch arg.Type() {
		case lua.LTNumber:
			info += fmt.Sprintf("%v\t", arg.(lua.LNumber))
		case lua.LTString:
			info += fmt.Sprintf("%v\t", string(arg.(lua.LString)))
		case lua.LTBool:
			b := bool(arg.(lua.LBool))
			if b {
				info += "true\t"
			} else {
				info += "false\t"
			}
		case lua.LTNil:
			info += "nil\t"
		case lua.LTTable:
			m := make(map[string]interface{})
			var err error
			m, err = utils.Table2Map(arg.(*lua.LTable))
			if err != nil {
				fmt.Println("Table2Map error", err)
			}
			tablestr, _ := json.MarshalIndent(&m, "", "    ")
			info += string(tablestr) + "\n"
		// 根据需要处理其他类型
		default:
			info += fmt.Sprintf("arg%d is of type %s\n", i, arg.Type().String())
		}
	}

	lm.logs = append(lm.logs, info)

	return 0
}

func (lm *LogModule) Pop() string {
	lm.Lock()
	defer lm.Unlock()

	n := len(lm.logs)
	if n == 0 {
		return ""
	}

	// 从头部取出元素
	x := lm.logs[0]
	lm.logs = lm.logs[1:n]

	return x
}

func (lm *LogModule) Clean() {
	lm.Lock()
	lm.logs = []string{}
	lm.Unlock()
}
