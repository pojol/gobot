package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
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

	httpMod   *script.HttpModule
	protoMod  *script.ProtoModule
	utilsMod  *script.UtilsModule
	base64Mod *script.Base64Module

	L *lua.LState
	sync.Mutex
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

func (b *Bot) GetMetadata() (string, string, error) {

	if b.preloadErr != "" {
		return b.preloadErr, "", nil
	}

	meta, err := utils.Table2Map(b.L.GetGlobal("meta").(*lua.LTable))
	if err != nil {
		return "", "", err
	}

	metabyt, err := json.Marshal(&meta)
	if err != nil {
		return "", "", err
	}

	change, err := utils.Table2Map(b.L.GetGlobal("change").(*lua.LTable))
	if err != nil {
		return "", "", err
	}

	changebyt, err := json.Marshal(&change)
	if err != nil {
		return "", "", err
	}

	return string(metabyt), string(changebyt), nil

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

func NewWithBehaviorTree(path string, bt *behavior.Tree, tmpl string) *Bot {

	bot := &Bot{
		id:        uuid.New().String(),
		tree:      bt,
		cur:       bt,
		L:         lua.NewState(),
		name:      tmpl,
		httpMod:   script.NewHttpModule(&http.Client{}),
		protoMod:  &script.ProtoModule{},
		utilsMod:  &script.UtilsModule{},
		base64Mod: &script.Base64Module{},
	}

	rand.Seed(time.Now().UnixNano())

	bot.L.PreloadModule("proto", bot.protoMod.Loader)
	bot.L.PreloadModule("http", bot.httpMod.Loader)
	bot.L.PreloadModule("utils", bot.utilsMod.Loader)
	bot.L.PreloadModule("base64", bot.base64Mod.Loader)

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		err := bot.L.DoFile(path + v)
		if err != nil {
			bot.preloadErr = fmt.Sprintf("load script %v err : %v", path+v, err.Error())
		}
	}

	return bot
}

func (b *Bot) run_selector(nod *behavior.Tree, next bool) (bool, error) {

	if next {
		for k := range nod.Children {
			ok, _ := b.run_nod(nod.Children[k], true)
			if ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_assert(nod *behavior.Tree, next bool) (bool, error) {

	err := b.L.DoString(nod.Code)
	if err != nil {
		return false, err
	}

	err = b.L.CallByParam(lua.P{
		Fn:      b.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return false, err
	}
	v := b.L.Get(-1)
	b.L.Pop(1)

	if lua.LVAsBool(v) {
		if next {
			b.run_children(nod, nod.Children)
		}

		return true, nil
	}

	return false, errors.New("assert failed")
}

func (b *Bot) run_sequence(nod *behavior.Tree, next bool) (bool, error) {
	if next {
		for k := range nod.Children {
			ok, _ := b.run_nod(nod.Children[k], true)
			if !ok {
				break
			}
		}
	}

	return true, nil
}

func (b *Bot) run_condition(nod *behavior.Tree, next bool) (bool, error) {

	err := b.L.DoString(nod.Code)
	if err != nil {
		return false, err
	}

	err = b.L.CallByParam(lua.P{
		Fn:      b.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return false, err
	}

	v := b.L.Get(-1)
	b.L.Pop(1)

	if lua.LVAsBool(v) {
		if next {
			b.run_children(nod, nod.Children)
		}

		return true, nil
	}

	return false, nil
}

func (b *Bot) run_wait(nod *behavior.Tree, next bool) (bool, error) {
	time.Sleep(time.Millisecond * time.Duration(nod.Wait))

	if next {
		b.run_children(nod, nod.Children)
	}

	return true, nil
}

func (b *Bot) run_loop(nod *behavior.Tree, next bool) (bool, error) {

	if nod.Loop == 0 {
		for {
			if next {
				b.run_children(nod, nod.Children)
			}
			time.Sleep(time.Millisecond)
		}
	} else {

		if next {
			for i := 0; i < int(nod.Loop); i++ {
				b.run_children(nod, nod.Children)
				time.Sleep(time.Millisecond)
			}
		}
	}

	return true, nil
}

func (b *Bot) run_script(nod *behavior.Tree, next bool) (bool, error) {

	err := b.L.DoString(nod.Code)
	if err != nil {
		return false, err
	}

	err = b.L.CallByParam(lua.P{
		Fn:      b.L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return false, err
	}
	//v := b.L.Get(-1)
	b.L.Pop(1)

	// mergo.MergeWithOverwrite(&b.metadata, t)

	if next {
		b.run_children(nod, nod.Children)
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
		ok, _ = b.run_sequence(nod, next)
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

func (b *Bot) run_children(parent *behavior.Tree, children []*behavior.Tree) {
	for k := range children {
		b.run_nod(children[k], true)
	}
}

func (b *Bot) Run(doneCh chan string, errch chan ErrInfo) {

	go func() {

		defer func() {
			if err := recover(); err != nil {
				fmt.Println("run err", err)
				errch <- ErrInfo{
					ID:  b.id,
					Err: err.(error),
				}
			}
		}()

		b.run_children(b.tree, b.tree.Children)
		doneCh <- b.id
	}()

}

func (b *Bot) GetReport() []script.Report {
	return b.httpMod.GetReport()
}

func (b *Bot) Close() {
	b.L.Close()
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
