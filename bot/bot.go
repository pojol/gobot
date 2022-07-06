package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
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

type Bot struct {
	id   string
	name string

	preloadErr string

	tree *behavior.Tree

	prev *behavior.Tree
	cur  *behavior.Tree

	sync.Mutex
	bs *botState

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

func NewWithBehaviorTree(path string, bt *behavior.Tree, name string, idx int32, globalScript string) *Bot {

	bot := &Bot{
		id:   strconv.Itoa(int(idx)),
		tree: bt,
		cur:  bt,
		bs:   luaPool.Get(),
		name: name,
	}

	rand.Seed(time.Now().UnixNano())

	// 加载预定义全局脚本文件
	DoString(bot.bs.L, globalScript)

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

func (b *Bot) run_selector(nod *behavior.Tree, next bool) (bool, error) {

	if next {
		for k := range nod.Children {
			ok, err := b.run_nod(nod.Children[k], true)

			if err != nil {
				return false, err
			}

			if ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_assert(nod *behavior.Tree, next bool) (bool, error) {

	err := DoString(b.bs.L, nod.Code)
	if err != nil {
		return false, err
	}

	err = b.bs.L.CallByParam(lua.P{
		Fn:      b.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return false, err
	}
	v := b.bs.L.Get(-1)
	b.bs.L.Pop(1)

	if lua.LVAsBool(v) {
		if next {
			err = b.run_children(nod, nod.Children)
			if err != nil {
				return false, err
			}
		}

		return true, nil
	}

	return false, fmt.Errorf("node %v assert failed", nod.ID)
}

func (b *Bot) run_sequence(nod *behavior.Tree, next bool) (bool, error) {
	if next {
		for k := range nod.Children {
			ok, err := b.run_nod(nod.Children[k], true)

			if err != nil {
				return false, err
			}

			if !ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_condition(nod *behavior.Tree, next bool) (bool, error) {

	err := DoString(b.bs.L, nod.Code)
	if err != nil {
		return false, err
	}

	err = b.bs.L.CallByParam(lua.P{
		Fn:      b.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return false, err
	}

	v := b.bs.L.Get(-1)
	b.bs.L.Pop(1)

	if lua.LVAsBool(v) {
		if next {
			err = b.run_children(nod, nod.Children)
			if err != nil {
				return false, err
			}
		}

		return true, nil
	}

	return false, nil
}

func (b *Bot) run_wait(nod *behavior.Tree, next bool) (bool, error) {
	time.Sleep(time.Millisecond * time.Duration(nod.Wait))

	if next {
		err := b.run_children(nod, nod.Children)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (b *Bot) run_loop(nod *behavior.Tree, next bool) (bool, error) {

	var err error

	if nod.Loop == 0 {
		for {
			if next {
				err = b.run_children(nod, nod.Children)
				if err != nil {
					goto ext
				}
			}
			time.Sleep(time.Millisecond)
		}
	} else {

		if next {
			for i := 0; i < int(nod.Loop); i++ {
				err = b.run_children(nod, nod.Children)
				if err != nil {
					goto ext
				}
			}
		}
	}

ext:
	return true, err
}

func (b *Bot) run_script(nod *behavior.Tree, next bool) (bool, error) {

	err := DoString(b.bs.L, nod.Code)
	if err != nil {
		return false, err
	}

	err = b.bs.L.CallByParam(lua.P{
		Fn:      b.bs.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return false, err
	}
	//v := b.bs.L.Get(-1)
	b.bs.L.Pop(1)

	// mergo.MergeWithOverwrite(&b.metadata, t)

	if next {
		err = b.run_children(nod, nod.Children)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (b *Bot) run_nod(nod *behavior.Tree, next bool) (bool, error) {

	var ok bool
	var err error

	switch nod.Ty {
	case behavior.SELETE:
		ok, err = b.run_selector(nod, next)
	case behavior.SEQUENCE:
		ok, err = b.run_sequence(nod, next)
	case behavior.CONDITION:
		ok, err = b.run_condition(nod, next)
	case behavior.WAIT:
		ok, _ = b.run_wait(nod, next)
	case behavior.LOOP:
		ok, err = b.run_loop(nod, next)
	case behavior.ACTION:
		ok, err = b.run_script(nod, next)
	case behavior.ASSERT:
		ok, err = b.run_assert(nod, next)
	case behavior.ROOT:
		ok = true
	default:
		ok = false
		err = fmt.Errorf("run node type err %s", nod.Ty)
	}

	if err != nil {
		fmt.Println("run node err", nod.ID, nod.Ty, err.Error())
	}

	return ok, err
}

func (b *Bot) run_children(parent *behavior.Tree, children []*behavior.Tree) error {
	var err error

	for k := range children {
		_, err = b.run_nod(children[k], true)
		if err != nil {
			break
		}
	}

	return err
}

func (b *Bot) Run(doneCh chan string, errch chan ErrInfo) {

	go func() {

		defer func() {
			if err := recover(); err != nil {
				errch <- ErrInfo{
					ID:  b.id,
					Err: err.(error),
				}
			}
		}()

		err := b.run_children(b.tree, b.tree.Children)
		if err != nil {
			errch <- ErrInfo{
				ID:  b.id,
				Err: err,
			}
			return
		}

		doneCh <- b.id
	}()

}

func (b *Bot) RunByBlock() error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run panic", err)
		}
	}()

	err := b.run_children(b.tree, b.tree.Children)
	if err != nil {
		fmt.Println("run block err", err)
	}

	return err
}

func (b *Bot) GetReport() []script.Report {
	return b.bs.httpMod.GetReport()
}

func (b *Bot) Close() {
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

	return SSucc
}
