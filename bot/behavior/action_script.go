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
	base Node

	code string
}

func (a *ScriptAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)

	a.code = t.Code
}

func (a *ScriptAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *ScriptAction) getThread() int {
	return a.base.getThread()
}

func (a *ScriptAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *ScriptAction) onTick(t *Tick) {
	var v lua.LValue
	var err error
	a.base.onTick(t)

	err = pool.DoString(t.bs.L, a.code)
	if err != nil {
		err = fmt.Errorf("%v node %v dostring \n%w", a.base.Type(), a.base.ID(), err)
		goto ext
	}

	err = t.bs.L.CallByParam(lua.P{
		Fn:      t.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		err = fmt.Errorf("%v node %v execute \n%w", a.base.Type(), a.base.ID(), err)
		goto ext
	}

	v = t.bs.L.Get(-1)
	t.bs.L.Pop(1)

ext:

	if a.base.mode == Step {

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
			Number: a.getThread(),
			CurNod: a.base.ID(),
			Change: changeStr,
		}
		t.blackboard.ThreadFillInfo(info, err)
	}
}

func (a *ScriptAction) onNext(t *Tick) {

	if a.base.ChildrenNum() > 0 && !a.base.GetFreeze() {
		a.base.SetFreeze(true)
		t.blackboard.Append([]INod{a.base.Children()[0]})
	} else {
		a.base.parent.onNext(t)
	}

}

func (a *ScriptAction) onReset() {
	a.base.SetFreeze(false)
	for _, child := range a.base.Children() {
		child.onReset()
	}
}
