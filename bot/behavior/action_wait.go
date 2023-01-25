package behavior

import (
	"time"
)

type WaitAction struct {
	INod
	base Node

	wait    int64
	endtime int64
}

func (a *WaitAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)
	a.wait = int64(t.Wait)
}

func (a *WaitAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *WaitAction) getThread() int {
	return a.base.getThread()
}

func (a *WaitAction) getMode() Mode {
	return a.base.getMode()
}

func (a *WaitAction) getID() string {
	return a.base.ID()
}

func (a *WaitAction) getType() string {
	return WAIT
}

func (a *WaitAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *WaitAction) onTick(t *Tick) error {
	a.base.onTick(t)

	if a.endtime == 0 {
		a.endtime = time.Now().UnixNano()/1000000 + int64(a.wait)
	}

	return nil
}

func (a *WaitAction) onNext(t *Tick) {

	var currTime int64 = time.Now().UnixNano() / 1000000
	if currTime >= a.endtime {
		a.endtime = 0

		if a.base.ChildrenNum() > 0 {
			t.blackboard.Append([]INod{a.base.Children()[0]})
		} else {
			a.base.parent.onNext(t)
		}

	} else {
		t.blackboard.Append([]INod{a})
	}

}

func (a *WaitAction) onReset() {
	a.endtime = 0

	for _, child := range a.base.Children() {
		child.onReset()
	}
}
