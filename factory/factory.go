package factory

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/pojol/apibot/behavior"
	"github.com/pojol/apibot/bot"
	"github.com/pojol/apibot/utils"
)

type BehaviorInfo struct {
	Name       string
	RootID     string
	Dat        []byte
	UpdateTime int64
}

type BatchBotInfo struct {
	Behavior string
	Num      int32
}
type BatchInfo struct {
	Batch []BatchBotInfo
}

type Factory struct {
	parm      Parm
	batchBots map[string]*bot.Bot
	debugBots map[string]*bot.Bot

	behaviorLst []BehaviorInfo

	pipelineCache []BatchInfo
	running       bool

	translateCh chan *bot.Bot
	doneCh      chan string
	errCh       chan bot.ErrInfo

	batch     utils.SizeWaitGroup
	batchDone chan interface{}

	colorer *color.Color

	lock sync.Mutex
	exit *utils.Switch
}

func Create(opts ...Option) (*Factory, error) {

	p := Parm{
		frameRate:  time.Second * 1,
		lifeTime:   time.Minute,
		Interrupt:  true,
		batchSize:  1024,
		ScriptPath: "script/",
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
	}

	go f.loop()

	Global = f
	return f, nil
}

var Global *Factory

// Close 关闭机器人工厂
func (f *Factory) Close() {
	f.exit.Open()
}

func (f *Factory) AddBehavior(rootid string, name string, byt []byte) {

	flag := false
	idx := 0
	for k, v := range f.behaviorLst {
		if v.Name == name {
			flag = true
			idx = k
			break
		}
	}

	if flag {
		f.behaviorLst = append(f.behaviorLst[:idx], f.behaviorLst[idx+1:]...)
	}

	f.behaviorLst = append(f.behaviorLst, BehaviorInfo{
		Name:       name,
		RootID:     rootid,
		Dat:        byt,
		UpdateTime: time.Now().Unix(),
	})
}

func (f *Factory) GetBehaviors() []BehaviorInfo {
	info := []BehaviorInfo{}
	for _, v := range f.behaviorLst {
		info = append(info, BehaviorInfo{
			Name:       v.Name,
			UpdateTime: v.UpdateTime,
		})
	}

	return info
}

func (f *Factory) FindBehavior(name string) (BehaviorInfo, error) {

	var info BehaviorInfo
	err := fmt.Errorf("FindBehavior err not found %v", name)

	for _, v := range f.behaviorLst {
		if v.Name == name {
			info = v
			err = nil
			break
		}
	}

	return info, err

}

func (f *Factory) Append(info BatchInfo) {
	f.pipelineCache = append(f.pipelineCache, info)
}

func (f *Factory) CreateBot(name string) *bot.Bot {
	var b *bot.Bot

	for _, v := range f.behaviorLst {
		if v.Name == name {

			tree, err := behavior.New(v.Dat)
			if err != nil {
				return nil
			}

			b = bot.NewWithBehaviorTree(f.parm.ScriptPath, tree)
			f.batchBots[b.ID()] = b
			break
		}
	}

	return b
}

func (f *Factory) CreateDebugBot(name string) *bot.Bot {
	var b *bot.Bot

	for _, v := range f.behaviorLst {
		if v.Name == name {

			tree, err := behavior.New(v.Dat)
			if err != nil {
				return nil
			}

			b = bot.NewWithBehaviorTree(f.parm.ScriptPath, tree)
			f.debugBots[b.ID()] = b
			break
		}
	}

	return b
}

func (f *Factory) FindBot(botid string) *bot.Bot {

	if _, ok := f.debugBots[botid]; ok {
		return f.debugBots[botid]
	}

	return nil
}

func (f *Factory) RmvBot(botid string) {
	delete(f.debugBots, botid)

	/*
		if f.parm.mock != nil {
			f.parm.mock.Reset(botid)
		}
	*/
}

func (f *Factory) push(bot *bot.Bot) {
	f.batch.Add()

	f.batchBots[bot.ID()] = bot
}

func (f *Factory) pop(id string, err error) {
	f.batch.Done()

	if err != nil && f.parm.Interrupt {
		panic(err)
	}

	if _, ok := f.batchBots[id]; ok {

		//f.pushReport(f.bots[id])
		if err != nil {
			f.colorer.Printf("%v\n", utils.Red(err.Error()))
		}
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

	for {
		select {
		case bot := <-f.translateCh:
			f.push(bot)
			bot.Run(f.exit, f.doneCh, f.errCh)
		case id := <-f.doneCh:
			f.pop(id, nil)
		case err := <-f.errCh:
			f.pop(err.ID, err.Err)
		case <-f.batchDone:
			goto ext
		}
	}

ext:

	// report
	fmt.Println("run done")
	f.running = false
}
