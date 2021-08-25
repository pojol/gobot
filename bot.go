package bot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/pojol/apibot/behavior"
	"github.com/pojol/apibot/expression"
	"github.com/pojol/apibot/plugins"
)

/*
	{
		"$ne" : {
			"meta.token" : ""
		}
	}
*/

type Bot struct {
	name string

	url      string
	metadata map[string]interface{}
	tree     *BehaviorTree

	defaultPost behavior.IPOST
}

type BehaviorTree struct {
	Ty   string      `mapstructure:"ty"`
	Api  string      `mapstructure:"api"`
	Loop int32       `mapstructure:"loop"`
	Parm interface{} `mapstructure:"parm"`
	Expr string      `mapstructure:"expr"`

	Children []BehaviorTree `mapstructure:"children"`
}

func NewWithBehaviorFile(f []byte, url string, meta interface{}) (*Bot, error) {
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

	md := make(map[string]interface{})
	md["Token"] = ""

	return &Bot{
		metadata:    md,
		url:         url,
		tree:        tree,
		defaultPost: &behavior.HTTPPost{URL: url},
	}, nil

}

func (b *Bot) run_selector(nod *BehaviorTree) (bool, error) {

	for k := range nod.Children {
		ok, _ := b.run_nod(&nod.Children[k])
		if ok {
			break
		}
	}

	return true, nil
}

func (b *Bot) run_condition(nod *BehaviorTree) (bool, error) {

	eg, err := expression.Parse(nod.Expr)
	if err != nil {
		return false, err
	}

	if eg.DecideWithMap(b.metadata) {
		b.run_children(nod.Children)
		return true, nil
	}

	return false, nil
}

func (b *Bot) run_loop(nod *BehaviorTree) (bool, error) {

	if nod.Loop == 0 {
		for {
			b.run_children(nod.Children)
			time.Sleep(time.Millisecond)
		}
	} else {
		for i := 0; i < int(nod.Loop); i++ {
			b.run_children(nod.Children)
			time.Sleep(time.Millisecond)
		}
	}

	return true, nil
}

func (b *Bot) run_http(nod *BehaviorTree) (bool, error) {

	p := plugins.Get("jsonp")
	if p == nil {
		return false, fmt.Errorf("can't find serialization plugin %v", "jsonp")
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

	b.run_children(nod.Children)
	return true, nil
}

func (b *Bot) run_nod(nod *BehaviorTree) (bool, error) {

	var ok bool
	var err error

	switch nod.Ty {
	case "SelectorNode":
		ok, err = b.run_selector(nod)
	case "ConditionNode":
		ok, err = b.run_condition(nod)
	case "LoopNode":
		ok, err = b.run_loop(nod)
	case "HTTPActionNode":
		ok, err = b.run_http(nod)
	}

	return ok, err
}

func (b *Bot) run_children(children []BehaviorTree) {
	for k := range children {
		b.run_nod(&children[k])
	}
}

func (b *Bot) Run() {
	b.run_children(b.tree.Children)

	fmt.Println(b.metadata)
}
