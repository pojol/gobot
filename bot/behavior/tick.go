package behavior

import (
	"fmt"

	"github.com/pojol/gobot/bot/state"
)

type Tick struct {
	blackboard *Blackboard
	bs         *state.BotState
	tick       int
}

func (t *Tick) Do() error {

	nods := t.blackboard.GetOpenNods()
	fmt.Println(t.tick)

	for _, n := range nods {

		ns := n.onTick(t)
		if ns == NSFail {
			return n.GetErr()
		}

	}

	t.blackboard.Reset()

	for _, n := range nods {
		n.onNext(t)
	}

	t.tick++
	return nil

}
