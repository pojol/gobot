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

func (a *RootAction) getType() string {
	return ROOT
}

func (a *RootAction) getBase() *Node {
	return &a.base
}

func (a *RootAction) onTick(t *Tick) error {
	a.base.onTick(t)

	return nil
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
