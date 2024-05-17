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
	"github.com/pojol/gobot/driver/bot/behavior"
	"github.com/pojol/gobot/driver/bot/pool"
	script "github.com/pojol/gobot/driver/script/module"
	"github.com/pojol/gobot/driver/utils"
	lua "github.com/yuin/gopher-lua"
)

type ErrInfo struct {
	ID  string
	Err error
}

type Bot struct {
	id   string
	name string

	batch      string
	preloadErr string
	mod        behavior.Mode

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
	table, ok := b.bs.L.GetGlobal("bot").(*lua.LTable)

	var tableerr error
	var byt []byte

	if b.preloadErr != "" {
		byt, _ = json.Marshal(&tablemap)
		goto ext
	}

	if ok {
		tablemap, tableerr = utils.Table2Map(table)
		if tableerr != nil {
			tablemap["Err"] = tableerr.Error()
			goto ext
		}

		byt, tableerr = json.Marshal(&tablemap)
		if tableerr != nil {
			tablemap["Err"] = tableerr.Error()
			byt, _ = json.Marshal(&tablemap)
			goto ext
		}

	} else {

		tablemap["Err"] = errors.New("the meta field is not obtained")
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

func NewWithBehaviorTree(path string, bt *behavior.Tree, mode behavior.Mode, name, batch string, idx int32, globalScript string) *Bot {

	bb := &behavior.Blackboard{
		Nods:      []behavior.INod{bt.GetRoot()},
		Threadlst: []behavior.ThreadInfo{{Number: 1}},
	}

	var state *pool.BotState
	var id string

	if mode == behavior.Thread { // batch mode
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
		id:    id,
		bb:    bb,
		batch: batch,
		mod:   mode,
		tick:  behavior.NewTick(bb, state, strconv.Itoa(int(idx))),
		bt:    bt,
		bs:    state,
		name:  name,
	}

	rand.Seed(time.Now().UnixNano())

	// 加载预定义全局脚本文件
	if globalScript != "" {
		pool.DoString(bot.bs.L, globalScript)
	}
	script.RegisterMessageType(bot.bs.L)
	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		err := pool.DoFile(bot.bs.L, path+v)
		if err != nil {
			fmt.Println("err", err.Error())
			bot.preloadErr = fmt.Sprintf("load script %v err : %v", path+v, err.Error())
		}
	}

	err := bot.bs.L.DoString(`bot.Meta.ID = "` + bot.id + `"`)
	if err != nil {
		fmt.Println("set bot id", err.Error())
	}
	err = bot.bs.L.DoString(`bot.Meta.Batch = "` + batch + `"`)
	if err != nil {
		fmt.Println("set bot batch", err.Error())
	}
	err = bot.bs.L.DoString(`bot.Meta.Name = "` + bot.name + `"`)
	if err != nil {
		fmt.Println("set bot name", err.Error())
	}

	bot.addLog(fmt.Sprintf("create bot id %v name %v success", bot.id, bot.name))

	return bot
}

func (b *Bot) loopThread(doneCh chan<- string, errch chan<- ErrInfo) {

	for {
		state, end, logs := b.tick.Do(b.mod)
		if len(logs) != 0 {
			for _, log := range logs {
				b.addLog(log)
			}
		}

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
		state, end, logs := b.tick.Do(b.mod)
		if len(logs) != 0 {
			for _, log := range logs {
				b.addLog(log)
			}
		}

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
	wsreport := b.bs.WebsocketMod.GetReport()

	report := []script.Report{}
	report = append(report, httpreport...)
	report = append(report, tcpreport...)
	report = append(report, wsreport...)

	return report
}

func (b *Bot) close() {

	if b.bb.HaveErr() {
		info := b.bb.ThreadInfo()
		for _, v := range info {
			fmt.Println("bot", b.Name(), v.CurNod, "err", v.ErrMsg)
		}
	}

	if b.mod == behavior.Thread {
		b.bs.L.DoString(`
		bot = {
			Meta = {
				ID = "",
				Name = "",
				Batch = "",
				Err = "",
			}
		}
	`)
		pool.PutState(b.bs)
		b.bt.Reset()
		behavior.Put(b.name, b.bt)
	} else {
		pool.FreeState(b.bs)
	}

	b.addLog(fmt.Sprintf("close bot id %v name %v success", b.id, b.name))
}

// PopLog - 弹出一条日志
func (b *Bot) PopLog() string {
	line := b.bs.LogMod.Pop()

	return line
}

func (b *Bot) addLog(log string) {
	fmt.Println("=>", log)

	log = time.Now().Format("2006-01-02 15:04:05") + " =================>\n" + log

	if b.mod != behavior.Thread {
		b.bs.LogMod.Push(log)
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

	// 这边的错误日志需要记录下
	state, end, logs := b.tick.Do(b.mod)
	if len(logs) != 0 {
		for _, log := range logs {
			b.addLog(log)
		}
	}

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
