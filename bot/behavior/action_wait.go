package behavior

import (
	"fmt"
	"time"
)

type WaitAction struct {
	INod

	child  []INod
	parent INod

	id string
	ty string

	wait    int64
	endtime int64

	threadnum int
}

func (a *WaitAction) Init(t *Tree, parent INod) {
	a.id = t.ID
	a.ty = t.Ty
	a.wait = int64(t.Wait)

	a.parent = parent
}

func (a *WaitAction) ID() string {
	return a.id
}

func (a *WaitAction) setThread(num int) {
	if a.threadnum == 0 {
		a.threadnum = num
	}
}

func (a *WaitAction) getThread() int {
	if a.threadnum != 0 {
		return a.threadnum
	} else {
		return a.parent.getThread()
	}
}

func (a *WaitAction) AddChild(child INod) {
	a.child = append(a.child, child)
}

func (a *WaitAction) onTick(t *Tick) NodStatus {
	fmt.Println("\t", a.ty, a.id)
	if a.endtime == 0 {
		a.endtime = time.Now().UnixNano()/1000000 + int64(a.wait)
	}

	t.blackboard.ThreadFillInfo(ThreadInfo{
		Num:    a.getThread(),
		ErrMsg: "",
		CurNod: a.id,
	})

	return NSSucc
}

func (a *WaitAction) onNext(t *Tick) {

	var currTime int64 = time.Now().UnixNano() / 1000000
	if currTime >= a.endtime {
		a.endtime = 0

		if len(a.child) > 0 {
			t.blackboard.Append([]INod{a.child[0]})
		} else {
			a.parent.onNext(t)
		}

	} else {
		t.blackboard.Append([]INod{a})
	}

}

func (a *WaitAction) onReset() {
	a.endtime = 0

	for _, child := range a.child {
		child.onReset()
	}
}
