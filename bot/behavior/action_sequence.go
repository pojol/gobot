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

func (a *SequenceAction) getMode() Mode {
	return a.base.getMode()
}

func (a *SequenceAction) getType() string {
	return SEQUENCE
}

func (a *SequenceAction) getID() string {
	return a.base.ID()
}

func (a *SequenceAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *SequenceAction) onTick(t *Tick) error {
	a.base.onTick(t)

	return nil
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
