package bot

import (
	"fmt"

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

func (b *Bot) Post(p behavior.POST) {
	b.behaviors = append(b.behaviors, &botBehavior{
		behavior: "POST",
		post:     p,
	})
}

func (b *Bot) Jump(j behavior.Jump) {
	b.behaviors = append(b.behaviors, &botBehavior{
		behavior: "Jump",
		jump:     j,
	})
}

func (b *Bot) Delay(d behavior.Delay) {
	b.behaviors = append(b.behaviors, &botBehavior{
		behavior: "Delay",
		delay:    d,
	})
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
