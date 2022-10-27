package behavior

import (
	"fmt"

	"github.com/pojol/gobot/bot/pool"
	lua "github.com/yuin/gopher-lua"
)

type ConditionAction struct {
	INod
	base Node

	code string
	succ bool
}

func (a *ConditionAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)

	a.code = t.Code
}

func (a *ConditionAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *ConditionAction) getThread() int {
	return a.base.getThread()
}

func (a *ConditionAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *ConditionAction) onTick(t *Tick) {
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
	a.succ = lua.LVAsBool(v)

ext:
	if a.base.mode == Step {
		t.blackboard.ThreadFillInfo(ThreadInfo{
			Number: a.base.getThread(),
			CurNod: a.base.ID(),
		}, err)
	}
}

func (a *ConditionAction) onNext(t *Tick) {

	if (a.base.ChildrenNum() > 0 && !a.base.GetFreeze()) && a.succ {
		a.base.SetFreeze(true)
		child := a.base.Children()[0]
		t.blackboard.Append([]INod{child})
	} else {
		a.base.parent.onNext(t)
	}

}

func (a *ConditionAction) onReset() {
	a.succ = false
	a.base.SetFreeze(false)

	for _, child := range a.base.Children() {
		child.onReset()
	}
}
