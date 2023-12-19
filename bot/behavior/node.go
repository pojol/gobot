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

//
const (
	Error = "Error"
	Succ  = "Succ"
	Exit  = "Exit"
	Break = "Break"
)

type INod interface {
	Init(*Tree, INod, Mode)
	AddChild(INod)

	getBase() *Node
	getType() string

	onTick(*Tick) error
	onNext(*Tick)
	onReset()
}

type Node struct {
	id string
	ty string

	child  []INod
	parent INod

	mode Mode

	freeze       bool
	threadNumber int
}

func (n *Node) Init(t *Tree, parent INod, mode Mode) {
	n.id = t.ID
	n.ty = t.Ty
	n.mode = mode

	n.parent = parent
}

func (a *Node) ID() string {
	return a.id
}

func (a *Node) Type() string {
	return a.ty
}

func (a *Node) GetFreeze() bool {
	return a.freeze
}

// SetFreeze 使节点无效（已经执行过的节点）
func (a *Node) SetFreeze(f bool) {
	a.freeze = f
}

func (a *Node) ChildrenNum() int {
	return len(a.child)
}

func (a *Node) Children() []INod {
	return a.child
}

func (a *Node) onTick(t *Tick) {
}

func (a *Node) getMode() Mode {
	return a.mode
}

func (a *Node) setThread(number int) {
	if a.threadNumber == 0 {
		a.threadNumber = number
	}
}

func (a *Node) getThread() int {

	if a.threadNumber != 0 {
		return a.threadNumber
	} else {
		return a.parent.getBase().getThread()
	}
}

func (a *Node) AddChild(child INod) {
	a.child = append(a.child, child)
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
