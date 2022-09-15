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

	Loop int32  `xml:"loop"`
	Code string `xml:"code"`

	Step int

	Parent   *Tree
	Children []*Tree `xml:"children"`
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

func (tree *Tree) Link(nod *Tree) {

	tree.Parent = nod
	for k := range tree.Children {
		tree.Children[k].Link(tree)
	}
}

func New(f []byte) (*Tree, error) {

	tree := &Tree{}
	err := xml.Unmarshal([]byte(f), &tree)
	if err != nil {
		panic(err)
	}
	/*
		err = mapstructure.Decode(m, tree)
		if err != nil {
			return nil, fmt.Errorf("behavior tree decode fail %v", err.Error())
		}
	*/
	tree.Parent = nil
	for k := range tree.Children {
		tree.Children[k].Link(tree)
	}

	return tree, nil
}

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

/*
func (tree *Tree) Next() []*Tree {

	var children []*Tree

	if tree.Step < len(tree.Children) {
		nextidx := tree.Step
		tree.Step++
		return tree.Children[nextidx]
	} else {

		if tree.Ty == LOOP {
			if tree.Loop == 0 { // 永远循环
				tree.Step = 0
				tree.resetChildren()
				return tree
			} else {
				tree.LoopStep++
				if tree.LoopStep < tree.Loop {
					tree.Step = 0
					tree.resetChildren()
					return tree
				}
			}
		}

		if tree.Parent != nil {
			return tree.Parent.Next()
		}
	}

	return children
}
*/
