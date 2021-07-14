package bot

import (
	"fmt"

	"github.com/pojol/apibot/behavior"
)

type botBehavior struct {
	behavior behavior.BehaviorType

	post  behavior.IPOST
	delay behavior.IDelay
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

func (b *Bot) Post(p behavior.IPOST) {
	b.behaviors = append(b.behaviors, &botBehavior{
		behavior: behavior.PostTy,
		post:     p,
	})
}

func (b *Bot) Delay(d behavior.IDelay) {
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
