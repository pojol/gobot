package behavior

type SelectAction struct {
	INod
	base Node

	step int
}

func (a *SelectAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)
}

func (a *SelectAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *SelectAction) getThread() int {
	return a.base.getThread()
}

func (a *SelectAction) getID() string {
	return a.base.ID()
}

func (a *SelectAction) getType() string {
	return SELETE
}

func (a *SelectAction) getMode() Mode {
	return a.base.getMode()
}

func (a *SelectAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *SelectAction) onTick(t *Tick) error {
	a.base.onTick(t)

	if a.base.ChildrenNum() <= 0 {
		goto ext
	}

	if a.step != 0 {
		self := a.base.Children()[a.step-1].(*ConditionAction)
		if self.succ {
			a.step = a.base.ChildrenNum()
			goto ext
		}
	}

ext:

	return nil
}

func (a *SelectAction) onNext(t *Tick) {

	if a.step < a.base.ChildrenNum() {
		a.step++
		t.blackboard.Append([]INod{a.base.Children()[a.step-1]})
	} else {
		a.base.parent.onNext(t)
	}

}

func (a *SelectAction) onReset() {
	a.step = 0

	for _, child := range a.base.Children() {
		child.onReset()
	}

}
