package behavior

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pojol/gobot/bot/pool"
	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
)

type Tick struct {
	blackboard *Blackboard
	bs         *pool.BotState
	botid      string
}

var (
	ErrorTickHave = errors.New("tick thread have error")
	ErrorNodeHave = errors.New("node script have error")
)

func NewTick(bb *Blackboard, state *pool.BotState, botid string) *Tick {
	t := &Tick{
		blackboard: bb,
		bs:         state,
		botid:      botid,
	}
	return t
}

func (t *Tick) stateCheck() (string, string, error) {

	var r1, r2 lua.LValue
	var err error
	var changestr string
	var changeByt []byte
	var changetab map[string]interface{}
	state := Succ

	r1 = t.bs.L.Get(-1)
	if r1.Type() != lua.LTNil {
		t.bs.L.Pop(1)
	}
	r2 = t.bs.L.Get(-1)
	if r2.Type() != lua.LTNil {
		t.bs.L.Pop(1)
		state = r2.String()
		if state != Error {
			tab, ok := r1.(*lua.LTable)
			if ok {
				changetab, err = utils.Table2Map(tab)
				if err != nil {
					goto ext
				}

				changeByt, err = json.Marshal(&changetab)
				if err != nil {
					goto ext
				}
				changestr = string(changeByt)
			}

		} else {
			changestr = r1.String()
			fmt.Println(changestr)
		}

	} else {
		goto ext //没有返回值，不需要处理
	}

ext:
	//
	return state, changestr, err
}

func (t *Tick) Do() (error, bool) {

	nods := t.blackboard.GetOpenNods()
	t.blackboard.ThreadInfoReset()

	var err, parseerr error
	var state, msg string

	for _, n := range nods {
		err = n.onTick(t)

		if n.getMode() == Step {
			fmt.Println(n.getType(), n.getID())

			if n.getType() == SCRIPT {
				state, msg, parseerr = t.stateCheck()
			}

			threadInfo := ThreadInfo{
				Number: n.getThread(),
				CurNod: n.getID(),
				Change: msg,
			}

			if err != nil {
				threadInfo.ErrMsg = fmt.Sprintf("tick err %v", err.Error())
			}
			if parseerr != nil {
				threadInfo.ErrMsg = fmt.Sprintf("%v parse err %v", threadInfo.ErrMsg, parseerr.Error())
			}

			t.blackboard.ThreadFillInfo(threadInfo)

			if state != Succ {

				if state == Exit || state == Break {
					return nil, true
				} else if state == Error {
					// 节点脚本出错，脚本逻辑自行抛出的错误
					return ErrorNodeHave, false
				}

			}

			// tick中有错误，可能是脚本解析出错，也可能是返回值解析出错
			if t.blackboard.HaveErr() {
				return ErrorTickHave, false
			}
		}
	}

	t.blackboard.Reset()

	for _, n := range nods {
		n.onNext(t)
	}

	if t.blackboard.end {
		return nil, true
	}

	return nil, false
}
