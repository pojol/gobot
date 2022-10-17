package factory

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/pojol/gobot/bot"
	"github.com/pojol/gobot/bot/behavior"
	"github.com/pojol/gobot/utils"
)

type BatchInfo struct {
	ID     string
	Name   string
	Cur    int32
	Max    int32
	Errors int32
}

type Batch struct {
	ID        string
	Name      string
	cursorNum int32
	CurNum    int32
	TotalNum  int32
	BatchNum  int32
	Errors    int32

	tree         *behavior.Tree
	path         string
	globalScript []string

	bots    map[string]*bot.Bot
	colorer *color.Color
	rep     *ReportDetail

	bwg  utils.SizeWaitGroup
	exit *utils.Switch

	pipeline  chan *bot.Bot
	done      chan interface{}
	BatchDone chan interface{}

	botDoneCh chan string
	botErrCh  chan bot.ErrInfo
}

func CreateBatch(scriptPath, name string, num int, tbyt []byte, batchsize int32, globalScript []string) *Batch {

	tree, err := behavior.Load(tbyt)
	if err != nil {
		return nil
	}

	b := &Batch{
		ID:           uuid.New().String(),
		Name:         name,
		path:         scriptPath,
		globalScript: globalScript,
		CurNum:       0,
		BatchNum:     batchsize,
		TotalNum:     int32(num),
		bwg:          utils.NewSizeWaitGroup(int(batchsize)),
		exit:         utils.NewSwitch(),
		tree:         tree,
		pipeline:     make(chan *bot.Bot, num),
		done:         make(chan interface{}, 1),
		BatchDone:    make(chan interface{}, 1),
		botDoneCh:    make(chan string),
		botErrCh:     make(chan bot.ErrInfo),

		colorer: color.New(),
		bots:    make(map[string]*bot.Bot),
	}

	fmt.Println("create", num, "bot", "pipeline size", batchsize)
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

func (b *Batch) Report() ReportDetail {
	return *b.rep
}

func (b *Batch) push(bot *bot.Bot) {
	b.bwg.Add()
	fmt.Println("bot", bot.ID(), "running", atomic.LoadInt32(&b.cursorNum), "=>", b.TotalNum)

	b.bots[bot.ID()] = bot
}

func (b *Batch) pop(id string) {
	b.bwg.Done()
	atomic.AddInt32(&b.CurNum, 1)

	if atomic.LoadInt32(&b.CurNum) >= b.TotalNum {
		b.done <- 1
	}
}

func (b *Batch) loop() {

	b.rep = &ReportDetail{
		ID:        b.ID,
		Name:      b.Name,
		BeginTime: time.Now(),
		UrlMap:    make(map[string]*urlDetail),
	}

	for {
		select {
		case botptr := <-b.pipeline:
			b.push(botptr)
			botptr.Run(b.botDoneCh, b.botErrCh, bot.Batch)
		case id := <-b.botDoneCh:
			if _, ok := b.bots[id]; ok {
				b.pushReport(b.rep, b.bots[id])
			}
			b.pop(id)
		case err := <-b.botErrCh:
			if _, ok := b.bots[err.ID]; ok {
				b.pushReport(b.rep, b.bots[err.ID])
			}
			atomic.AddInt32(&b.Errors, 1)
			b.pop(err.ID)
		case <-b.done:
			goto ext
		}
	}
ext:
	b.record()
	b.exit.Done()
	b.BatchDone <- 1
}

func (b *Batch) run() {

	go func() {

		for {

			if b.exit.HasOpend() {
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
				b.pipeline <- bot.NewWithBehaviorTree(b.path, b.tree, b.Name, atomic.LoadInt32(&b.cursorNum), b.globalScript)
			}

			b.bwg.Wait()
			time.Sleep(time.Millisecond * 100)
		}

	}()

}

func (b *Batch) Close() {

}

func (b *Batch) pushReport(rep *ReportDetail, bot *bot.Bot) {
	rep.BotNum++
	robotReport := bot.GetReport()

	rep.ReqNum += len(robotReport)
	for _, v := range robotReport {
		if _, ok := rep.UrlMap[v.Api]; !ok {
			rep.UrlMap[v.Api] = &urlDetail{}
		}

		rep.UrlMap[v.Api].ReqNum++
		rep.UrlMap[v.Api].AvgNum += int64(v.Consume)
		rep.UrlMap[v.Api].ReqSize += int64(v.ReqBody)
		rep.UrlMap[v.Api].ResSize += int64(v.ResBody)
		if v.Err != "" {
			rep.ErrNum++
			rep.UrlMap[v.Api].ErrNum++
		}
	}

}

func (b *Batch) record() {

	fmt.Println("+--------------------------------------------------------------------------------------------------------+")
	fmt.Printf("Req url%-33s Req count %-5s Average time %-5s Body req/res %-5s Succ rate %-10s\n", "", "", "", "", "")

	arr := []string{}
	for k := range b.rep.UrlMap {
		arr = append(arr, k)
	}
	sort.Strings(arr)

	var reqtotal int64

	for _, sk := range arr {
		v := b.rep.UrlMap[sk]
		var avg string
		if v.AvgNum == 0 {
			avg = "0 ms"
		} else {
			avg = strconv.Itoa(int(v.AvgNum/int64(v.ReqNum))) + "ms"
		}

		succ := strconv.Itoa(v.ReqNum-v.ErrNum) + "/" + strconv.Itoa(v.ReqNum)

		reqsize := strconv.Itoa(int(v.ReqSize/1024)) + "kb"
		ressize := strconv.Itoa(int(v.ResSize/1024)) + "kb"

		reqtotal += int64(v.ReqNum)

		u, _ := url.Parse(sk)
		if v.ErrNum != 0 {
			b.colorer.Printf("%-40s %-15d %-18s %-18s %-10s\n", u.Path, v.ReqNum, avg, reqsize+" / "+ressize, utils.Red(succ))
		} else {
			fmt.Printf("%-40s %-15d %-18s %-18s %-10s\n", u.Path, v.ReqNum, avg, reqsize+" / "+ressize, succ)
		}
	}
	fmt.Println("+--------------------------------------------------------------------------------------------------------+")

	durations := int(time.Since(b.rep.BeginTime).Seconds())
	if durations <= 0 {
		durations = 1
	}

	qps := int(reqtotal / int64(durations))
	duration := strconv.Itoa(durations) + "s"

	b.rep.Tps = qps
	b.rep.Dura = duration

	if b.rep.ErrNum != 0 {
		b.colorer.Printf("robot : %d match to %d APIs req count : %d duration : %s qps : %d errors : %v\n", b.rep.BotNum, len(b.rep.UrlMap), b.rep.ReqNum, duration, qps, utils.Red(b.rep.ErrNum))
	} else {
		fmt.Printf("robot : %d match to %d APIs req count : %d duration : %s qps : %d errors : %d\n", b.rep.BotNum, len(b.rep.UrlMap), b.rep.ReqNum, duration, qps, b.rep.ErrNum)
	}

}
