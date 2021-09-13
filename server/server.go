package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pojol/apibot/behavior"
	"github.com/pojol/apibot/bot"
	"github.com/pojol/apibot/factory"
)

type Response struct {
	Code int
	Msg  string
	Body interface{}
}

type Err int32

const (
	Succ Err = 200 + iota
	ErrContentRead
	ErrJsonInvalid
	ErrPluginLoad
	ErrMetaData
	ErrEnd
	ErrBreak
	ErrCantFindBot
	ErrCreateBot
)

var errmap map[Err]string = map[Err]string{
	ErrContentRead: "failed to read request content",
	ErrJsonInvalid: "wrong file format",
	ErrPluginLoad:  "failed to plugin load",
	ErrMetaData:    "failed to get meta data",
	ErrEnd:         "run to the end",
	ErrBreak:       "run to the break",
	ErrCantFindBot: "can't find bot",
	ErrCreateBot:   "failed to create bot, the behavior tree file needs to be uploaded to the server before creation",
}

func Upload(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	code := Succ
	var tree *behavior.Tree

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	tree, err = behavior.New(bts)
	if err != nil {
		fmt.Println(err.Error())
		code = ErrJsonInvalid
		goto EXT
	}
	factory.Global.AppendBehavior(tree.ID, bts)

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]

	ctx.JSON(http.StatusOK, res)
	return nil
}

type RunRequest struct {
	BotID string
}

type RunResponse struct {
	Blackboard string
}

func Run(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	req := &RunRequest{}
	body := &RunResponse{}
	var b *bot.Bot
	code := Succ

	var err error

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, &req)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	b = factory.Global.FindBot(req.BotID)
	if b == nil {
		code = ErrCantFindBot
		goto EXT
	}
	b.Run(nil, nil, nil)
	defer factory.Global.RmvBot(req.BotID)

	body.Blackboard, err = b.GetMetadata()
	if err != nil {
		fmt.Println(err.Error())
		code = ErrMetaData
	}

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

type StepRequest struct {
	BotID string
}

type StepResponse struct {
	Prev       string
	Cur        string
	Blackboard string
}

func Step(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	req := &StepRequest{}
	body := &StepResponse{}
	code := Succ
	var b *bot.Bot

	var err error
	var s bot.State

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, &req)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	b = factory.Global.FindBot(req.BotID)
	if b == nil {
		code = ErrCantFindBot
		goto EXT
	}

	s = b.RunStep()
	body.Blackboard, err = b.GetMetadata()
	fmt.Println("blackboard", body.Blackboard, err)
	body.Cur = b.GetCurNodeID()
	body.Prev = b.GetPrevNodeID()

	if s == bot.SEnd {
		code = ErrEnd
		defer factory.Global.RmvBot(req.BotID)
		goto EXT
	} else if s == bot.SBreak {
		code = ErrBreak
		defer factory.Global.RmvBot(req.BotID)
		goto EXT
	}

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

type createRequest struct {
	TreeID string
}

type createResponse struct {
	BotID string
}

func Create(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	req := &createRequest{}
	body := &createResponse{}
	code := Succ
	var b *bot.Bot

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, &req)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	b = factory.Global.CreateBot(req.TreeID)
	if b == nil {
		fmt.Println("create bot err", req.TreeID)
		code = ErrCreateBot
		goto EXT
	}

	body.BotID = b.ID()

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

func Route(e *echo.Echo) {

	e.POST("/upload", Upload) // 上传行为树模版文件
	e.POST("/create", Create) // 创建一个bot
	e.POST("/run", Run)       // 运行一个bot
	e.POST("/step", Step)     // 单步运行一个bot

}
