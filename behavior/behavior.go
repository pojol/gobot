package behavior

import (
	"encoding/xml"
)

const (
	ROOT      = "RootNode"
	SELETE    = "SelectorNode"
	SEQUENCE  = "SequenceNode"
	CONDITION = "ConditionNode"
	WAIT      = "WaitNode"
	LOOP      = "LoopNode"
	ASSERT    = "AssertNode"
	PARALLEL  = "ParallelNode"
)

type Tree struct {
	ID string `xml:"id"`
	Ty string `xml:"ty"`

	Wait int32 `xml:"wait"`

	Loop int32 `xml:"loop"` // 用于记录循环节点的循环次数

	Freeze bool

	Step int
	Nods int

	Code string `xml:"code"`
	Succ bool

	Parent   *Tree
	Children []*Tree `xml:"children"`
}

var tickStrategy = map[string]func(int, *Tree) []*ThreadTree{
	SELETE:    selectTick,
	SEQUENCE:  sequenceTick,
	CONDITION: conditionTick,
	WAIT:      waitTick,
	LOOP:      loopTick,
	ASSERT:    assertTick,
	"Script":  scriptTick,
	PARALLEL:  parallelTick,
}

func (t *Tree) Tick(thread int) []*ThreadTree {

	if t.Ty == ROOT {
		return []*ThreadTree{}
	}

	return tickStrategy[t.Ty](thread, t)

}

func IsScriptNode(ty string) bool {

	if ty == ROOT ||
		ty == SELETE ||
		ty == SEQUENCE ||
		ty == LOOP ||
		ty == PARALLEL {
		return false
	}

	return true
}

func (t *Tree) link(parent *Tree) {

	t.Parent = parent
	for k := range t.Children {
		t.Children[k].link(t)
	}

}

func New(f []byte) (*Tree, error) {

	tree := &Tree{
		Parent: nil,
	}
	err := xml.Unmarshal([]byte(f), &tree)
	if err != nil {
		panic(err)
	}

	for k := range tree.Children {
		tree.Children[k].link(tree)
	}

	return tree, nil
}

/*
func (tree *Tree) resetChildren() {
	for k := range tree.Children {

		tree.Children[k].Step = 0

		if len(tree.Children[k].Children) > 0 {
			tree.Children[k].resetChildren()
		}

	}
}

func (tree *Tree) Next(ret bool) []*Tree {

	children := []*Tree{}

	if tree == nil {
		return children
	}

	switch tree.Ty {
	case SEQUENCE:
		if !ret && tree.Step != 0 {
			tree.Step = len(tree.Children)
		}
		if tree.Step < len(tree.Children) {
			children = append(children, tree.Children[tree.Step])
			tree.Step++
			goto ext
		}
		return tree.Parent.Next(false)
	case SELETE:
		if ret {
			tree.Step = len(tree.Children)
		}
		if tree.Step < len(tree.Children) {
			children = append(children, tree.Children[tree.Step])
			tree.Step++
			goto ext
		}
		return tree.Parent.Next(false)
	case PARALLEL:
		break
	case LOOP:
		tree.Step++
		if tree.Step < int(tree.Loop) {
			tree.resetChildren()
			children = append(children, tree.Children[0])
		} else {
			return tree.Parent.Next(false)
		}
	case CONDITION:
		if ret && tree.Step == 0 && len(tree.Children) != 0 {
			tree.Step++
			children = append(children, tree.Children[0])
		} else {
			return tree.Parent.Next(ret)
		}
	default:
		if tree.Step == 0 && len(tree.Children) != 0 {
			tree.Step++
			children = append(children, tree.Children[0])
		} else {
			return tree.Parent.Next(true)
		}
	}

ext:
	return children
}
*/
