package bot

import (
	"encoding/json"
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
	num   int
	preid string
	err   error
}

type Transaction struct {
	cur    *behavior.Tree
	parent *behavior.Tree
	thread *Thread
}

type Bot struct {
	id   string
	name string
	mode RunMode

	preloadErr string

	tree *behavior.Tree
	prev *behavior.Tree
	cur  *behavior.Tree

	threadnum     int32
	threadDoneNum int32

	threadChan chan *Transaction
	waitChan   chan *Transaction
	waitLst    []*Transaction
	next       *utils.Switch

	threadLst  []*Thread
	threadDone chan interface{}

	sync.Mutex
	bs *botState

	donech chan<- string
	errch  chan<- ErrInfo

	runtimeErr string
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

func (b *Bot) GetMetadata() (string, string, string, bool) {

	var metaStr, changeStr string

	if b.preloadErr != "" {
		return b.preloadErr, "", "", true
	}

	metaTable, ok := b.bs.L.GetGlobal("meta").(*lua.LTable)
	if ok {
		meta, err := utils.Table2Map(metaTable)
		if err != nil {
			return "", "", err.Error(), false
		}

		metabyt, err := json.Marshal(&meta)
		if err != nil {
			return "", "", err.Error(), false
		}

		metaStr = string(metabyt)
	} else {
		b.runtimeErr += "\nThe meta field is not obtained"
	}

	changeTable, ok := b.bs.L.GetGlobal("change").(*lua.LTable)
	if ok {
		change, err := utils.Table2Map(changeTable)
		if err != nil {
			return "", "", err.Error(), false
		}

		changebyt, err := json.Marshal(&change)
		if err != nil {
			return "", "", err.Error(), false
		}

		changeStr = string(changebyt)
	} else {
		b.runtimeErr += "\nThe change field is not obtained"
	}

	return metaStr, changeStr, b.runtimeErr, true

}

func (b *Bot) GetCurNodeID() string {
	if b.cur != nil {
		return b.cur.ID
	}
	return ""
}

func (b *Bot) GetPrevNodeID() string {
	if b.prev != nil {
		return b.prev.ID
	}
	return ""
}

func NewWithBehaviorTree(path string, bt *behavior.Tree, name string, idx int32, globalScript []string) *Bot {

	bot := &Bot{
		id:         strconv.Itoa(int(idx)),
		tree:       bt,
		cur:        bt,
		bs:         luaPool.Get(),
		name:       name,
		threadChan: make(chan *Transaction, 1),
		waitChan:   make(chan *Transaction, 1),
		threadDone: make(chan interface{}),
		next:       utils.NewSwitch(),
	}

	rand.Seed(time.Now().UnixNano())

	// 加载预定义全局脚本文件
	for _, gs := range globalScript {
		DoString(bot.bs.L, gs)
	}

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		err := DoFile(bot.bs.L, path+v)
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
		if v.num == t.num {
			b.threadLst[k] = t
			f = true
			break
		}
	}

	if !f {
		b.threadLst = append(b.threadLst, t)
	}

	fmt.Println("thread info ===>")
	for _, v := range b.threadLst {
		fmt.Println("\tthread"+strconv.Itoa(v.num), v.preid, v.err)
	}

	b.Unlock()

}

func (b *Bot) do(t *Transaction) bool {
	var ok bool

	switch t.cur.Ty {
	case behavior.SELETE:
		ok = b.strategySelete(t)
	case behavior.SEQUENCE:
		ok = b.strategySequence(t)
	case behavior.CONDITION:
		ok = b.strategyCondition(t)
	case behavior.WAIT:
		ok = b.strategyWait(t)
	case behavior.LOOP:
		ok = b.strategyLoop(t)
	case behavior.ASSERT:
		ok = b.strategyAssert(t)
	case behavior.PARALLEL:
		ok = b.strategyParallel(t)
	case behavior.ROOT:
		ok = true
	default:
		ok = b.strategyScript(t)
	}

	return ok
}

func (b *Bot) strategySelete(t *Transaction) bool {

	batch := func() bool {
		for k := range t.cur.Children {

			tr := &Transaction{
				cur:    t.cur.Children[k],
				parent: t.cur,
				thread: &Thread{
					num:   t.thread.num,
					preid: t.cur.ID,
				},
			}
			ok := b.do(tr)
			b.fillThreadInfo(tr.thread)

			if tr.thread.err != nil {
				return false
			}

			if ok {
				break
			}

		}
		return true
	}

	step := func() bool {

		tr := &Transaction{
			cur:    t.parent.Next(),
			parent: t.parent,
			thread: &Thread{
				num:   t.thread.num,
				preid: t.cur.ID,
			},
		}
		b.waitChan <- tr

		return true
	}

	if b.mode == Batch {
		return batch()
	}
	return step()
}

func (b *Bot) strategySequence(t *Transaction) bool {

	for k := range t.cur.Children {

		tr := &Transaction{
			cur: t.cur.Children[k],
			thread: &Thread{
				num:   t.thread.num,
				preid: t.cur.ID,
			},
		}

		ok := b.do(tr)
		b.fillThreadInfo(tr.thread)

		if tr.thread.err != nil {
			return false
		}

		if !ok {
			break
		}
	}

	return true
}

func (b *Bot) strategyParallel(t *Transaction) bool {

	for k := range t.cur.Children {

		tr := &Transaction{
			cur: t.cur.Children[k],
			thread: &Thread{
				num:   int(atomic.AddInt32(&b.threadnum, 1)),
				preid: t.cur.ID,
			},
		}
		b.fillThreadInfo(tr.thread)

		b.threadChan <- tr

	}

	return true
}

func (b *Bot) strategyWait(t *Transaction) bool {

	time.Sleep(time.Millisecond * time.Duration(t.cur.Wait))

	if b.mode == Batch && len(t.cur.Children) != 0 {

		tr := &Transaction{
			cur: t.cur.Children[0],
			thread: &Thread{
				num:   t.thread.num,
				preid: t.cur.ID,
			},
		}
		b.fillThreadInfo(tr.thread)
		b.do(t)

	}

	return true
}

func (b *Bot) strategyCondition(t *Transaction) bool {

	tr := &Transaction{
		thread: &Thread{
			num:   t.thread.num,
			preid: t.cur.ID,
		},
	}
	defer b.fillThreadInfo(tr.thread)

	err := DoString(b.bs.L, t.cur.Code)
	if err != nil {
		tr.thread.err = err
		return false
	}

	err = b.bs.L.CallByParam(lua.P{
		Fn:      b.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		t.thread.err = err
		return false
	}

	v := b.bs.L.Get(-1)
	b.bs.L.Pop(1)

	if lua.LVAsBool(v) {
		if b.mode == Batch && len(t.cur.Children) != 0 {

			tr.cur = t.cur.Children[0]
			b.do(tr)
		}

	}

	return true
}

func (b *Bot) strategyLoop(t *Transaction) bool {

	for i := 0; i < int(t.cur.Loop); i++ {
		if b.mode == Batch && len(t.cur.Children) != 0 {
			tr := &Transaction{
				cur: t.cur.Children[0],
				thread: &Thread{
					num:   t.thread.num,
					preid: t.cur.ID,
				},
			}
			b.fillThreadInfo(tr.thread)
			b.do(tr)
		}
	}

	return true
}

func (b *Bot) strategyScript(t *Transaction) bool {

	tr := &Transaction{
		thread: &Thread{
			num:   t.thread.num,
			preid: t.cur.ID,
		},
	}
	defer b.fillThreadInfo(tr.thread)

	tr.thread.err = DoString(b.bs.L, t.cur.Code)
	if tr.thread.err != nil {
		return false
	}

	tr.thread.err = b.bs.L.CallByParam(lua.P{
		Fn:      b.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if tr.thread.err != nil {
		return false
	}

	b.bs.L.Pop(1)

	if b.mode == Batch && len(t.cur.Children) != 0 {
		tr.cur = t.cur.Children[0]
		b.do(tr)
	}

	return true
}

func (b *Bot) strategyAssert(t *Transaction) bool {

	tr := &Transaction{
		thread: &Thread{
			num:   t.thread.num,
			preid: t.cur.ID,
		},
	}

	tr.thread.err = DoString(b.bs.L, t.cur.Code)
	if tr.thread.err != nil {
		return false
	}

	tr.thread.err = b.bs.L.CallByParam(lua.P{
		Fn:      b.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if tr.thread.err != nil {
		return false
	}
	v := b.bs.L.Get(-1)
	b.bs.L.Pop(1)

	if lua.LVAsBool(v) {

		if len(t.cur.Children) != 0 {
			tr.cur = t.cur.Children[0]
			b.fillThreadInfo(tr.thread)
			b.do(tr)
		}

		return true
	}

	return false
}

func (b *Bot) loop() {

	for {

		select {
		case t := <-b.threadChan:
			b.do(t)
			if t.thread.err != nil {
				b.errch <- ErrInfo{ID: b.id, Err: t.thread.err}
				goto ext
			}

			b.threadDoneNum++

			if atomic.LoadInt32(&b.threadnum) == b.threadDoneNum {
				b.donech <- b.id
				goto ext
			}
		case w := <-b.waitChan:
			b.waitLst = append(b.waitLst, w)
		case <-b.next.Done():
			if b.next.HasOpend() {

				for _, v := range b.waitLst {
					b.threadChan <- v
				}

				b.waitLst = b.waitLst[:0]
			}
			b.next.Close()
		}

	}

ext:
	// cleanup
	b.close()
}

func (b *Bot) Run(doneCh chan<- string, errch chan<- ErrInfo, mode RunMode) {

	b.donech = doneCh
	b.errch = errch
	b.mode = mode

	tn := int(atomic.AddInt32(&b.threadnum, 1))

	go b.loop()

	if len(b.tree.Children) != 0 {
		b.threadChan <- &Transaction{
			cur:    b.cur.Children[0],
			parent: b.cur,
			thread: &Thread{num: tn, preid: b.tree.ID},
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
	return b.bs.httpMod.GetReport()
}

func (b *Bot) close() {
	b.bs.L.DoString(`
		meta = {}
	`)
	luaPool.Put(b.bs)
}

type State int32

// 系统内部错误
const (
	SEnd State = 1 + iota
	SBreak
	SSucc
)

func (b *Bot) RunStep() State {

	/*
		if b.cur == nil {
			return SEnd
		}

		b.Lock()
		defer b.Unlock()

		f, err := b.run_nod(b.cur, false)
		if err != nil {
			b.runtimeErr = err.Error()
			return SBreak
		}
		// step 中使用了sleep之后，会有多个goroutine执行接下来的程序
		// fmt.Println(goid.Get())

		if b.cur.Parent != nil {
			if b.cur.Parent.Ty == behavior.SELETE && f {
				b.cur.Parent.Step = len(b.cur.Parent.Children)
			}
		}

		if f && b.cur.Step < len(b.cur.Children) {
			// down
			nextidx := b.cur.Step
			b.cur.Step++
			next := b.cur.Children[nextidx]
			b.prev = b.cur
			b.cur = next
		} else {
			// right
			if b.cur.Parent != nil {
				b.prev = b.cur
				b.cur = b.cur.Parent.Next()
				if b.cur == nil {
					return SEnd
				}
			}
		}
	*/
	return SSucc
}
