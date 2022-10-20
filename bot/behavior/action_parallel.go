package behavior

import "fmt"

type ParallelAction struct {
	INod

	child  []INod
	parent INod

	id string
	ty string

	freeze bool

	threadnum int
}

func (a *ParallelAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty

	a.parent = parent
}

func (a *ParallelAction) ID() string {
	return a.id
}

func (a *ParallelAction) setThread(num int) {

}

func (a *ParallelAction) getThread() int {
	return a.parent.getThread()
}

func (a *ParallelAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *ParallelAction) onTick(t *Tick) NodStatus {
	fmt.Println("\t", a.ty, a.id)

	t.blackboard.ThreadFillInfo(ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
	})

	return NSSucc
}

func (a *ParallelAction) onNext(t *Tick) {
	if !a.freeze {
		a.freeze = true

		for _, children := range a.child {
			t.blackboard.Append([]INod{children})

			newthreadnum := t.blackboard.ThreadCurNum() + 1
			t.blackboard.ThreadAdd(newthreadnum)

			children.setThread(newthreadnum)
		}

	} else {
		a.threadnum++
		fmt.Println("end thread")

		if a.threadnum >= len(a.child) {
			a.parent.onNext(t)
		}
	}
}

func (a *ParallelAction) onReset() {
	a.freeze = false
	a.threadnum = 0

	for _, child := range a.child {
		child.onReset()
	}
}
