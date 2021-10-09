package factory

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
	"github.com/pojol/gobot-driver/behavior"
	"github.com/pojol/gobot-driver/bot"
	"github.com/pojol/gobot-driver/utils"
)

type urlDetail struct {
	ReqNum int
	ErrNum int
	AvgNum int64

	ReqSize int64
	ResSize int64
}

type Report struct {
	ID     string
	Name   string
	BotNum int
	ReqNum int
	ErrNum int
	Tps    int
	Dura   string

	BeginTime time.Time

	UrlMap map[string]*urlDetail
}

type BatchBotInfo struct {
	Behavior string
	Num      int32
}
type BatchInfo struct {
	Batch []BatchBotInfo
}

type Factory struct {
	parm          Parm
	bfile         *BehaviorFile
	reportHistory []*Report

	batchBots map[string]*bot.Bot
	debugBots map[string]*bot.Bot

	pipelineCache []BatchInfo
	running       bool

	translateCh chan *bot.Bot
	doneCh      chan string
	errCh       chan bot.ErrInfo

	batch     utils.SizeWaitGroup
	batchDone chan interface{}

	IncID int64

	colorer *color.Color

	lock sync.Mutex
	exit *utils.Switch
}

func Create(opts ...Option) (*Factory, error) {

	p := Parm{
		frameRate:   time.Second * 1,
		lifeTime:    time.Minute,
		Interrupt:   true,
		batchSize:   1024,
		ReportLimit: 10,
		ScriptPath:  "script/",
	}

	for _, opt := range opts {
		opt(&p)
	}

	f := &Factory{
		parm:        p,
		batchBots:   make(map[string]*bot.Bot),
		debugBots:   make(map[string]*bot.Bot),
		exit:        utils.NewSwitch(),
		translateCh: make(chan *bot.Bot),
		doneCh:      make(chan string),
		errCh:       make(chan bot.ErrInfo),
		batchDone:   make(chan interface{}, 1),
		colorer:     color.New(),
		batch:       utils.New(p.batchSize),
		bfile:       NewBehaviorFile(),
	}

	go f.loop()

	Global = f
	return f, nil
}

var Global *Factory

func (f *Factory) pushReport(rep *Report, bot *bot.Bot) {
	f.lock.Lock()
	defer f.lock.Unlock()

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

// Report 输出报告
func (f *Factory) Report(rep *Report) {

	f.lock.Lock()
	defer f.lock.Unlock()

	fmt.Println("+--------------------------------------------------------------------------------------------------------+")
	fmt.Printf("Req url%-33s Req count %-5s Average time %-5s Body req/res %-5s Succ rate %-10s\n", "", "", "", "", "")

	arr := []string{}
	for k := range rep.UrlMap {
		arr = append(arr, k)
	}
	sort.Strings(arr)

	var reqtotal int64

	for _, sk := range arr {
		v := rep.UrlMap[sk]
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
			f.colorer.Printf("%-40s %-15d %-18s %-18s %-10s\n", u.Path, v.ReqNum, avg, reqsize+" / "+ressize, utils.Red(succ))
		} else {
			fmt.Printf("%-40s %-15d %-18s %-18s %-10s\n", u.Path, v.ReqNum, avg, reqsize+" / "+ressize, succ)
		}
	}
	fmt.Println("+--------------------------------------------------------------------------------------------------------+")

	durations := int(time.Since(rep.BeginTime).Seconds())
	if durations <= 0 {
		durations = 1
	}

	qps := int(reqtotal / int64(durations))
	duration := strconv.Itoa(durations) + "s"

	rep.Tps = qps
	rep.Dura = duration

	if rep.ErrNum != 0 {
		f.colorer.Printf("robot : %d req count : %d duration : %s qps : %d errors : %v\n", rep.BotNum, rep.ReqNum, duration, qps, utils.Red(rep.ErrNum))
	} else {
		fmt.Printf("robot : %d req count : %d duration : %s qps : %d errors : %d\n", rep.BotNum, rep.ReqNum, duration, qps, rep.ErrNum)
	}

}

func (f *Factory) GetReport() []*Report {
	return f.reportHistory
}

// Close 关闭机器人工厂
func (f *Factory) Close() {
	f.exit.Open()
}

func (f *Factory) AddBehavior(name string, byt []byte) {
	err := f.bfile.Upset(name, byt)
	if err != nil {
		fmt.Println("AddBehavior ", err.Error())
	}
}

func (f *Factory) RmvBehavior(name string) {
	err := f.bfile.Del(name)
	if err != nil {
		fmt.Println("RmvBehavior ", err.Error())
	}
}

func (f *Factory) GetBehaviors() []BehaviorInfo {
	lst, err := f.bfile.All()
	if err != nil {
		fmt.Println("GetBehaviors ", err.Error())
	}

	return lst
}

func (f *Factory) FindBehavior(name string) (BehaviorInfo, error) {
	return f.bfile.Find(name)
}

func (f *Factory) Append(info BatchInfo) error {
	for _, v := range info.Batch {
		_, err := f.bfile.Find(v.Behavior)
		if err != nil {
			return err
		}
	}

	f.pipelineCache = append(f.pipelineCache, info)
	return nil
}

func (f *Factory) CreateBot(name string) *bot.Bot {
	var b *bot.Bot

	info, err := f.bfile.Find(name)
	if err != nil {
		return nil
	}

	tree, err := behavior.New(info.Dat)
	if err != nil {
		return nil
	}

	b = bot.NewWithBehaviorTree(f.parm.ScriptPath, tree, name)
	f.batchBots[b.ID()] = b

	return b
}

func (f *Factory) CreateDebugBot(name string) *bot.Bot {
	var b *bot.Bot

	info, err := f.bfile.Find(name)
	if err != nil {
		return nil
	}

	tree, err := behavior.New(info.Dat)
	if err != nil {
		return nil
	}

	b = bot.NewWithBehaviorTree(f.parm.ScriptPath, tree, name)
	f.debugBots[b.ID()] = b

	return b
}

func (f *Factory) FindBot(botid string) *bot.Bot {

	if _, ok := f.debugBots[botid]; ok {
		return f.debugBots[botid]
	}

	return nil
}

func (f *Factory) RmvBot(botid string) {

	if _, ok := f.debugBots[botid]; ok {
		f.debugBots[botid].Close()
		delete(f.debugBots, botid)
	}

}

func (f *Factory) push(bot *bot.Bot) {
	f.batch.Add()

	f.batchBots[bot.ID()] = bot
}

func (f *Factory) pop(id string, err error, rep *Report) {
	f.batch.Done()

	if err != nil && f.parm.Interrupt {
		panic(err)
	}

	if _, ok := f.batchBots[id]; ok {

		f.pushReport(rep, f.batchBots[id])
		if err != nil {
			f.colorer.Printf("%v\n", utils.Red(err.Error()))
		}
		f.batchBots[id].Close()
		delete(f.batchBots, id)

	}

	fmt.Println(len(f.batchBots))
	if len(f.batchBots) == 0 {
		f.batchDone <- 1
	}

}

func (f *Factory) loop() {
	for {
		f.lock.Lock()
		if len(f.pipelineCache) > 0 && !f.running {
			info := f.pipelineCache[0]
			f.pipelineCache = f.pipelineCache[1:]

			fmt.Println("pop", info)

			go f.router()
			f.running = true

			for _, v := range info.Batch {
				for i := 0; i < int(v.Num); i++ {
					f.translateCh <- f.CreateBot(v.Behavior)
				}
			}

		}
		f.lock.Unlock()
		time.Sleep(time.Millisecond)
	}
}

func (f *Factory) router() {

	rep := &Report{
		ID:        strconv.Itoa((time.Now().YearDay() * 1000) + int(atomic.LoadInt64(&f.IncID))),
		BeginTime: time.Now(),
		UrlMap:    make(map[string]*urlDetail),
	}

	for {
		select {
		case bot := <-f.translateCh:
			f.push(bot)
			rep.Name = bot.Name()
			bot.Run(f.exit, f.doneCh, f.errCh)
		case id := <-f.doneCh:
			f.pop(id, nil, rep)
		case err := <-f.errCh:
			f.pop(err.ID, err.Err, rep)
		case <-f.batchDone:
			goto ext
		}
	}
ext:
	atomic.AddInt64(&f.IncID, 1)
	// report
	f.Report(rep)
	if len(f.reportHistory) >= f.parm.ReportLimit {
		f.reportHistory = f.reportHistory[1:]
	}
	f.reportHistory = append(f.reportHistory, rep)
	f.running = false
}
