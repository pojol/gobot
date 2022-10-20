package behavior

import (
	"fmt"
)

type SequenceAction struct {
	INod

	child  []INod
	parent INod

	id string
	ty string

	step int

	threadnum int
}

func (a *SequenceAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty

	a.parent = parent
}

func (a *SequenceAction) ID() string {
	return a.id
}

func (a *SequenceAction) setThread(num int) {
	if a.threadnum == 0 {
		a.threadnum = num
	}
}

func (a *SequenceAction) getThread() int {
	if a.threadnum != 0 {
		return a.threadnum
	} else {
		return a.parent.getThread()
	}
}

func (a *SequenceAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *SequenceAction) onTick(t *Tick) NodStatus {
	fmt.Println("\t", a.ty, a.id)

	t.blackboard.ThreadFillInfo(ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
	})

	return NSSucc
}

func (a *SequenceAction) onNext(t *Tick) {

	if a.step < len(a.child) {
		a.step++
		t.blackboard.Append([]INod{a.child[a.step-1]})

	} else {
		a.parent.onNext(t)
	}

}

func (a *SequenceAction) onReset() {
	a.step = 0

	for _, child := range a.child {
		child.onReset()
	}
}
