package behavior

type ParallelAction struct {
	INod
	base Node
}

func (a *ParallelAction) Init(t *Tree, parent INod) {
	a.base.Init(t, parent)
}

func (a *ParallelAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *ParallelAction) getType() string {
	return PARALLEL
}

func (a *ParallelAction) getBase() *Node {
	return &a.base
}

func (a *ParallelAction) onTick(t *Tick) error {
	a.base.onTick(t)
	return nil
}

func (a *ParallelAction) onNext(t *Tick) {
	if !a.base.GetFreeze() {
		a.base.SetFreeze(true)

		for _, children := range a.base.Children() {
			t.blackboard.Append([]INod{children})

			newthreadnum := t.blackboard.ThreadCurNum() + 1
			t.blackboard.ThreadAdd(newthreadnum)

			children.getBase().setThread(newthreadnum)
		}

	} else {
		a.base.threadNumber++

		if a.base.threadNumber >= a.base.ChildrenNum() {
			a.base.parent.onNext(t)
		}
	}
}

func (a *ParallelAction) onReset() {
	a.base.SetFreeze(false)
	a.base.threadNumber = 0

	for _, child := range a.base.Children() {
		child.onReset()
	}
}
