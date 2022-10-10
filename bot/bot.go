package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pojol/gobot/behavior"
	script "github.com/pojol/gobot/script/module"
	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
)

type ErrInfo struct {
	ID  string
	Err error
}

type RunMode int

const (
	Debug RunMode = 1 + iota
	Batch
)

type Thread struct {
	Number int    // 当前线程的编号
	ErrMsg string // 错误信息

	PreNods  []string
	NextNods []string
}

type Transaction struct {
	next   []*behavior.Tree
	thread *Thread
}

type Bot struct {
	id   string
	name string
	mode RunMode

	preloadErr string

	tree    *behavior.Tree
	running bool

	threadnum     int32
	threadDoneNum int32

	threadChan chan *Transaction
	waitChan   chan *Transaction
	waitLst    []*Transaction
	step       *utils.Switch

	threadLst  []*Thread
	threadDone chan interface{}

	sync.RWMutex
	bs *behavior.BotState // lua state pool

	donech chan<- string
	errch  chan<- ErrInfo
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
	threadinfolst := []Thread{}

	b.RLock()
	defer b.RUnlock()

	for _, v := range b.threadLst {
		threadinfolst = append(threadinfolst, Thread{
			Number:   v.Number,
			ErrMsg:   v.ErrMsg,
			PreNods:  v.PreNods,
			NextNods: v.NextNods,
		})
	}

	info, err := json.Marshal(&threadinfolst)
	if err != nil {
		fmt.Println(err.Error())
	}

	return string(info)
}

func NewWithBehaviorTree(path string, bt *behavior.Tree, name string, idx int32, globalScript []string) *Bot {

	bot := &Bot{
		id:         strconv.Itoa(int(idx)),
		tree:       bt,
		bs:         behavior.GetState(),
		name:       name,
		threadChan: make(chan *Transaction, 1),
		waitChan:   make(chan *Transaction, 1),
		threadDone: make(chan interface{}),
		step:       utils.NewSwitch(),
	}

	rand.Seed(time.Now().UnixNano())

	// 加载预定义全局脚本文件
	for _, gs := range globalScript {
		behavior.DoString(bot.bs.L, gs)
	}

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		err := behavior.DoFile(bot.bs.L, path+v)
		if err != nil {
			fmt.Println("err", err.Error())
			bot.preloadErr = fmt.Sprintf("load script %v err : %v", path+v, err.Error())
		}
	}

	err := bot.bs.L.DoString(`meta.BotID = "` + bot.id + `"`)
	if err != nil {
		fmt.Println("set bot id", err.Error())
	}
	err = bot.bs.L.DoString(`meta.BotName = "` + bot.name + `"`)
	if err != nil {
		fmt.Println("set bot name", err.Error())
	}

	return bot
}

func (b *Bot) fillThreadInfo(t *Thread) {

	b.Lock()

	var f bool

	for k, v := range b.threadLst {
		if v.Number == t.Number {
			b.threadLst[k] = t
			f = true
			break
		}
	}

	if !f {
		b.threadLst = append(b.threadLst, t)
	}

	/*
		fmt.Println("thread info ===>")
		for _, v := range b.threadLst {
			fmt.Println("\tthread"+strconv.Itoa(v.number), v.nodeid, v.err)
		}
	*/

	b.Unlock()

}

func getNodsName(nods []*behavior.Tree) []string {

	names := []string{}
	for _, v := range nods {
		names = append(names, v.ID)
	}

	return names
}

func (b *Bot) next(parent *Transaction) bool {
	trans := make(map[int]*Transaction)
	var batch []*behavior.ThreadTree

	for _, v := range parent.next {
		tt := v.Tick(parent.thread.Number)
		if len(tt.Children) > 0 {
			batch = append(batch, tt)
		}
	}

	if len(batch) == 0 {
		return false
	}

	for _, b := range batch {

		trans[b.ThreadNum] = &Transaction{
			thread: &Thread{
				Number:   b.ThreadNum,
				PreNods:  getNodsName(parent.next),
				NextNods: getNodsName(b.Children),
			},
			next: b.Children,
		}

	}

	for _, v := range trans {
		if b.mode == Batch {
			b.threadChan <- v
		} else if b.mode == Debug {
			b.waitChan <- v
		}
	}

	return true
}

func (b *Bot) loop() {

	for {

		select {
		case t := <-b.threadChan:
			b.fillThreadInfo(t.thread)
			b.next(t)

			if t.thread.ErrMsg != "" {
				b.errch <- ErrInfo{ID: b.id, Err: errors.New(t.thread.ErrMsg)}
				goto ext
			}

			if atomic.LoadInt32(&b.threadnum) == b.threadDoneNum {
				b.donech <- b.id
				goto ext
			}

		case w := <-b.waitChan:
			b.Lock()
			for _, v := range w.next {
				fmt.Printf("\t thread:%v id:%v type:%v\n", w.thread.Number, v.ID, v.Ty)
			}

			b.waitLst = append(b.waitLst, w)
			b.Unlock()
		case <-b.step.Done():
			if b.step.HasOpend() {
				fmt.Println("install:")
				for _, v := range b.waitLst {
					fmt.Printf("\t thread:%v\n", v.thread.Number)
					b.threadChan <- v
				}

				b.waitLst = b.waitLst[:0]
			}
			b.step.Close()
		}
	}

ext:
	fmt.Println("clean")
	// cleanup
	b.close()
}

func (b *Bot) Run(doneCh chan<- string, errch chan<- ErrInfo, mode RunMode) {

	b.donech = doneCh
	b.errch = errch
	b.mode = mode

	go b.loop()

	if len(b.tree.Children) != 0 {
		b.threadChan <- &Transaction{
			next: []*behavior.Tree{b.tree.Children[0]},
			thread: &Thread{
				Number:  int(atomic.AddInt32(&b.threadnum, 1)),
				PreNods: []string{b.tree.Children[0].ID},
			},
		}
	}

}

func (b *Bot) RunByBlock() error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run panic", err)
		}
	}()

	donech := make(chan string)
	errch := make(chan ErrInfo)

	b.Run(donech, errch, Batch)

	select {
	case <-donech:
		return nil
	case e := <-errch:
		return e.Err
	}
}

func (b *Bot) GetReport() []script.Report {
	return b.bs.HttpMod.GetReport()
}

func (b *Bot) close() {
	b.bs.L.DoString(`
		meta = {}
	`)
	behavior.PutState(b.bs)
}

type State int32

// 系统内部错误
const (
	SEnd State = 1 + iota
	SBreak
	SSucc
)

func (b *Bot) RunStep() State {

	if b.running {
		fmt.Println("")
		fmt.Println("cmd : step ->")
		b.Lock()
		b.threadLst = b.threadLst[:0]
		b.Unlock()
		b.step.Open()

	} else {
		fmt.Println("cmd : create ->")

		donech := make(chan string)
		errch := make(chan ErrInfo)

		b.Run(donech, errch, Debug)
		b.running = true
	}

	return SSucc
}
