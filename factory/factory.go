package factory

import (
	"errors"
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
	parm       Parm
	pickCursor int
	bots       map[string]*bot.Bot

	behaviorLst []BehaviorInfo
	pipeline    []BatchInfo

	translateCh chan *bot.Bot
	doneCh      chan string
	errCh       chan bot.ErrInfo

	batch utils.SizeWaitGroup

	colorer *color.Color

	lock sync.Mutex
	wg   sync.WaitGroup
	exit *utils.Switch
}

func Create(opts ...Option) (*Factory, error) {

	p := Parm{
		frameRate: time.Second * 1,
		lifeTime:  time.Minute,
		Interrupt: true,
		batchSize: 1024,
	}

	for _, opt := range opts {
		opt(&p)
	}

	f := &Factory{
		parm:        p,
		bots:        make(map[string]*bot.Bot),
		exit:        utils.NewSwitch(),
		translateCh: make(chan *bot.Bot),
		doneCh:      make(chan string),
		errCh:       make(chan bot.ErrInfo),
		colorer:     color.New(),
		batch:       utils.New(p.batchSize),
	}

	f.wg.Add(1)

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
	f.pipeline = append(f.pipeline, info)
}

func (f *Factory) getRobot() *bot.Bot {

	if len(f.behaviorLst) <= 0 {
		panic(errors.New("no behavior tree list"))
	}

	if f.pickCursor >= len(f.behaviorLst) {
		f.pickCursor = 0
	}

	info := f.behaviorLst[f.pickCursor]
	f.pickCursor++

	tree, err := behavior.New(info.Dat)
	if err != nil {
		return nil
	}

	return bot.NewWithBehaviorTree(tree)
}

func (f *Factory) CreateBot(name string) *bot.Bot {
	var b *bot.Bot

	for _, v := range f.behaviorLst {
		if v.Name == name {

			tree, err := behavior.New(v.Dat)
			if err != nil {
				return nil
			}

			b = bot.NewWithBehaviorTree(tree)
			f.bots[b.ID()] = b
			break
		}
	}

	return b
}

func (f *Factory) FindBot(botid string) *bot.Bot {

	if _, ok := f.bots[botid]; ok {
		return f.bots[botid]
	}

	return nil
}

func (f *Factory) RmvBot(botid string) {
	delete(f.bots, botid)

	/*
		if f.parm.mock != nil {
			f.parm.mock.Reset(botid)
		}
	*/
}

// Run 运行
func (f *Factory) RunBatch() error {

	go f.router()

	if f.parm.tickCreateNum == 0 {
		f.parm.tickCreateNum = len(f.behaviorLst)
	}

	if f.parm.mode == FactoryModeStatic {
		f.static()
	} else if f.parm.mode == FactoryModeIncrease {
		f.increase()
		time.AfterFunc(f.parm.lifeTime, func() {
			f.exit.Open()
		})
	}

	<-f.exit.Done()
	f.wg.Wait()

	return nil
}

func (f *Factory) push(bot *bot.Bot) {
	f.batch.Add()

	f.bots[bot.ID()] = bot
}

func (f *Factory) pop(id string, err error) {
	f.batch.Done()

	if err != nil && f.parm.Interrupt {
		panic(err)
	}

	if _, ok := f.bots[id]; ok {

		//f.pushReport(f.bots[id])
		if err != nil {
			f.colorer.Printf("%v\n", utils.Red(err.Error()))
		}
		delete(f.bots, id)

	}

	if len(f.bots) == 0 && f.parm.mode == FactoryModeStatic {
		f.exit.Open()
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
		case <-f.exit.Done():
			goto ext
		}
	}

ext:

	// report
	f.wg.Done()
}

func (f *Factory) static() {

	for i := 0; i < f.parm.tickCreateNum; i++ {
		f.translateCh <- f.getRobot()
	}

	f.batch.Wait()

}

func (f *Factory) increase() {

	go func() {

		ticker := time.NewTicker(f.parm.frameRate)

		for {
			select {
			case <-ticker.C:

				if f.exit.HasOpend() {
					break
				}

				f.static()

			case <-f.exit.Done():
			}
		}

	}()

}
