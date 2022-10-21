package behavior

import "fmt"

type LoopAction struct {
	INod
	base Node

	loop    int
	curLoop int
}

func (a *LoopAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)
	a.loop = int(t.Loop)
}

func (a *LoopAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *LoopAction) getThread() int {
	return a.base.getThread()
}

func (a *LoopAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *LoopAction) onTick(t *Tick) {

	if a.base.mode == Step {
		var err error

		if a.loop <= 0 {
			err = fmt.Errorf("%v node %v", a.base.Type(), a.base.ID())
		}

		t.blackboard.ThreadFillInfo(ThreadInfo{
			Number: a.getThread(),
			CurNod: a.base.ID(),
		}, err)
	}

}

func (a *LoopAction) onNext(t *Tick) {

	if a.base.ChildrenNum() > 0 && a.curLoop < a.loop {
		a.curLoop++

		for _, child := range a.base.Children() {
			child.onReset()
		}

		t.blackboard.Append([]INod{a.base.Children()[0]})
	} else {
		a.base.parent.onNext(t)
	}

}

func (a *LoopAction) onReset() {

	a.curLoop = 0
	for _, child := range a.base.Children() {
		child.onReset()
	}

}
