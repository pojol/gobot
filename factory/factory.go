package factory

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pojol/gobot/bot"
	"github.com/pojol/gobot/bot/behavior"
	"github.com/pojol/gobot/database"
	"github.com/pojol/gobot/utils"
)

type TaskInfo struct {
	Name string
	Num  int32
}

type Factory struct {
	parm   Parm
	report *Report

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
		frameRate:   time.Second * 1,
		lifeTime:    time.Minute,
		Interrupt:   true,
		ReportLimit: 100,
		ScriptPath:  "script/",
		NoDBMode:    false,
	}

	for _, opt := range opts {
		opt(&p)
	}

	db := database.Create()
	f := &Factory{
		parm:      p,
		db:        db,
		debugBots: make(map[string]*bot.Bot),
		exit:      utils.NewSwitch(),
		report:    NewReport(int32(p.ReportLimit), db),
	}

	go f.taskLoop()

	Global = f
	return f, nil
}

var Global *Factory

func (f *Factory) GetReport() []database.ReportInfo {
	return f.report.Info()
}

// Close 关闭机器人工厂
func (f *Factory) Close() {
	f.exit.Open()
}

func (f *Factory) AddBehavior(name string, byt []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()

	err := f.db.UpsetFile(name, byt)
	if err != nil {
		fmt.Println("AddBehavior ", err.Error())
	}

}

func (f *Factory) RmvBehavior(name string) {
	err := f.db.DelFile(name)
	if err != nil {
		fmt.Println("RmvBehavior ", err.Error())
	}
}

func (f *Factory) UpdateBehaviorTags(name string, tags []byte) []database.BehaviorInfo {
	err := f.db.UpdateTags(name, tags)
	if err != nil {
		fmt.Println("UpdateBehaviorTags", err.Error())
	}

	return f.GetBehaviors()
}

func (f *Factory) GetBehaviors() []database.BehaviorInfo {
	lst, err := f.db.GetAllFiles()
	if err != nil {
		fmt.Println("GetBehaviors ", err.Error())
	}

	return lst
}

func (f *Factory) UploadConfig(name string, dat []byte) error {

	if name == "" {
		return errors.New("upload config err : meaningless naming")
	}

	_name := strings.ToLower(name)

	return f.db.ConfigUpset(_name, dat)
}

func (f *Factory) GetConfig(name string) (database.TemplateConfig, error) {

	_name := strings.ToLower(name)

	return f.db.ConfigFind(_name)
}

func (f *Factory) RemoveConfig(name string) error {

	_name := strings.ToLower(name)

	return f.db.ConfigRemove(_name)
}

func (f *Factory) GetConfigList() ([]string, error) {
	return f.db.ConfigList()
}

func (f *Factory) FindBehavior(name string) (database.BehaviorInfo, error) {
	return f.db.FindFile(name)
}

func (f *Factory) AddTask(name string, cnt int32) error {
	_, err := f.db.FindFile(name)
	if err != nil {
		return err
	}

	f.pipelineCache = append(f.pipelineCache, TaskInfo{Name: name, Num: cnt})
	return nil
}

func (f *Factory) CreateTask(name string, num int) *Batch {

	var dat []byte

	info, err := f.db.FindFile(name)
	if err != nil {
		return nil
	}

	dat = info.Dat

	return CreateBatch(f.parm.ScriptPath,
		name,
		num,
		dat,
		int32(config.GetChannelSize()),
		string(config.GetGlobalDefine()))
}

func (f *Factory) CreateDebugBot(name string, fbyt []byte) *bot.Bot {
	var b *bot.Bot

	tree, err := behavior.Load(fbyt, behavior.Step)
	if err != nil {
		return nil
	}

	b = bot.NewWithBehaviorTree(f.parm.ScriptPath, tree, name, 1, string(config.GetGlobalDefine()))
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

			b := f.CreateTask(info.Name, int(info.Num))

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
	f.report.Append(b.Report())
	b.Close()

	fmt.Println("pop batch", b.ID, b.Name)

	s := bot.BotStatusUnknow
	if b.Report().ErrNum > 0 {
		s = bot.BotStatusFail
	} else {
		s = bot.BotStatusSucc
	}
	f.db.UpdateState(b.Name, s)
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
