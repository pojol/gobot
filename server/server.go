package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pojol/apibot/bot"
	"github.com/pojol/apibot/mock"
	"github.com/pojol/apibot/plugins"
)

type Response struct {
	Code int
	Msg  string
	Body interface{}
}

func Upload(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	code := 200
	msg := ""

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = -1 // tmp
		msg = err.Error()
		goto EXT
	}

	if !json.Valid(bts) {
		code = -2
		msg = "json invalid"
		goto EXT
	}

	err = plugins.Load("./json.so")
	if err != nil {
		msg = err.Error()
		code = -3
		goto EXT
	}

	mockServer.Reset()
	behaviorBuffer.Reset()
	behaviorBuffer.Write(bts)

	mbot, _ = bot.NewWithBehaviorFile(behaviorBuffer.Bytes(), mockServer.Url())

EXT:
	res.Code = code
	res.Msg = msg

	ctx.JSON(http.StatusOK, res)
	return nil
}

type RunResponse struct {
	Blackboard string
}

func Run(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	body := &RunResponse{}
	code := 200
	msg := ""

	var err error

	if behaviorBuffer.Len() == 0 {
		code = -1
		msg = "not behavior data, need upload!"
		goto EXT
	}

	mbot.Run()

	body.Blackboard, err = mbot.GetMetadata()
	if err != nil {
		msg = err.Error()
		code = -3
	}

EXT:
	res.Code = code
	res.Msg = msg
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

type StepResponse struct {
	Prev       string
	Cur        string
	Blackboard string
}

func Step(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	body := &StepResponse{}
	code := 200
	msg := ""

	var err error
	var s bot.State

	if behaviorBuffer.Len() == 0 {
		code = -1
		msg = "not behavior data, need upload!"
		goto EXT
	}

	s = mbot.RunStep()

	body.Blackboard, err = mbot.GetMetadata()
	if err != nil {
		msg = err.Error()
		code = -3
	}
	body.Cur = mbot.GetCurNodeID()
	body.Prev = mbot.GetPrevNodeID()

	if s == bot.SEnd {
		mbot = nil
		behaviorBuffer.Reset()
		code = -4
		msg = "执行到末尾，重新上传行为文件开始！"
		goto EXT
	} else if s == bot.SBreak {
		body.Blackboard, _ = mbot.GetMetadata()
		code = -5
		msg = "执行遇到请求错误!"
		mbot = nil
		behaviorBuffer.Reset()
		goto EXT
	}

EXT:
	res.Code = code
	res.Msg = msg
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

var mbot *bot.Bot
var behaviorBuffer bytes.Buffer
var mockServer *mock.MockServer

func Route(e *echo.Echo) {

	mockServer = mock.NewServer()
	behaviorBuffer.Reset()

	e.POST("/upload", Upload)
	e.POST("/run", Run)
	e.POST("/step", Step)

}
