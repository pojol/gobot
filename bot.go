package bot

import (
	"fmt"
	"reflect"

	"github.com/pojol/apibot/behavior"
)

type botBehavior struct {
	behavior string

	post  behavior.POST
	jump  behavior.Jump
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

func (b *Bot) NewBehavor(behav interface{}) {

	name := reflect.TypeOf(behav).Name()
	switch name {
	case "POST":
		p, ok := behav.(behavior.POST)
		if ok {
			b.behaviors = append(b.behaviors, &botBehavior{
				behavior: name,
				post:     p,
			})
		}

	case "Jump":
		b.behaviors = append(b.behaviors, &botBehavior{
			behavior: name,
			jump:     behav.(behavior.Jump),
		})
	case "Delay":
		b.behaviors = append(b.behaviors, &botBehavior{
			behavior: name,
			delay:    behav.(behavior.Delay),
		})
	}
}

func (b *Bot) Run() {
	var err error

	for _, v := range b.behaviors {
		if v.behavior == "POST" {
			err = v.post.Exec()
			if err != nil {
				goto ext
			}
		}
	}

ext:
	if err != nil {
		fmt.Println(err.Error())
	}
}
