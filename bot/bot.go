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
	Preid  string // 上一个节点的id
	Curid  string // 当前线程处理的节点id
	Errmsg string // 当前线程遇到的错误
	Ret    bool   // 判定类节点返回值
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

	sync.Mutex
	bs *botState // lua state pool

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

	b.Lock()
	defer b.Unlock()

	for _, v := range b.threadLst {
		fmt.Println("\t", v.Number, v.Curid, v.Errmsg)
		threadinfolst = append(threadinfolst, Thread{
			Number: v.Number,
			Preid:  v.Preid,
			Curid:  v.Curid,
			Errmsg: v.Errmsg,
			Ret:    v.Ret,
		})
	}

	info, err := json.Marshal(&threadinfolst)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("clean thread info")
	b.threadLst = b.threadLst[:0]

	return string(info)
}

func NewWithBehaviorTree(path string, bt *behavior.Tree, name string, idx int32, globalScript []string) *Bot {

	bot := &Bot{
		id:         strconv.Itoa(int(idx)),
		tree:       bt,
		bs:         luaPool.Get(),
		name:       name,
		threadChan: make(chan *Transaction, 1),
		waitChan:   make(chan *Transaction, 1),
		threadDone: make(chan interface{}),
		step:       utils.NewSwitch(),
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

func (b *Bot) next(parent *Transaction) {
	var trans []*Transaction

	fmt.Println("next", parent.cur.ID, parent.cur.Ty)

	children := parent.cur.Next(parent.thread.Ret)
	for k := range children {

		t := &Thread{
			Number: parent.thread.Number,
			Preid:  parent.cur.ID,
			Curid:  children[k].ID,
		}

		trans = append(trans, &Transaction{
			cur:    children[k],
			parent: parent.parent,
			thread: t,
		})
	}

	for _, v := range trans {
		if b.mode == Batch {
			b.threadChan <- v
		} else if b.mode == Debug {
			b.waitChan <- v
		}
	}
}

func (b *Bot) execute(t *Transaction) {

	fmt.Println("execute", t.cur.ID, t.cur.Ty)

	doscript := func(id string, ty string, code string) (error, bool) {
		err := DoString(b.bs.L, code)
		if err != nil {
			return err, false
		}

		err = b.bs.L.CallByParam(lua.P{
			Fn:      b.bs.L.GetGlobal("execute"),
			NRet:    1,
			Protect: true,
		}, lua.LString(id))
		if err != nil {
			return err, false
		}
		ret := b.bs.L.Get(-1)
		b.bs.L.Pop(1)

		if ty == behavior.CONDITION || ty == behavior.ASSERT {
			return nil, lua.LVAsBool(ret)
		}

		return nil, true
	}

	// 和控制节点不同，实际的脚本节点应该携带两个语意
	// 1. 执行脚本
	// 2. 进入到下一个节点
	getchildren := func(nod *behavior.Tree) {
		children := nod.Next(t.thread.Ret)
		if len(children) == 1 {
			t := &Transaction{
				cur:    children[0],
				parent: nod,
				thread: &Thread{
					Number: t.thread.Number,
					Preid:  nod.ID,
					Curid:  children[0].ID,
				},
			}

			if b.mode == Batch {
				b.threadChan <- t
			} else if b.mode == Debug {
				b.waitChan <- t
			}
		}
	}

	nod := t.cur

	switch nod.Ty {
	case behavior.WAIT:
		if nod.Wait >= 0 {
			time.Sleep(time.Millisecond * time.Duration(nod.Wait))
		}
		getchildren(nod)
	case behavior.LOOP, behavior.PARALLEL, behavior.SELETE, behavior.SEQUENCE, behavior.ROOT:
		break
	default: // script
		err, ok := doscript(nod.ID, nod.Ty, nod.Code)
		t.thread.Errmsg = err.Error()
		t.thread.Ret = ok
		getchildren(nod)
	}

}

func (b *Bot) loop() {

	for {

		select {
		case t := <-b.threadChan:
			fmt.Println("fill thread info")
			b.fillThreadInfo(t.thread)
			if !behavior.IsScriptNode(t.cur.Ty) {
				b.next(t)
			} else {
				b.execute(t)
			}

			if t.thread.Errmsg != "" {
				b.errch <- ErrInfo{ID: b.id, Err: errors.New(t.thread.Errmsg)}
				goto ext
			}

			if atomic.LoadInt32(&b.threadnum) == b.threadDoneNum {
				b.donech <- b.id
				goto ext
			}

		case w := <-b.waitChan:
			b.Lock()
			fmt.Println("\t", "append", w.cur.ID, w.cur.Ty)
			b.waitLst = append(b.waitLst, w)
			b.Unlock()
		case <-b.step.Done():
			if b.step.HasOpend() {
				fmt.Println("process wait list", len(b.waitLst))
				for _, v := range b.waitLst {
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
			cur:    b.tree.Children[0],
			parent: b.tree,
			thread: &Thread{Number: int(atomic.AddInt32(&b.threadnum, 1)), Curid: b.tree.ID},
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

	if b.running {

		b.step.Open()

	} else {
		donech := make(chan string)
		errch := make(chan ErrInfo)

		b.Run(donech, errch, Debug)
		b.running = true
	}

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
