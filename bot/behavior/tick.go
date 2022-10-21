package behavior

import (
	"errors"

	"github.com/pojol/gobot/bot/pool"
)

type Tick struct {
	blackboard *Blackboard
	bs         *pool.BotState
}

func NewTick(bb *Blackboard, state *pool.BotState) *Tick {
	t := &Tick{
		blackboard: bb,
		bs:         state,
	}
	return t
}

func (t *Tick) Do() (error, bool) {

	nods := t.blackboard.GetOpenNods()
	t.blackboard.ThreadInfoReset()

	for _, n := range nods {
		n.onTick(t)
		if t.blackboard.HaveErr() {
			return errors.New("thread have err"), false
		}
	}

	t.blackboard.Reset()

	for _, n := range nods {
		n.onNext(t)
	}

	if t.blackboard.end {
		return nil, true
	}

	return nil, false
}
