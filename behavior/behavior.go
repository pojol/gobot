package behavior

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const (
	ROOT       = "RootNode"
	SELETE     = "SelectorNode"
	SEQUENCE   = "SequenceNode"
	CONDITION  = "ConditionNode"
	WAIT       = "WaitNode"
	LOOP       = "LoopNode"
	HTTPACTION = "HTTPActionNode"
)

type Tree struct {
	ID   string `mapstructure:"id"`
	Ty   string `mapstructure:"ty"`
	Api  string `mapstructure:"api"`
	Wait int32  `mapstructure:"wait"`

	Loop     int32 `mapstructure:"loop"`
	LoopStep int32

	Parm interface{} `mapstructure:"parm"`
	Expr string      `mapstructure:"expr"`

	Step int

	Parent   *Tree
	Children []*Tree `mapstructure:"children"`
}

func (tree *Tree) Link(nod *Tree) {

	tree.Parent = nod
	for k := range tree.Children {
		tree.Children[k].Link(tree)
	}
}

func New(f []byte) (*Tree, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(f, &m)
	if err != nil {
		return nil, fmt.Errorf("behavior file unmarshal fail %v", err.Error())
	}

	tree := &Tree{}

	err = mapstructure.Decode(m, tree)
	if err != nil {
		return nil, fmt.Errorf("behavior tree decode fail %v", err.Error())
	}
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
