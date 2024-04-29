package utils

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {

	s := NewSwitch()
	tick := int32(0)

	go func() {
		for {
			<-s.Done()
			atomic.AddInt32(&tick, 1)
		}
	}()

	s.Open()
	assert.Equal(t, s.HasOpend(), true)
	s.Close()
	assert.Equal(t, s.HasOpend(), false)

	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, tick, int32(2))
}
