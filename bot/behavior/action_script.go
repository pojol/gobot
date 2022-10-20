package behavior

import (
	"encoding/json"
	"fmt"

	"github.com/pojol/gobot/bot/pool"
	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
)

type ScriptAction struct {
	INod

	child  []INod
	parent INod

	id   string
	ty   string
	code string

	freeze bool
	err    error

	threadnum int
}

func (a *ScriptAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty
	a.code = t.Code

	a.parent = parent
}

func (a *ScriptAction) ID() string {
	return a.id
}

func (a *ScriptAction) setThread(num int) {
	if a.threadnum == 0 {
		a.threadnum = num
	}
}

func (a *ScriptAction) getThread() int {
	if a.threadnum != 0 {
		return a.threadnum
	} else {
		return a.parent.getThread()
	}
}

func (a *ScriptAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *ScriptAction) onTick(t *Tick) NodStatus {
	fmt.Println("\t", a.ty, a.id)
	err := pool.DoString(t.bs.L, a.code)
	if err != nil {
		a.err = err
		return NSErr
	}

	err = t.bs.L.CallByParam(lua.P{
		Fn:      t.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		a.err = err
		return NSErr
	}

	v := t.bs.L.Get(-1)
	t.bs.L.Pop(1)

	var changeStr string

	tab, ok := v.(*lua.LTable)
	if ok {
		change, err := utils.Table2Map(tab)
		if err != nil {
			fmt.Println("script response 2 map err", err.Error())
		}

		changeByt, err := json.Marshal(&change)
		if err != nil {
			fmt.Println("marshal change info err", err.Error())
		}
		changeStr = string(changeByt)
	}

	info := ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
		Change: changeStr,
	}
	t.blackboard.ThreadFillInfo(info)

	return NSSucc
}

func (a *ScriptAction) onNext(t *Tick) {

	if len(a.child) > 0 && !a.freeze {
		a.freeze = true
		t.blackboard.Append([]INod{a.child[0]})
	} else {
		a.parent.onNext(t)
	}

}

func (a *ScriptAction) onReset() {
	a.freeze = false
	for _, child := range a.child {
		child.onReset()
	}
}
