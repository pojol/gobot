package behavior

import (
	"fmt"
)

type LoopAction struct {
	INod
	base Node

	loop    int64
	curLoop int64
}

func (a *LoopAction) Init(t *Tree, parent INod, mode Mode) {
	a.base.Init(t, parent, mode)
	a.loop = int64(t.Loop)
	if a.loop == 0 {
		a.loop = 100000000000 // 无限循环 一千亿
	}
}

func (a *LoopAction) AddChild(nod INod) {
	a.base.AddChild(nod)
}

func (a *LoopAction) getType() string {
	return LOOP
}

func (a *LoopAction) getBase() *Node {
	return &a.base
}

func (a *LoopAction) onTick(t *Tick) error {

	a.base.onTick(t)

	if a.loop <= 0 {
		return fmt.Errorf("%v node %v", a.base.Type(), a.base.ID())
	}

	return nil
}

func (a *LoopAction) onNext(t *Tick) {

	if a.base.ChildrenNum() > 0 && a.curLoop < a.loop {
		a.curLoop++

		for _, child := range a.base.Children() {
			child.onReset()
		}

		t.blackboard.Append([]INod{a.base.Children()[0]})
	} else {
		a.base.parent.onNext(t)
	}

}

func (a *LoopAction) onReset() {

	a.curLoop = 0
	for _, child := range a.base.Children() {
		child.onReset()
	}

}
