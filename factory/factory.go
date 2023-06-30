package factory

import (
	"fmt"
	"sync"
	"time"

	"github.com/pojol/gobot/bot"
	"github.com/pojol/gobot/bot/behavior"
	"github.com/pojol/gobot/database"
	"github.com/pojol/gobot/utils"
)

type TaskInfo struct {
	Name string
	Cur  int32
	Num  int32
}

type Factory struct {
	parm Parm

	debugBots map[string]*bot.Bot

	pipelineCache []TaskInfo
	batches       []*Batch
	db            *database.Cache

	batchLock sync.Mutex

	lock sync.Mutex
	exit *utils.Switch
}

func Create(opts ...Option) (*Factory, error) {

	p := Parm{
		lifeTime:    time.Minute,
		Interrupt:   true,
		ReportLimit: 100,
		ScriptPath:  "script/",
		NoDBMode:    false,
	}

	for _, opt := range opts {
		opt(&p)
	}

	db, err := database.Init(p.NoDBMode)
	if err != nil {
		panic(err)
	}
	f := &Factory{
		parm:      p,
		db:        db,
		debugBots: make(map[string]*bot.Bot),
		exit:      utils.NewSwitch(),
	}

	go f.taskLoop()

	fmt.Println("create bot driver", "mode", p.NoDBMode)

	Global = f
	return f, nil
}

var Global *Factory

// Close 关闭机器人工厂
func (f *Factory) Close() {
	f.exit.Open()
}

func (f *Factory) CheckTaskHistory() {
	tasklst, _ := database.GetTask().List()
	for _, task := range tasklst {
		fmt.Println("recover task", task.Name, task.CurNumber, task.TotalNumber)
		f.AddBatch(task.Name, task.CurNumber, task.TotalNumber)

		// 删除旧表
		database.GetTask().Rmv(task.ID)
	}
}

func (f *Factory) AddBatch(name string, cur, total int32) error {

	_, err := database.GetBehavior().Find(name)
	if err != nil {
		return err
	}

	f.pipelineCache = append(f.pipelineCache, TaskInfo{Name: name, Cur: cur, Num: total})
	return nil
}

func (f *Factory) createBatch(name string, cur, num int32) *Batch {

	var dat []byte

	info, err := database.GetBehavior().Find(name)
	if err != nil {
		return nil
	}

	dat = info.File
	cfg, err := database.GetConfig().Get()
	if err != nil {
		return nil
	}

	return CreateBatch(name, cur, num, dat, BatchConfig{
		batchsize:     int32(cfg.ChannelSize),
		globalScript:  string(cfg.GlobalCode),
		scriptPath:    f.parm.ScriptPath,
		enqeueneDelay: int32(cfg.EnqueneDelay),
	})
}

func (f *Factory) CreateDebugBot(name string, fbyt []byte) *bot.Bot {
	var b *bot.Bot

	tree, err := behavior.Load(fbyt, behavior.Step)
	if err != nil {
		return nil
	}

	cfg, err := database.GetConfig().Get()
	if err != nil {
		return nil
	}

	b = bot.NewWithBehaviorTree(f.parm.ScriptPath, tree, name, "", 1, string(cfg.GlobalCode))
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
	delete(f.debugBots, botid)
}

func (f *Factory) taskLoop() {
	for {
		f.lock.Lock()
		if len(f.pipelineCache) > 0 {
			info := f.pipelineCache[0]
			f.pipelineCache = f.pipelineCache[1:]

			b := f.createBatch(info.Name, info.Cur, info.Num)

			f.pushBatch(b)
			<-b.BatchDone
			f.popBatch()
		}
		f.lock.Unlock()

		time.Sleep(time.Millisecond)
	}
}

func (f *Factory) pushBatch(b *Batch) {
	fmt.Println("push batch", b.ID, b.Name)

	f.batchLock.Lock()
	f.batches = append(f.batches, b)
	f.batchLock.Unlock()
}

func (f *Factory) popBatch() {

	f.batchLock.Lock()
	b := f.batches[0]
	database.GetReport().Append(b.Report())
	b.Close()

	fmt.Println("pop batch", b.ID, b.Name)

	s := bot.BotStatusUnknow
	if b.Report().ErrNum > 0 {
		s = bot.BotStatusFail
	} else {
		s = bot.BotStatusSucc
	}
	database.GetBehavior().UpdateStatus(b.Name, s)
	f.batches = f.batches[1:]
	f.batchLock.Unlock()
}

func (f *Factory) GetBatchInfo() []BatchInfo {
	var lst []BatchInfo
	f.batchLock.Lock()
	for _, v := range f.batches {
		lst = append(lst, v.Info())
	}
	f.batchLock.Unlock()
	return lst
}
