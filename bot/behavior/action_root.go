package behavior

import "fmt"

type RootAction struct {
	INod
	base Node
}

func (a *RootAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)
	a.base.threadNumber = 1
}

func (a *RootAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *RootAction) getThread() int {
	return a.base.getThread()
}

func (a *RootAction) setThread(tn int) {
	a.base.setThread(tn)
}

func (a *RootAction) onTick(t *Tick) {
	a.base.onTick(t)

	if a.base.mode == Step {
		t.blackboard.ThreadFillInfo(ThreadInfo{
			Number: a.getThread(),
			ErrMsg: "",
			CurNod: a.base.ID(),
		}, nil)
	}

}

func (a *RootAction) onNext(t *Tick) {
	if a.base.ChildrenNum() > 0 && !a.base.GetFreeze() {
		a.base.SetFreeze(true)
		t.blackboard.Append([]INod{a.base.Children()[0]})
	} else {
		fmt.Println(t.botid, "end")
		t.blackboard.End()
	}
}

func (a *RootAction) onReset() {

}
