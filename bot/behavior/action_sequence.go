package behavior

import "fmt"

type SequenceAction struct {
	INod
	Nod
	step int
}

func (a *SequenceAction) Init(t *Tree) {
	a.Nod.Init(t)
}
func (a *SequenceAction) ID() string {
	return a.Nod.id
}
func (a *SequenceAction) AddChild(child INod, parent INod) {
	a.Nod.AddChild(child, parent)
}

func (a *SequenceAction) Close(t *Tick) {
}

func (a *SequenceAction) OnTick(t *Tick) NodStatus {
	fmt.Println(a.Nod.tree.Ty, a.Nod.id)
	ns := NSSucc

	if a.step < len(a.Nod.child) {
		a.step++
	}

	return ns
}

func (a *SequenceAction) onNext(t *Tick) {

	if a.step < len(a.Nod.child) {
		t.blackboard.Append([]INod{a.Nod.child[a.step]})
	} else {
		a.parent.onNext(t)
	}

}

func (a *SequenceAction) onReset() {
	a.step = 0

	for _, child := range a.Nod.child {
		child.onReset()
	}
}
