package behavior

import "fmt"

type RootAction struct {
	INod
	Nod
}

func (a *RootAction) Init(t *Tree) {
	a.Nod.Init(t)
}
func (a *RootAction) ID() string {
	return a.Nod.id
}

func (a *RootAction) AddChild(child INod, parent INod) {
	a.Nod.AddChild(child, parent)
}

func (a *RootAction) Close(t *Tick) {
}

func (a *RootAction) onTick(t *Tick) NodStatus {
	fmt.Println(a.Nod.tree.Ty, a.Nod.id)
	return NSSucc
}

func (a *RootAction) onNext(t *Tick) {
	if len(a.Nod.child) > 0 {
		t.blackboard.Append([]INod{a.Nod.child[0]})
	}
}

func (a *RootAction) onReset() {

}
