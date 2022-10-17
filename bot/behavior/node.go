package behavior

import "fmt"

/*
const (
	COMPOSITE = "composite"
	DECORATOR = "decorator"
	ACTION    = "action"
	CONDITION = "condition"
)
*/

// Returning status
type NodStatus int

const (
	NSSucc NodStatus = 1 + iota
	NSErr
	NSFail
)

type CreateActionFunc func() interface{}

const (
	ROOT      = "RootNode"
	SELETE    = "SelectorNode"
	SEQUENCE  = "SequenceNode"
	CONDITION = "ConditionNode"
	WAIT      = "WaitNode"
	LOOP      = "LoopNode"
	PARALLEL  = "ParallelNode"
	SCRIPT    = "ScriptNode"
)

var actionFactory map[string]CreateActionFunc = map[string]CreateActionFunc{
	ROOT:      func() interface{} { return &RootAction{} },
	SELETE:    func() interface{} { return &SelectAction{} },
	SEQUENCE:  func() interface{} { return &SequenceAction{} },
	CONDITION: func() interface{} { return &ConditionAction{} },
	WAIT:      func() interface{} { return &WaitAction{} },
	LOOP:      func() interface{} { return &LoopAction{} },
	PARALLEL:  func() interface{} { return &ParallelAction{} },
	SCRIPT:    func() interface{} { return &ScriptAction{} },
}

func NewNode(name string) interface{} {

	if _, ok := actionFactory[name]; ok {
		return actionFactory[name]()
	}

	return actionFactory[SCRIPT]()
}

type INod interface {
	Init(*Tree)
	ID() string
	AddChild(INod, INod)
	Close(*Tick)

	onTick(*Tick) NodStatus
	onNext(*Tick)
	onReset()

	GetErr() error
}

type Nod struct {
	tree *Tree

	child  []INod
	parent INod

	succ bool

	id   string
	wait int
	loop int
	code string

	err error
}

func (n *Nod) Init(bt *Tree) {
	fmt.Println(bt.Ty, bt.ID, "init")
	n.id = bt.ID
	n.tree = bt
	n.wait = int(bt.Wait)
	n.loop = int(bt.Loop)
	n.code = bt.Code
}

func (n *Nod) AddChild(child INod, parent INod) {
	n.parent = parent
	n.child = append(n.child, child)
}

func GetErrInfo() error {
	return nil
}
