package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pojol/braid-go"
	"github.com/pojol/braid-go/module/meta"
	"github.com/pojol/gobot/driver/bot"
	"github.com/pojol/gobot/driver/bot/behavior"
	"github.com/pojol/gobot/driver/constant"
	"github.com/pojol/gobot/driver/database"
	script "github.com/pojol/gobot/driver/script/module"
	"github.com/pojol/gobot/driver/utils"
)

type BatchInfo struct {
	ID     string
	Name   string
	Cur    int32
	Max    int32
	Errors int32
}

type Batch struct {
	ID           string
	Name         string
	cursorNum    int32
	CurNum       int32
	TotalNum     int32
	BatchNum     int32
	reportTick   int32
	enqueneDelay int32
	Errors       int32

	treeData     []byte
	path         string
	globalScript string

	beginTime time.Time
	bots      map[string]*bot.Bot
	reports   []script.Report

	bwg  utils.SizeWaitGroup
	exit *utils.Switch

	pipeline  chan *bot.Bot
	done      chan interface{}
	BatchDone chan interface{}

	botDoneCh chan string
	botErrCh  chan bot.ErrInfo
}

type BatchConfig struct {
	batchsize     int32
	globalScript  string
	scriptPath    string
	enqeueneDelay int32
}

type BatchReport struct {
	ID      string
	Reports []script.Report
}

func CreateBatch(name, id string, cur, total int32, tbyt []byte, cfg BatchConfig) *Batch {

	b := &Batch{
		ID:           id,
		Name:         name,
		path:         cfg.scriptPath,
		globalScript: cfg.globalScript,
		enqueneDelay: cfg.enqeueneDelay,
		CurNum:       cur,
		BatchNum:     cfg.batchsize,
		TotalNum:     total,
		bwg:          utils.NewSizeWaitGroup(int(cfg.batchsize)),
		exit:         utils.NewSwitch(),
		treeData:     tbyt,
		pipeline:     make(chan *bot.Bot, cfg.batchsize),
		done:         make(chan interface{}, 1),
		BatchDone:    make(chan interface{}, 1),
		botDoneCh:    make(chan string),
		botErrCh:     make(chan bot.ErrInfo),
		beginTime:    time.Now(),

		bots: make(map[string]*bot.Bot),
	}

	fmt.Println("create batch", id, "size", total)
	database.GetTask().New(database.TaskTable{
		ID:          b.ID,
		Name:        name,
		TotalNumber: b.TotalNum,
		CurNumber:   0,
	})

	go b.loop()
	b.run()

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
	b.bots[bot.ID()] = bot
}

func (b *Batch) pop(id string) {
	b.bwg.Done()
	atomic.AddInt32(&b.CurNum, 1)
}

func (b *Batch) loop() {

	for {
		select {
		case botptr := <-b.pipeline:
			b.push(botptr)
			botptr.RunByThread(b.botDoneCh, b.botErrCh)
		case id := <-b.botDoneCh:
			if _, ok := b.bots[id]; ok {
				b.pushReport(b.bots[id])
			}
			b.pop(id)
		case err := <-b.botErrCh:
			if _, ok := b.bots[err.ID]; ok {
				b.pushReport(b.bots[err.ID])
			}
			atomic.AddInt32(&b.Errors, 1)
			b.pop(err.ID)
		case <-b.done:
			goto ext
		}
	}
ext:
	b.exit.Done()
	b.BatchDone <- 1
}

func (b *Batch) run() {

	go func() {

		for {

			if b.exit.HasOpend() {
				fmt.Println("break running")
				break
			}

			var curbatchnum int32
			last := b.TotalNum - atomic.LoadInt32(&b.CurNum)
			if b.BatchNum < last {
				curbatchnum = b.BatchNum
			} else {
				curbatchnum = last
			}

			for i := 0; i < int(curbatchnum); i++ {
				atomic.AddInt32(&b.cursorNum, 1)
				b.bwg.Add()

				b.pipeline <- bot.NewWithBehaviorTree(b.path, behavior.Get(b.Name), behavior.Thread, b.Name, b.ID, atomic.LoadInt32(&b.cursorNum), b.globalScript)
				time.Sleep(time.Millisecond * time.Duration(b.enqueneDelay))
			}

			b.bwg.Wait()
			database.GetTask().Update(b.ID, atomic.LoadInt32(&b.CurNum))
			fmt.Println("batch", b.ID, "end", atomic.LoadInt32(&b.CurNum), "=>", b.TotalNum)
			if atomic.LoadInt32(&b.CurNum) >= b.TotalNum {
				b.done <- 1
			}

			time.Sleep(time.Millisecond * 100)
		}

	}()

}

func (b *Batch) Close() {

}

func (b *Batch) Report() []script.Report {
	return b.reports
}

func (b *Batch) GetBeginTime() time.Time {
	return b.beginTime
}

func (b *Batch) pushReport(bot *bot.Bot) {

	rep := BatchReport{
		ID:      b.ID,
		Reports: bot.GetReport(),
	}
	dat, err := json.Marshal(rep)
	if err != nil {
		fmt.Println("batch.pushReport", err.Error())
		return
	}

	atomic.AddInt32(&b.reportTick, 1)

	if constant.GetClusterState() {
		braid.Topic("batch.report").Pub(context.TODO(), &meta.Message{
			Body: dat,
		})
	} else {
		b.reports = append(b.reports, rep.Reports...) // 暂存下
	}

}
