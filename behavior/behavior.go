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

	Loop     int32 `xml:"loop"`
	LoopStep int32

	Code string `xml:"code"`

	Step int

	Parent   *Tree
	Children []*Tree `xml:"children"`
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

// 这边的重置，应该是重置loop节点名下的所有children
func (tree *Tree) resetChildren() {
	for k := range tree.Children {

		tree.Children[k].Step = 0
		if tree.Children[k].Ty == LOOP {
			tree.Children[k].LoopStep = 0
		}

		if len(tree.Children[k].Children) > 0 {
			tree.Children[k].resetChildren()
		}

	}
}

func (tree *Tree) Next() *Tree {

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

	return nil
}
