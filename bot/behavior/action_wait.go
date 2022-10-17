package behavior

import (
	"fmt"
	"time"
)

type WaitAction struct {
	INod
	Nod
	endtime int64
}

func (a *WaitAction) Init(t *Tree) {
	a.Nod.Init(t)
}
func (a *WaitAction) ID() string {
	return a.Nod.id
}
func (a *WaitAction) AddChild(child INod, parent INod) {
	a.Nod.AddChild(child, parent)
}
func (a *WaitAction) Close(t *Tick) {
}

func (a *WaitAction) onTick(t *Tick) NodStatus {
	fmt.Println(a.Nod.tree.Ty, a.Nod.id)
	if a.endtime == 0 {
		a.endtime = time.Now().Unix() + int64(a.wait)
	}

	return NSSucc
}

func (a *WaitAction) onNext(t *Tick) {

	if time.Now().Unix() >= a.endtime {

		if len(a.Nod.child) > 0 {
			t.blackboard.Append([]INod{a.Nod.child[0]})
		} else {
			a.parent.onNext(t)
		}

	}

}

func (a *WaitAction) onReset() {
	a.endtime = 0

	for _, child := range a.Nod.child {
		child.onReset()
	}
}
