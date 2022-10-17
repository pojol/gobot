package behavior

import (
	"fmt"

	"github.com/pojol/gobot/bot/state"
	lua "github.com/yuin/gopher-lua"
)

type ConditionAction struct {
	INod
	Nod
}

func (a *ConditionAction) Init(t *Tree) {
	a.Nod.Init(t)
}
func (a *ConditionAction) ID() string {
	return a.Nod.id
}
func (a *ConditionAction) AddChild(child INod, parent INod) {
	a.Nod.AddChild(child, parent)
}
func (a *ConditionAction) Close(t *Tick) {
}

func (a *ConditionAction) onTick(t *Tick) NodStatus {
	fmt.Println(a.Nod.tree.Ty, a.Nod.id)
	err := state.DoString(t.bs.L, a.code)
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

	a.succ = lua.LVAsBool(v)

	return NSSucc
}

func (a *ConditionAction) onNext(t *Tick) {
	if len(a.Nod.child) > 0 {
		t.blackboard.Append([]INod{a.Nod.child[0]})
	} else {
		a.parent.onNext(t)
	}

}

func (a *ConditionAction) onReset() {
	a.Nod.succ = false

	for _, child := range a.Nod.child {
		child.onReset()
	}
}
