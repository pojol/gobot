package bot

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/pojol/apibot/behavior"
	"github.com/pojol/apibot/expression"
	"github.com/pojol/apibot/plugins"
)

type Bot struct {
	url      string
	metadata map[string]interface{}
	tree     *BehaviorTree

	prev *BehaviorTree
	cur  *BehaviorTree

	sync.Mutex

	defaultPost behavior.IPOST
}

type BehaviorTree struct {
	ID   string `mapstructure:"id"`
	Ty   string `mapstructure:"ty"`
	Api  string `mapstructure:"api"`
	Wait int32  `mapstructure:"wait"`

	Loop     int32 `mapstructure:"loop"`
	LoopStep int32

	Parm interface{} `mapstructure:"parm"`
	Expr string      `mapstructure:"expr"`

	Step int

	Parent   *BehaviorTree
	Children []*BehaviorTree `mapstructure:"children"`
}

func (tree *BehaviorTree) link(nod *BehaviorTree) {

	tree.Parent = nod
	for k := range tree.Children {
		tree.Children[k].link(tree)
	}
}

func (b *Bot) GetMetadata() (string, error) {
	byt, err := json.Marshal(b.metadata)
	if err != nil {
		return "", err
	}

	return string(byt), nil
}

func (b *Bot) GetCurNodeID() string {
	if b.cur != nil {
		return b.cur.ID
	}
	return ""
}

func (b *Bot) GetPrevNodeID() string {
	if b.prev != nil {
		return b.prev.ID
	}
	return ""
}

func NewWithBehaviorFile(f []byte, url string) (*Bot, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(f, &m)
	if err != nil {
		return nil, fmt.Errorf("behavior file unmarshal fail %v", err.Error())
	}

	tree := &BehaviorTree{}

	err = mapstructure.Decode(m, tree)
	if err != nil {
		return nil, fmt.Errorf("behavior tree decode fail %v", err.Error())
	}
	tree.Parent = nil
	for k := range tree.Children {
		tree.Children[k].link(tree)
	}

	md := make(map[string]interface{})
	md["Token"] = ""

	return &Bot{
		metadata:    md,
		url:         url,
		tree:        tree,
		cur:         tree,
		defaultPost: &behavior.HTTPPost{URL: url},
	}, nil

}

func (b *Bot) run_selector(nod *BehaviorTree, next bool) (bool, error) {

	if next {
		for k := range nod.Children {
			ok, _ := b.run_nod(nod.Children[k])
			if ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_sequence(nod *BehaviorTree, next bool) (bool, error) {
	if next {
		for k := range nod.Children {
			ok, _ := b.run_nod(nod.Children[k])
			if !ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_condition(nod *BehaviorTree, next bool) (bool, error) {

	eg, err := expression.Parse(nod.Expr)
	if err != nil {
		return false, err
	}

	if eg.DecideWithMap(b.metadata) {
		if next {
			b.run_children(nod, nod.Children)
		}

		return true, nil
	}

	return false, nil
}

func (b *Bot) run_wait(nod *BehaviorTree, next bool) (bool, error) {
	time.Sleep(time.Second * time.Duration(nod.Wait))

	if next {
		b.run_children(nod, nod.Children)
	}

	return true, nil
}

func (b *Bot) run_loop(nod *BehaviorTree, next bool) (bool, error) {

	if nod.Loop == 0 {
		for {
			if next {
				b.run_children(nod, nod.Children)
			}
			time.Sleep(time.Millisecond)
		}
	} else {

		if next {
			for i := 0; i < int(nod.Loop); i++ {
				b.run_children(nod, nod.Children)
				time.Sleep(time.Millisecond)
			}
		}
	}

	return true, nil
}

func (b *Bot) run_http(nod *BehaviorTree, next bool) (bool, error) {

	p := plugins.Get("jsonparse")
	if p == nil {
		return false, fmt.Errorf("can't find serialization plugin %v", "jsonparse")
	}

	byt, err := p.Marshal(nod.Parm)
	if err != nil {
		return false, err
	}

	resBody, err := b.defaultPost.Do(byt, nod.Api)
	if err != nil {
		return false, err
	}
	t := make(map[string]interface{})
	err = p.Unmarshal(resBody, &t)
	if err != nil {
		return false, err
	}

	mergo.MergeWithOverwrite(&b.metadata, t)

	if next {
		b.run_children(nod, nod.Children)
	}

	return true, nil
}

func (b *Bot) run_nod(nod *BehaviorTree) (bool, error) {

	var ok bool
	var err error

	switch nod.Ty {
	case "SelectorNode":
		ok, err = b.run_selector(nod, true)
	case "SequenceNode":
		ok, _ = b.run_sequence(nod, true)
	case "ConditionNode":
		ok, err = b.run_condition(nod, true)
	case "WaitNode":
		ok, _ = b.run_wait(nod, true)
	case "LoopNode":
		ok, err = b.run_loop(nod, true)
	case "HTTPActionNode":
		ok, err = b.run_http(nod, true)
	}

	return ok, err
}

func (b *Bot) run_children(parent *BehaviorTree, children []*BehaviorTree) {
	for k := range children {
		b.run_nod(children[k])
	}
}

func (b *Bot) Run() {
	b.run_children(b.tree, b.tree.Children)
}

func (b *Bot) step(nod *BehaviorTree) bool {

	var ok bool

	switch nod.Ty {
	case "SelectorNode":
		ok, _ = b.run_selector(nod, false)
	case "SequenceNode":
		ok, _ = b.run_sequence(nod, false)
	case "ConditionNode":
		ok, _ = b.run_condition(nod, false)
	case "WaitNode":
		ok, _ = b.run_wait(nod, false)
	case "LoopNode":
		ok, _ = b.run_loop(nod, false)
	case "HTTPActionNode":
		ok, _ = b.run_http(nod, false)
	default:
		ok = true
	}

	fmt.Println("step", ok, nod)

	return ok
}

func (loopnod *BehaviorTree) resetChildren() {
	for k := range loopnod.Children {

		loopnod.Children[k].Step = 0

		if len(loopnod.Children[k].Children) > 0 {
			loopnod.Children[k].resetChildren()
		}

	}
}

func (b *Bot) next(nod *BehaviorTree) *BehaviorTree {
	if nod.Step < len(nod.Children) {
		nextidx := nod.Step
		nod.Step++
		return nod.Children[nextidx]
	} else {

		if nod.Ty == "LoopNode" {
			if nod.Loop == 0 { // 永远循环
				nod.Step = 0
				nod.resetChildren()
				return nod
			} else {
				nod.LoopStep++
				if nod.LoopStep < nod.Loop {
					nod.Step = 0
					nod.resetChildren()
					return nod
				}
			}
		}

		if nod.Parent != nil {
			return b.next(nod.Parent)
		}
	}

	return nil
}

func (b *Bot) RunStep() bool {
	if b.cur == nil {
		return false
	}

	b.Lock()
	defer b.Unlock()

	f := b.step(b.cur)
	// step 中使用了sleep之后，会有多个goroutine执行接下来的程序
	// fmt.Println(goid.Get())

	if f && b.cur.Step < len(b.cur.Children) {
		// down
		nextidx := b.cur.Step
		b.cur.Step++
		next := b.cur.Children[nextidx]
		b.prev = b.cur
		b.cur = next
	} else {
		// right
		if b.cur.Parent != nil {
			b.prev = b.cur
			b.cur = b.next(b.cur.Parent)
			if b.cur == nil {
				return false
			}
		}
	}

	return true
}
