package behavior

import "fmt"

type LoopAction struct {
	INod
	Nod
	loop int
}

func (a *LoopAction) Init(t *Tree) {
	a.Nod.Init(t)
}
func (a *LoopAction) ID() string {
	return a.Nod.id
}

func (a *LoopAction) AddChild(child INod, parent INod) {
	a.Nod.AddChild(child, parent)
}

func (a *LoopAction) Close(t *Tick) {
}

func (a *LoopAction) onTick(t *Tick) NodStatus {
	fmt.Println(a.Nod.tree.Ty, a.Nod.id)
	a.loop++
	return NSSucc
}

func (a *LoopAction) onNext(t *Tick) {

	childnum := len(a.Nod.child)

	if childnum > 0 && a.loop < a.Nod.loop {
		t.blackboard.Append([]INod{a.Nod.child[0]})
	} else {
		a.Nod.parent.onNext(t)
	}
}

func (a *LoopAction) onReset() {

	a.loop = 0
	for _, child := range a.Nod.child {
		child.onReset()
	}

}
