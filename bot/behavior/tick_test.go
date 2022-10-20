package behavior

import (
	"testing"
	"time"

	"github.com/pojol/gobot/bot/pool"
	"github.com/stretchr/testify/assert"
)

func TestTick(t *testing.T) {

	tree, err := Load([]byte(compose))
	assert.Equal(t, err, nil)

	bb := &Blackboard{
		Nods:      []INod{tree.GetRoot()},
		Threadlst: []ThreadInfo{{Num: 1}},
	}

	tick := &Tick{
		blackboard: bb,
		bs:         pool.GetState(),
	}

	for i := 0; i < 150; i++ {
		tick.Do()
		time.Sleep(time.Millisecond * 50)
	}
}
