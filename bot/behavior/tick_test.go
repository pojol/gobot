package behavior

import (
	"testing"

	"github.com/pojol/gobot/bot/state"
	"github.com/stretchr/testify/assert"
)

func TestTick(t *testing.T) {

	tree, err := Load([]byte(compose))
	assert.Equal(t, err, nil)

	bb := &Blackboard{}
	bb.Append([]INod{tree.root})

	tick := &Tick{
		blackboard: bb,
		bs:         state.GetState(),
	}

	for i := 0; i < 10; i++ {
		tick.Do()
	}

}
