package script

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type LogModule struct {
	logs []string
	sync.Mutex
}

func (lm *LogModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"info": lm.info,
	})
	L.Push(mod)
	return 1
}

// info - 脚本端添加一行日志
func (lm *LogModule) info(L *lua.LState) int {
	lm.Lock()
	defer lm.Unlock()

	info := L.ToString(1)
	lm.logs = append(lm.logs, info)

	return 0
}

// Push - golang端添加一行日志
func (lm *LogModule) Push(log string) {
	lm.Lock()
	lm.logs = append(lm.logs, log)
	lm.Unlock()
}

func (lm *LogModule) Pop() string {
	lm.Lock()
	defer lm.Unlock()

	n := len(lm.logs)
	if n == 0 {
		return ""
	}
	x := lm.logs[n-1]
	lm.logs = lm.logs[0 : n-1]

	return x
}

func (lm *LogModule) Clean() {
	lm.Lock()
	lm.logs = []string{}
	lm.Unlock()
}
