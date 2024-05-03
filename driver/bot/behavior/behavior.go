package behavior

import (
	"encoding/xml"
)

type Mode int

const (
	Thread Mode = 1 + iota // 线程运行模式（一般用于压测
	Block                  // 阻塞运行模式（一般用于http调用
	Step                   // 步进运行模式（一般用于调试
)

type Tree struct {
	ID string `xml:"id"`
	Ty string `xml:"ty"`

	Wait int32 `xml:"wait"`

	Loop int32  `xml:"loop"` // 用于记录循环节点的循环x次数
	Code string `xml:"code"`

	root INod

	Children []*Tree `xml:"children"`
}

func (t *Tree) GetRoot() INod {
	return t.root
}

func (t *Tree) link(self INod, parent INod) {

	self.Init(t, parent)

	for k := range t.Children {
		child := NewNode(t.Children[k].Ty).(INod)
		self.AddChild(child)
		t.Children[k].link(child, self)
	}

}

func Load(f []byte) (*Tree, error) {

	tree := &Tree{}

	err := xml.Unmarshal([]byte(f), &tree)
	if err != nil {
		panic(err)
	}

	tree.root = NewNode(tree.Ty).(INod)
	tree.root.Init(tree, nil)

	for k := range tree.Children {

		cn := NewNode(tree.Children[k].Ty).(INod)
		cn.Init(tree.Children[k], tree.root)
		tree.root.AddChild(cn)

		tree.Children[k].link(cn, tree.root)
	}

	return tree, nil
}

func (t *Tree) Reset() {
	t.root.onReset()
}
