package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pojol/braid-go"
	"github.com/pojol/braid-go/components"
	"github.com/pojol/braid-go/components/discoverk8s"
	"github.com/pojol/braid-go/module"
	"github.com/pojol/braid-go/module/meta"
	"github.com/pojol/gobot/driver/bot"
	"github.com/pojol/gobot/driver/bot/behavior"
	"github.com/pojol/gobot/driver/constant"
	"github.com/pojol/gobot/driver/database"
	"github.com/pojol/gobot/driver/utils"
	"github.com/redis/go-redis/v9"
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

	bnet      *braid.Braid
	statechan module.IChannel

	debugBots map[string]*bot.Bot

	pipelineCache     []TaskInfo
	batchDoneChannels map[*Batch]chan struct{}

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
		DBType:      "sqlite",
	}

	for _, opt := range opts {
		opt(&p)
	}

	var bnet *braid.Braid
	var statechan module.IChannel

	db, err := database.Init(p.DBType)
	if err != nil {
		panic(err)
	}

	if p.ClusterMode {
		constant.SetClusterState(true)

		bnet, _ = braid.NewService(
			"bot",
			os.Getenv("POD_NAME"),
			&components.DefaultDirector{
				Opts: &components.DirectorOpts{
					RedisCliOpts: &redis.Options{
						Addr: os.Getenv("REDIS_ADDR"),
					},
					DiscoverOpts: []discoverk8s.Option{
						discoverk8s.WithNamespace("bot"),
						discoverk8s.WithSelectorTag("bot"),
					},
				},
			},
		)

		bnet.Init()
		bnet.Run()

		statechan, err = braid.Topic(meta.TopicElectionChangeState).Sub(context.TODO(), "election"+uuid.NewString())
		if err != nil {
			panic(err)
		}

		statechan.Arrived(func(msg *meta.Message) error {
			smsg := meta.DecodeStateChangeMsg(msg)
			constant.SetServerState(smsg.State)
			return nil
		})
	}

	f := &Factory{
		parm:              p,
		bnet:              bnet,
		statechan:         statechan,
		db:                db,
		batchDoneChannels: make(map[*Batch]chan struct{}),
		debugBots:         make(map[string]*bot.Bot),
		exit:              utils.NewSwitch(),
		report:            BuildReport(p.ServiceID),
	}

	go f.taskLoop()
	if constant.GetClusterState() {
		f.watch()
		f.report.watch()
	}

	Global = f
	return f, nil
}

var Global *Factory

// Close 关闭机器人工厂
func (f *Factory) Close() {
	f.report.Close()

	if f.bnet != nil {
		f.bnet.Close()
		f.statechan.Close()
	}

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
			f.batchDoneChannels[b] = b.BatchDone
			f.pushBatch(b)
		}
		f.lock.Unlock()

		for batch, ch := range f.batchDoneChannels {
			select {
			case <-ch:
				f.popBatch(batch.ID)
			default:
				continue
			}
		}

		time.Sleep(time.Millisecond)
	}
}

func (f *Factory) pushBatch(b *Batch) {
	fmt.Println("push batch", b.ID, b.Name)

	f.batchLock.Lock()
	f.batches = append(f.batches, b)
	f.batchLock.Unlock()
}

func (f *Factory) popBatch(id string) {

	f.batchLock.Lock()
	for i, batch := range f.batches {
		if batch.ID == id {

			batch.Close()
			fmt.Println("pop batch", batch.ID, batch.Name)

			if !constant.GetClusterState() {
				rep := Report{
					rep: &database.ReportDetail{
						ID:        batch.ID,
						Name:      batch.Name,
						BeginTime: batch.GetBeginTime().Unix(),
						ApiMap:    make(map[string]*database.ApiDetail),
					},
				}
				rep.Record(batch.TotalNum, batch.Report())
				rep.Generate()
			}

			f.batches = append(f.batches[:i], f.batches[i+1:]...)
			break
		}
	}

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
