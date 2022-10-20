package behavior

import (
	"fmt"
)

type SelectAction struct {
	INod

	child  []INod
	parent INod

	id string
	ty string

	step int
	err  error

	threadnum int
}

func (a *SelectAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty

	a.parent = parent
}

func (a *SelectAction) ID() string {
	return a.id
}

func (a *SelectAction) setThread(num int) {
	if a.threadnum == 0 {
		a.threadnum = num
	}
}

func (a *SelectAction) getThread() int {
	if a.threadnum != 0 {
		return a.threadnum
	} else {
		return a.parent.getThread()
	}
}

func (a *SelectAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *SelectAction) onTick(t *Tick) NodStatus {

	fmt.Println("\t", a.ty, a.id)

	childnum := len(a.child)
	if childnum <= 0 {
		a.err = fmt.Errorf("node %v not children", a.id)
		return NSErr
	}

	if a.step != 0 {
		self := a.child[a.step-1].(*ConditionAction)
		if self.succ {
			a.step = len(a.child)
			return NSSucc
		}
	}

	t.blackboard.ThreadFillInfo(ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
	})

	return NSSucc
}

func (a *SelectAction) onNext(t *Tick) {

	if a.step < len(a.child) {
		a.step++
		t.blackboard.Append([]INod{a.child[a.step-1]})
	} else {
		a.parent.onNext(t)
	}

}

func (a *SelectAction) onReset() {
	a.step = 0

	for _, child := range a.child {
		child.onReset()
	}

}
