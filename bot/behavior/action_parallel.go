package behavior

type ParallelAction struct {
	INod
	base Node
}

func (a *ParallelAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)
}

func (a *ParallelAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *ParallelAction) getThread() int {
	return a.base.getThread()
}

func (a *ParallelAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *ParallelAction) onTick(t *Tick) {
	a.base.onTick(t)

	if a.base.mode == Step {
		t.blackboard.ThreadFillInfo(ThreadInfo{
			Number: a.getThread(),
			ErrMsg: "",
			CurNod: a.base.ID(),
		}, nil)
	}

}

func (a *ParallelAction) onNext(t *Tick) {
	if !a.base.GetFreeze() {
		a.base.SetFreeze(true)

		for _, children := range a.base.Children() {
			t.blackboard.Append([]INod{children})

			newthreadnum := t.blackboard.ThreadCurNum() + 1
			t.blackboard.ThreadAdd(newthreadnum)

			children.setThread(newthreadnum)
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
