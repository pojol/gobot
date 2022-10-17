package behavior

import "fmt"

type SelectAction struct {
	INod
	Nod
	step int
}

func (a *SelectAction) Init(t *Tree) {
	a.Nod.Init(t)
}
func (a *SelectAction) ID() string {
	return a.Nod.id
}
func (a *SelectAction) AddChild(child INod, parent INod) {
	a.Nod.AddChild(child, parent)
}

func (a *SelectAction) Close(t *Tick) {
}

func (a *SelectAction) onTick(t *Tick) NodStatus {

	fmt.Println(a.Nod.tree.Ty, a.Nod.id)

	childnum := len(a.Nod.child)
	if childnum <= 0 {
		a.err = fmt.Errorf("node %v not children", a.Nod.id)
		return NSErr
	}

	self := a.Nod.child[a.step].(*ConditionAction)
	if self.succ {
		a.step = len(a.Nod.child) - 1
		return NSSucc
	}

	if a.step < len(a.Nod.child) {
		a.step++
	}

	return NSSucc
}

func (a *SelectAction) onNext(t *Tick) {

	self := a.Nod.child[a.step-1].(*ConditionAction)

	if a.step < len(a.Nod.child) && self.succ {
		t.blackboard.Append([]INod{a.Nod.child[a.step]})
	} else {
		a.parent.onNext(t)
	}

}

func (a *SelectAction) onReset() {
	a.step = 0

	for _, child := range a.Nod.child {
		child.onReset()
	}
}
