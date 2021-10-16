package factory

import (
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/pojol/gobot-driver/behavior"
	"github.com/pojol/gobot-driver/bot"
	"github.com/pojol/gobot-driver/utils"
)

type BatchInfo struct {
	ID     string
	Name   string
	Cur    int32
	Max    int32
	Errors int32
}

type Batch struct {
	ID       string
	Name     string
	CurNum   int32
	TotalNum int32
	Errors   int32

	bwg     *utils.SizeWaitGroup
	bwgDone chan interface{}

	pipeline chan *bot.Bot
	done     chan interface{}

	botDoneCh chan string
	botErrCh  chan bot.ErrInfo
}

func CreateBatch(scriptPath, name string, num int, tbyt []byte, bwg *utils.SizeWaitGroup, done chan interface{}) *Batch {

	tree, err := behavior.New(tbyt)
	if err != nil {
		return nil
	}

	b := &Batch{
		ID:        uuid.New().String(),
		Name:      name,
		CurNum:    0,
		TotalNum:  int32(num),
		bwg:       bwg,
		bwgDone:   done,
		pipeline:  make(chan *bot.Bot, num),
		done:      make(chan interface{}, 1),
		botDoneCh: make(chan string),
		botErrCh:  make(chan bot.ErrInfo),
	}

	for i := 0; i < num; i++ {
		b.pipeline <- bot.NewWithBehaviorTree(scriptPath, tree, name)
	}

	go b.loop()
	return b
}

func (b *Batch) Info() BatchInfo {
	cur := atomic.LoadInt32(&b.CurNum)

	return BatchInfo{
		ID:     b.ID,
		Name:   b.Name,
		Cur:    cur,
		Max:    b.TotalNum,
		Errors: atomic.LoadInt32(&b.Errors),
	}
}

func (b *Batch) push(bot *bot.Bot) {
	b.bwg.Add()
	atomic.AddInt32(&b.CurNum, 1)
}

func (b *Batch) pop(id string) {
	b.bwg.Done()

	if atomic.LoadInt32(&b.CurNum) == b.TotalNum {
		b.done <- 1
	}
}

func (b *Batch) loop() {
	for {
		select {
		case bot := <-b.pipeline:
			b.push(bot)
			bot.Run(b.botDoneCh, b.botErrCh)
		case id := <-b.botDoneCh:
			b.pop(id)
		case err := <-b.botErrCh:
			atomic.AddInt32(&b.Errors, 1)
			b.pop(err.ID)
		case <-b.done:
			goto ext
		}
	}
ext:
	b.bwgDone <- 1
}
