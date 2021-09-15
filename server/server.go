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
	"github.com/pojol/apibot/utils"
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

func UploadWithBlob(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	code := Succ
	var name string
	var tree *behavior.Tree

	name = ctx.Request().Header.Get("FileName")
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
	factory.Global.AddBehavior(tree.ID, name, bts)

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]

	ctx.JSON(http.StatusOK, res)
	return nil
}

func UploadWithFile(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	code := Succ
	var upload *utils.UploadFile
	var fbyte []byte
	var name string
	var tree *behavior.Tree

	f, header, err := ctx.Request().FormFile("file")
	if err != nil {
		code = ErrContentRead
		fmt.Println(err.Error())
		goto EXT
	}

	upload = utils.NewUploadFile(f, header)
	if upload.GetFileExt() != ".xml" {
		code = ErrJsonInvalid
		goto EXT
	}
	fbyte = upload.ReadBytes()
	if len(fbyte) == 0 {
		code = ErrContentRead
		goto EXT
	}

	name = upload.FileName()
	tree, err = behavior.New(fbyte)
	if err != nil {
		fmt.Println(err.Error())
		code = ErrJsonInvalid
		goto EXT
	}
	factory.Global.AddBehavior(tree.ID, name, fbyte)

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]

	ctx.JSON(http.StatusOK, res)
	return nil
}

type behaviorInfo struct {
	Name   string
	Update int64
}
type BehaviorListRes struct {
	Bots []behaviorInfo
}

func GetList(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}
	body := &BehaviorListRes{}

	info := factory.Global.GetBehaviors()
	for _, v := range info {
		body.Bots = append(body.Bots, behaviorInfo{
			Name:   v.Name,
			Update: v.UpdateTime,
		})
	}

	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

type FindBehaviorReq struct {
	Name string
}

type FindBehaviorRes struct {
	Info factory.BehaviorInfo
}

func GetBlob(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	req := &FindBehaviorReq{}
	info := factory.BehaviorInfo{}

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, &req)
	if err != nil {
		fmt.Println(err.Error())
		goto EXT
	}

	info, err = factory.Global.FindBehavior(req.Name)
	if err != nil {
		fmt.Println(err.Error())
		goto EXT
	}

	fmt.Println("get blob", info.Name, info.RootID)

EXT:
	ctx.Blob(http.StatusOK, "text/plain;charset=utf-8", info.Dat)
	return nil
}

type RunRequest struct {
	Info factory.BatchInfo
}

type RunResponse struct {
}

func Run(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	req := &RunRequest{}
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

	for _, v := range req.Info.Batch {
		if v.Behavior == "" {
			goto EXT
		}
		if v.Num == 0 {
			goto EXT
		}
	}

	factory.Global.Append(req.Info)
	factory.Global.RunBatch()

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]

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
	Name string
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

	b = factory.Global.CreateBot(req.Name)
	if b == nil {
		fmt.Println("create bot err", req.Name)
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

	e.POST("/upload.blob", UploadWithBlob) // 上传行为树模版文件
	e.POST("/upload.file", UploadWithFile)

	e.POST("/get.list", GetList)
	e.POST("/get.blob", GetBlob)
	e.POST("/create", Create) // 创建一个bot
	e.POST("/run", Run)       // 运行bot
	e.POST("/step", Step)     // 单步运行一个bot

}
