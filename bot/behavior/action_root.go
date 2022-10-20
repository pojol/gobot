package behavior

import (
	"fmt"
)

type RootAction struct {
	INod

	child []INod

	id string
	ty string

	freeze bool

	threadnum int
}

func (a *RootAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty

	a.threadnum = 1
}

func (a *RootAction) ID() string {
	return a.id
}

func (a *RootAction) setThread(num int) {
}

func (a *RootAction) getThread() int {
	return a.threadnum
}

func (a *RootAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *RootAction) onTick(t *Tick) NodStatus {
	fmt.Println("\t", a.ty, a.id)

	t.blackboard.ThreadFillInfo(ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
	})

	return NSSucc
}

func (a *RootAction) onNext(t *Tick) {
	if len(a.child) > 0 && !a.freeze {
		a.freeze = true
		t.blackboard.Append([]INod{a.child[0]})
	} else {
		t.blackboard.End()
	}
}

func (a *RootAction) onReset() {

}
