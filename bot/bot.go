package bot

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/pojol/apibot/behavior"
	"github.com/pojol/apibot/expression"
	"github.com/pojol/apibot/plugins"
	"github.com/pojol/apibot/utils"
)

type ErrInfo struct {
	ID  string
	Err error
}

type Bot struct {
	id string

	url      string
	metadata map[string]interface{}
	tree     *behavior.Tree

	prev *behavior.Tree
	cur  *behavior.Tree

	sync.Mutex

	post behavior.IPOST
}

func (b *Bot) ID() string {
	return b.id
}

func (b *Bot) GetMetadata() (string, error) {
	byt, err := json.Marshal(b.metadata)
	if err != nil {
		fmt.Println("meta marshal err", err.Error())
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

func NewWithBehaviorTree(bt *behavior.Tree, mockip string) *Bot {

	md := make(map[string]interface{})
	md["Token"] = ""

	bot := &Bot{
		id:       uuid.New().String(),
		metadata: md,
		url:      mockip,
		tree:     bt,
		cur:      bt,
		post:     &behavior.HTTPPost{URL: mockip},
	}

	return bot
}

func (b *Bot) run_selector(nod *behavior.Tree, next bool) (bool, error) {

	if next {
		for k := range nod.Children {
			ok, _ := b.run_nod(nod.Children[k], true)
			if ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_sequence(nod *behavior.Tree, next bool) (bool, error) {
	if next {
		for k := range nod.Children {
			ok, _ := b.run_nod(nod.Children[k], true)
			if !ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_condition(nod *behavior.Tree, next bool) (bool, error) {

	// parse 可以提前到 behavior tree 构造阶段
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

func (b *Bot) run_wait(nod *behavior.Tree, next bool) (bool, error) {
	time.Sleep(time.Second * time.Duration(nod.Wait))

	if next {
		b.run_children(nod, nod.Children)
	}

	return true, nil
}

func (b *Bot) run_loop(nod *behavior.Tree, next bool) (bool, error) {

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

func (b *Bot) parse_right(rv interface{}) interface{} {

	return rv
}

func (b *Bot) run_http(nod *behavior.Tree, next bool) (bool, error) {

	p := plugins.Get("jsonparse")
	if p == nil {
		return false, fmt.Errorf("can't find serialization plugin %v", "jsonparse")
	}

	// nod.Parm need to judge the right value to see if there is a reference value.
	//  like meta.Token

	nod.Parm = b.parse_right(nod.Parm)

	byt, err := p.Marshal(nod.Parm)
	if err != nil {
		return false, fmt.Errorf("marshal plugin err %v", err.Error())
	}

	resBody, err := b.post.Do(byt, nod.Api)
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

func (b *Bot) run_nod(nod *behavior.Tree, next bool) (bool, error) {

	var ok bool
	var err error

	switch nod.Ty {
	case behavior.SELETE:
		ok, err = b.run_selector(nod, next)
	case behavior.SEQUENCE:
		ok, _ = b.run_sequence(nod, next)
	case behavior.CONDITION:
		ok, err = b.run_condition(nod, next)
	case behavior.WAIT:
		ok, _ = b.run_wait(nod, next)
	case behavior.LOOP:
		ok, err = b.run_loop(nod, next)
	case behavior.HTTPACTION:
		ok, err = b.run_http(nod, next)
	case behavior.ROOT:
		ok = true
	default:
		ok = false
		err = fmt.Errorf("run node type err %s", nod.Ty)
	}

	return ok, err
}

func (b *Bot) run_children(parent *behavior.Tree, children []*behavior.Tree) {
	for k := range children {
		b.run_nod(children[k], true)
	}
}

func (b *Bot) Run(sw *utils.Switch, doneCh chan string, errCh chan ErrInfo) {

	go func() {
		b.run_children(b.tree, b.tree.Children)
	}()

}

type State int32

// 系统内部错误
const (
	SEnd State = 1 + iota
	SBreak
	SSucc
)

func (b *Bot) RunStep() State {
	if b.cur == nil {
		return SEnd
	}

	b.Lock()
	defer b.Unlock()

	fmt.Println("step", b.cur.ID, b.cur.Step)
	f, err := b.run_nod(b.cur, false)
	if err != nil {
		b.metadata["err"] = err.Error()
		return SBreak
	}
	// step 中使用了sleep之后，会有多个goroutine执行接下来的程序
	// fmt.Println(goid.Get())

	if b.cur.Parent != nil {
		if b.cur.Parent.Ty == behavior.SELETE && f {
			b.cur.Parent.Step = len(b.cur.Parent.Children)
		}
	}

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
			b.cur = b.cur.Parent.Next()
			if b.cur == nil {
				return SEnd
			}
		}
	}

	return SSucc
}
