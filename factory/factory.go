package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pojol/braid-go"
	"github.com/pojol/braid-go/module"
	"github.com/pojol/braid-go/module/meta"
	"github.com/pojol/gobot/bot"
	"github.com/pojol/gobot/bot/behavior"
	"github.com/pojol/gobot/constant"
	"github.com/pojol/gobot/database"
	"github.com/pojol/gobot/utils"
)

type BotBatchInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Num  int    `json:"num"`
	Cnt  int    `json:"cnt"`
}

type TaskInfo struct {
	Name string
	ID   string
	Cur  int32
	Num  int32
}

type Factory struct {
	parm Parm

	debugBots map[string]*bot.Bot

	pipelineCache []TaskInfo

	batches []*Batch
	db      *database.Cache
	report  *FactoryReport

	batchLock sync.Mutex

	createchan module.IChannel

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
		report:    BuildReport(p.ServiceID),
	}

	go f.taskLoop()
	if constant.GetClusterState() {
		f.watch()
		f.report.watch()
	}

	fmt.Println("create bot driver", "mode", p.NoDBMode)

	Global = f
	return f, nil
}

var Global *Factory

// Close 关闭机器人工厂
func (f *Factory) Close() {
	f.report.Close()

	f.createchan.Close()
	f.exit.Open()
}

func (f *Factory) watch() {
	var err error
	f.createchan, err = braid.Topic("bot.batch.create").Sub(context.TODO(), "factory")
	if err != nil {
		fmt.Println("factory.watch", err.Error())
		return
	}

	f.createchan.Arrived(func(msg *meta.Message) error {
		info := &BotBatchInfo{}
		json.Unmarshal(msg.Body, info)

		fmt.Println("pop & add batch", info.ID, info.Name, info.Num)
		f.AddBatch(info.Name, info.ID, 0, int32(info.Num))

		return nil
	})

}

func (f *Factory) CheckTaskHistory() {
	tasklst, _ := database.GetTask().List()
	for _, task := range tasklst {
		fmt.Println("recover task", task.Name, task.ID, task.CurNumber, task.TotalNumber)
		f.AddBatch(task.Name, task.ID, task.CurNumber, task.TotalNumber)

		// 删除旧表
		database.GetTask().Rmv(task.ID)
	}
}

func (f *Factory) AddBatch(name, id string, cur, total int32) error {

	_, err := database.GetBehavior().Find(name)
	if err != nil {
		return err
	}

	f.pipelineCache = append(f.pipelineCache, TaskInfo{ID: id, Name: name, Cur: cur, Num: total})
	return nil
}

func (f *Factory) createBatch(name, id string, cur, num int32) *Batch {

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

	return CreateBatch(name, id, cur, num, dat, BatchConfig{
		batchsize:     int32(cfg.ChannelSize),
		globalScript:  string(cfg.GlobalCode),
		scriptPath:    f.parm.ScriptPath,
		enqeueneDelay: int32(cfg.EnqueneDelay),
	})
}

func (f *Factory) CreateDebugBot(name string, fbyt []byte) *bot.Bot {
	var b *bot.Bot

	tree, err := behavior.Load(fbyt)
	if err != nil {
		return nil
	}

	cfg, err := database.GetConfig().Get()
	if err != nil {
		return nil
	}

	b = bot.NewWithBehaviorTree(f.parm.ScriptPath, tree, behavior.Step, name, "", 1, string(cfg.GlobalCode))
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

			b := f.createBatch(info.Name, info.ID, info.Cur, info.Num)

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
	b.Close()

	fmt.Println("pop batch", b.ID, b.Name)
	if !constant.GetClusterState() {
		rep := Report{
			rep: &database.ReportDetail{
				ID:        b.ID,
				Name:      b.Name,
				BeginTime: b.GetBeginTime().Unix(),
				ApiMap:    make(map[string]*database.ApiDetail),
			},
		}
		rep.Record(b.TotalNum, b.Report())
		rep.Generate()
	}

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
