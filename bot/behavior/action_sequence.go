package behavior

type SequenceAction struct {
	INod
	base Node

	step int
}

func (a *SequenceAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)
}

func (a *SequenceAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *SequenceAction) getThread() int {
	return a.base.getThread()
}

func (a *SequenceAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *SequenceAction) onTick(t *Tick) {
	a.base.onTick(t)

	if a.base.mode == Step {
		t.blackboard.ThreadFillInfo(ThreadInfo{
			Number: a.getThread(),
			ErrMsg: "",
			CurNod: a.base.ID(),
		}, nil)
	}

}

func (a *SequenceAction) onNext(t *Tick) {

	if a.step < a.base.ChildrenNum() {
		a.step++
		t.blackboard.Append([]INod{a.base.Children()[a.step-1]})

	} else {
		a.base.parent.onNext(t)
	}

}

func (a *SequenceAction) onReset() {
	a.step = 0

	for _, child := range a.base.Children() {
		child.onReset()
	}
}
