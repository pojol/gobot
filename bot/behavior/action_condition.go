package behavior

import (
	"fmt"

	"github.com/pojol/gobot/bot/pool"
	lua "github.com/yuin/gopher-lua"
)

type ConditionAction struct {
	INod

	child  []INod
	parent INod

	id   string
	ty   string
	code string

	succ   bool
	freeze bool

	threadnum int

	err error
}

func (a *ConditionAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty
	a.code = t.Code

	a.parent = parent
}

func (a *ConditionAction) ID() string {
	return a.id
}

func (a *ConditionAction) setThread(num int) {
	if a.threadnum == 0 {
		a.threadnum = num
	}
}

func (a *ConditionAction) getThread() int {
	if a.threadnum != 0 {
		return a.threadnum
	} else {
		return a.parent.getThread()
	}
}

func (a *ConditionAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *ConditionAction) onTick(t *Tick) NodStatus {
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

	a.succ = lua.LVAsBool(v)
	t.blackboard.ThreadFillInfo(ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
	})

	return NSSucc
}

func (a *ConditionAction) onNext(t *Tick) {

	if len(a.child) > 0 && !a.freeze {
		a.freeze = true
		child := a.child[0]
		t.blackboard.Append([]INod{child})
	} else {
		a.parent.onNext(t)
	}

}

func (a *ConditionAction) onReset() {
	a.succ = false
	a.freeze = false

	for _, child := range a.child {
		child.onReset()
	}
}
