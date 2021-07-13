package bot

import (
	"fmt"

	"github.com/pojol/apibot/behavior"
)

type botBehavior struct {
	behavior behavior.BehaviorType

	post  behavior.POST
	delay behavior.Delay
}

type Bot struct {
	name      string
	metadata  interface{}
	behaviors []*botBehavior
}

func New(name string, meta interface{}) *Bot {

	return &Bot{
		name:     name,
		metadata: meta,
	}

}

func (b *Bot) Post(p behavior.POST) {
	b.behaviors = append(b.behaviors, &botBehavior{
		behavior: behavior.PostTy,
		post:     p,
	})
}

func (b *Bot) Delay(d behavior.Delay) {
	b.behaviors = append(b.behaviors, &botBehavior{
		behavior: behavior.DelayTy,
		delay:    d,
	})
}

func (b *Bot) Run() {
	var err error

	end := len(b.behaviors)
	begin := 0

	for {
		b := b.behaviors[begin]
		if b.behavior == behavior.PostTy {
			err = b.post.Do()
			if err != nil {
				goto ext
			}
		}

		begin++
		if begin == end {
			break
		}
	}

ext:
	if err != nil {
		fmt.Println(err.Error())
	}
}
