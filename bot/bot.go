package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pojol/gobot/bot/behavior"
	"github.com/pojol/gobot/bot/pool"
	script "github.com/pojol/gobot/script/module"
	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
)

type ErrInfo struct {
	ID  string
	Err error
}

type Bot struct {
	id   string
	name string

	preloadErr string

	bb   *behavior.Blackboard
	tick *behavior.Tick
	bt   *behavior.Tree

	sync.RWMutex
	bs *pool.BotState // lua state pool

}

const (
	BotStatusSucc   = "succ"
	BotStatusFail   = "fail"
	BotStatusUnknow = "unknow"
)

func (b *Bot) ID() string {
	return b.id
}

func (b *Bot) Name() string {
	return b.name
}

func (b *Bot) GetMetaInfo() string {
	tablemap := make(map[string]interface{})
	table, ok := b.bs.L.GetGlobal("meta").(*lua.LTable)

	var tableerr error
	var byt []byte

	if b.preloadErr != "" {
		byt, _ = json.Marshal(&tablemap)
		goto ext
	}

	if ok {
		tablemap, tableerr = utils.Table2Map(table)
		if tableerr != nil {
			tablemap["err"] = tableerr.Error()
			goto ext
		}

		byt, tableerr = json.Marshal(&tablemap)
		if tableerr != nil {
			tablemap["err"] = tableerr.Error()
			byt, _ = json.Marshal(&tablemap)
			goto ext
		}

	} else {

		tablemap["err"] = errors.New("the meta field is not obtained")
		byt, _ = json.Marshal(&tablemap)
		goto ext

	}

ext:
	return string(byt)
}

func (b *Bot) GetThreadInfo() string {

	lst := b.bb.ThreadInfo()

	info, err := json.Marshal(&lst)
	if err != nil {
		fmt.Println(err.Error())
	}

	return string(info)
}

func NewWithBehaviorTree(path string, bt *behavior.Tree, name, batch string, idx int32, globalScript string) *Bot {

	bb := &behavior.Blackboard{
		Nods:      []behavior.INod{bt.GetRoot()},
		Threadlst: []behavior.ThreadInfo{{Number: 1}},
	}

	var state *pool.BotState
	var id string

	if bt.GetMode() == behavior.Thread {
		state = pool.GetState()
		id = strconv.Itoa(int(idx))
	} else {
		state = pool.NewState()
		id = uuid.NewString()
	}

	if batch == "" {
		batch = "-"
	}

	bot := &Bot{
		id:   id,
		bb:   bb,
		tick: behavior.NewTick(bb, state, strconv.Itoa(int(idx))),
		bt:   bt,
		bs:   state,
		name: name,
	}

	rand.Seed(time.Now().UnixNano())

	// 加载预定义全局脚本文件
	if globalScript != "" {
		pool.DoString(bot.bs.L, globalScript)
	}

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		err := pool.DoFile(bot.bs.L, path+v)
		if err != nil {
			fmt.Println("err", err.Error())
			bot.preloadErr = fmt.Sprintf("load script %v err : %v", path+v, err.Error())
		}
	}

	err := bot.bs.L.DoString(`meta.BotID = "` + bot.id + `"`)
	if err != nil {
		fmt.Println("set bot id", err.Error())
	}
	err = bot.bs.L.DoString(`meta.BotBatch = "` + batch + `"`)
	if err != nil {
		fmt.Println("set bot batch", err.Error())
	}
	err = bot.bs.L.DoString(`meta.BotName = "` + bot.name + `"`)
	if err != nil {
		fmt.Println("set bot name", err.Error())
	}

	return bot
}

func (b *Bot) loopThread(doneCh chan<- string, errch chan<- ErrInfo) {

	for {
		state, end := b.tick.Do()
		if end {
			doneCh <- b.id
			goto ext
		}

		if state == behavior.Break || state == behavior.Error {
			errch <- ErrInfo{
				ID:  b.id,
				Err: nil,
			}
			goto ext
		}

		time.Sleep(time.Millisecond * 10)
	}

ext:
	b.close()
}

func (b *Bot) RunByThread(doneCh chan<- string, errch chan<- ErrInfo) {

	go b.loopThread(doneCh, errch)

}

func (b *Bot) RunByBlock() error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run panic", err)
		}
	}()

	for {
		state, end := b.tick.Do()
		if end {
			return nil
		}

		if state == behavior.Break || state == behavior.Exit {
			return behavior.ErrorNodeHaveErr
		}

		time.Sleep(time.Millisecond * 10)
	}

}

func (b *Bot) GetReport() []script.Report {
	httpreport := b.bs.HttpMod.GetReport()
	tcpreport := b.bs.TCPMod.GetReport()

	return append(httpreport, tcpreport...)
}

func (b *Bot) close() {

	if b.bt.GetMode() == behavior.Thread {
		b.bs.L.DoString(`
		meta = {}
	`)
		pool.PutState(b.bs)
	} else {
		pool.FreeState(b.bs)
	}

}

type State int32

// 系统内部错误
const (
	SEnd State = 1 + iota
	SBreak
	SSucc
)

var stepmu sync.Mutex

func (b *Bot) RunByStep() State {
	stepmu.Lock()
	defer stepmu.Unlock()

	state, end := b.tick.Do()
	if end {
		return SEnd
	}

	if state == behavior.Break {
		return SBreak
	} else if state == behavior.Exit {
		return SEnd
	}

	return SSucc
}
