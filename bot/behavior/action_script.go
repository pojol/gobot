package behavior

import (
	"fmt"

	"github.com/pojol/gobot/bot/state"
	lua "github.com/yuin/gopher-lua"
)

type ScriptAction struct {
	INod
	Nod
}

func (a *ScriptAction) Init(t *Tree) {
	a.Nod.Init(t)
}
func (a *ScriptAction) ID() string {
	return a.Nod.id
}
func (a *ScriptAction) AddChild(child INod, parent INod) {
	a.Nod.AddChild(child, parent)
}

func (a *ScriptAction) Close(t *Tick) {
}

func (a *ScriptAction) onTick(t *Tick) NodStatus {
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
	t.bs.L.Pop(1)

	return NSSucc
}

func (sa *ScriptAction) onNext(t *Tick) {

	if len(sa.Nod.child) > 0 {
		t.blackboard.Append([]INod{sa.Nod.child[0]})
	} else {
		sa.parent.onNext(t)
	}

}

func (a *ScriptAction) onReset() {
	for _, child := range a.Nod.child {
		child.onReset()
	}
}
