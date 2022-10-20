package behavior

import (
	"fmt"
)

type LoopAction struct {
	INod

	child  []INod
	parent INod

	id string
	ty string

	loop    int
	curLoop int

	threadnum int
}

func (a *LoopAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty

	a.loop = int(t.Loop)

	a.parent = parent
}

func (a *LoopAction) ID() string {
	return a.id
}

func (a *LoopAction) setThread(num int) {
	if a.threadnum == 0 {
		a.threadnum = num
	}
}

func (a *LoopAction) getThread() int {
	if a.threadnum != 0 {
		return a.threadnum
	} else {
		return a.parent.getThread()
	}
}

func (a *LoopAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *LoopAction) onTick(t *Tick) NodStatus {
	fmt.Println("\t", a.ty, a.id)

	t.blackboard.ThreadFillInfo(ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
	})

	return NSSucc
}

func (a *LoopAction) onNext(t *Tick) {

	childnum := len(a.child)
	if childnum > 0 && a.curLoop < a.loop {
		a.curLoop++

		for _, child := range a.child {
			child.onReset()
		}

		t.blackboard.Append([]INod{a.child[0]})
	} else {
		a.parent.onNext(t)
	}

}

func (a *LoopAction) onReset() {

	a.curLoop = 0
	for _, child := range a.child {
		child.onReset()
	}

}
