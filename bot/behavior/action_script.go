package behavior

import (
	"fmt"

	"github.com/pojol/gobot/bot/pool"
	lua "github.com/yuin/gopher-lua"
)

type ScriptAction struct {
	INod
	base Node

	code string
}

func (a *ScriptAction) Init(t *Tree, parent INod) {
	a.base.Init(t, parent)

	a.code = t.Code
}

func (a *ScriptAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *ScriptAction) getBase() *Node {
	return &a.base
}

func (a *ScriptAction) getType() string {
	return SCRIPT
}

func (a *ScriptAction) onTick(t *Tick) error {
	var err error
	a.base.onTick(t)

	err = pool.DoString(t.bs.L, a.code)
	if err != nil {
		err = fmt.Errorf("%v node %v dostring \n%w", a.base.ID(), a.base.Type(), err)
		goto ext
	}

	for i := 0; i < t.bs.L.GetTop(); i++ {
		t.bs.L.Pop(1) // clean stack
	}

	err = t.bs.L.CallByParam(lua.P{
		Fn:      t.bs.L.GetGlobal("execute"),
		NRet:    2,
		Protect: true,
	})
	if err != nil {
		err = fmt.Errorf("%v node %v execute \n%w", a.base.ID(), a.base.Type(), err)
		goto ext
	}

ext:
	return err
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
