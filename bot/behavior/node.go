package behavior

// Returning status
type NodStatus int

const (
	NSSucc NodStatus = 1 + iota
	NSErr
	NSFail
)

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

type INod interface {
	Init(*Tree, INod)
	ID() string
	AddChild(INod)

	getThread() int
	setThread(int)

	onTick(*Tick) NodStatus
	onNext(*Tick)
	onReset()

	GetErr() error
}

type CreateActionFunc func() interface{}

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
