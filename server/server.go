package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pojol/gobot/behavior"
	"github.com/pojol/gobot/bot"
	"github.com/pojol/gobot/database"
	"github.com/pojol/gobot/factory"
	"github.com/pojol/gobot/utils"
)

type Err int32

const (
	Succ           Err = 200
	Fail           Err = 1000
	ErrContentRead     = 1000 + iota
	ErrWrongInput
	ErrJsonUnmarshal
	ErrJsonInvalid
	ErrPluginLoad
	ErrMetaData
	ErrEnd
	ErrBreak
	ErrCantFindBot
	ErrCreateBot
	ErrEmptyBatch
	ErrTagsFormat
	ErrRunningErr
	ErrUploadConfig
	ErrGetConfig
)

var errmap map[Err]string = map[Err]string{
	ErrContentRead:   "failed to read request content",
	ErrJsonInvalid:   "wrong file format",
	ErrJsonUnmarshal: "json unmarshal err",
	ErrWrongInput:    "bad request parameter",
	ErrPluginLoad:    "failed to plugin load",
	ErrMetaData:      "failed to get meta data",
	ErrEnd:           "run to the end",
	ErrBreak:         "run to the break",
	ErrCantFindBot:   "can't find bot",
	ErrCreateBot:     "failed to create bot, the behavior tree file needs to be uploaded to the server before creation",
	ErrEmptyBatch:    "empty batch info",
}

func FileBlobUpload(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	code := Succ

	name := ctx.Request().Header.Get("FileName")
	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	if len(bts) == 0 {
		code = ErrContentRead // tmp
		fmt.Println("bytes is empty!")
		goto EXT
	}

	_, err = behavior.New(bts)
	if err != nil {
		fmt.Println(err.Error())
		code = ErrJsonInvalid
		goto EXT
	}
	factory.Global.AddBehavior(name, bts)

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]

	ctx.JSON(http.StatusOK, res)
	return nil
}

func FileTextUpload(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	code := Succ
	var upload *utils.UploadFile
	var fbyte []byte
	var name string

	f, header, err := ctx.Request().FormFile("file")
	if err != nil {
		code = ErrContentRead
		fmt.Println(err.Error())
		goto EXT
	}

	upload = utils.NewUploadFile(f, header)
	fbyte = upload.ReadBytes()
	if len(fbyte) == 0 {
		code = ErrContentRead
		goto EXT
	}

	name = upload.FileName()
	_, err = behavior.New(fbyte)
	if err != nil {
		fmt.Println(err.Error())
		code = ErrJsonInvalid
		goto EXT
	}
	factory.Global.AddBehavior(name, fbyte)

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]

	ctx.JSON(http.StatusOK, res)
	return nil
}

func FileRemove(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}
	req := &FileRemoveReq{}
	body := &FileRemoveRes{}

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, req)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	factory.Global.RmvBehavior(req.Name)
	for _, v := range factory.Global.GetBehaviors() {
		body.Bots = append(body.Bots, behaviorInfo{
			Name:   v.Name,
			Update: v.UpdateTime,
		})
	}

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

func FileGetList(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}
	body := &BehaviorListRes{}

	info := factory.Global.GetBehaviors()
	for _, v := range info {
		tags := []string{}
		json.Unmarshal(v.TagDat, &tags)

		body.Bots = append(body.Bots, behaviorInfo{
			Name:   v.Name,
			Update: v.UpdateTime,
			Status: v.Status,
			Tags:   tags,
		})
	}

	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

func FileSetTags(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}
	req := &SetBehaviorTagsReq{}
	body := &SetBehaviorTagsRes{}
	var info []database.BehaviorInfo
	var jdat []byte

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, &req)
	if err != nil {
		code = ErrJsonInvalid
		fmt.Println(err.Error())
		goto EXT
	}

	jdat, err = json.Marshal(req.NewTags)
	if err != nil {
		code = ErrTagsFormat
		goto EXT
	}

	info = factory.Global.UpdateBehaviorTags(req.Name, jdat)
	for _, v := range info {
		tags := []string{}
		json.Unmarshal(v.TagDat, &tags)
		body.Bots = append(body.Bots, behaviorInfo{
			Name:   v.Name,
			Update: v.UpdateTime,
			Status: v.Status,
			Tags:   tags,
		})
	}

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

func ConfigUpload(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}

	name := ctx.Request().Header.Get("FileName")

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	if len(bts) == 0 {
		code = ErrContentRead // tmp
		fmt.Println("bytes is empty!")
		goto EXT
	}

	err = factory.Global.UploadConfig(name, bts)
	if err != nil {
		code = ErrUploadConfig
		res.Msg = err.Error()
		goto EXT
	}

EXT:

	res.Code = int(code)
	ctx.JSON(http.StatusOK, res)
	return nil
}

func ConfigRemove(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}
	req := ConfigRemoveReq{}

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, &req)
	if err != nil {
		code = ErrJsonInvalid
		fmt.Println(err.Error())
		goto EXT
	}

	err = factory.Global.RemoveConfig(req.Name)
	if err != nil {
		fmt.Println("remove config", err.Error())
		code = ErrGetConfig
		goto EXT
	}

EXT:
	res.Code = int(code)
	ctx.JSON(http.StatusOK, res)
	return nil
}

func ConfigListInfo(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}

	lst, err := factory.Global.GetConfigList()
	if err != nil {
		fmt.Println("config get list err", err.Error())
		code = ErrGetConfig
		goto EXT
	}

	res.Body = ConfigGetListInfoRes{
		Lst: lst,
	}

EXT:
	res.Code = int(code)
	ctx.JSON(http.StatusOK, res)
	return nil
}

func ConfigGetInfo(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	code := Succ
	res := &Response{}
	name := ctx.Request().Header.Get("FileName")

	cfg, err := factory.Global.GetConfig(name)
	if err != nil {
		fmt.Println("get config", err.Error())
		code = ErrGetConfig
		res.Msg = err.Error()
		goto EXT
	}

EXT:
	res.Code = int(code)
	ctx.Blob(http.StatusOK, "text/plain;charset=utf-8", cfg.Dat)
	return nil
}

func FileGetBlob(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	req := &FindBehaviorReq{}
	info := database.BehaviorInfo{}

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

EXT:
	ctx.Blob(http.StatusOK, "text/plain;charset=utf-8", info.Dat)
	return nil
}

func GetReport(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	body := &ReportRes{
		Info: factory.Global.GetReport(),
	}

	res.Code = int(Succ)
	res.Msg = ""
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

func BotRun(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	req := &BotRunRequest{}
	code := Succ
	var info database.BehaviorInfo
	var tree *behavior.Tree
	var b *bot.Bot

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead
		fmt.Println(err.Error())
		goto EXT
	}

	err = json.Unmarshal(bts, &req)
	if err != nil {
		code = ErrJsonUnmarshal
		fmt.Println(err.Error())
		goto EXT
	}

	if req.Name == "" {
		code = ErrWrongInput
		goto EXT
	}

	info, err = factory.Global.FindBehavior(req.Name)
	if err != nil {
		code = Fail
		goto EXT
	}

	tree, err = behavior.New(info.Dat)
	if err != nil {
		code = Fail
		goto EXT
	}
	b = bot.NewWithBehaviorTree("script/", tree, req.Name, 1, factory.Global.GetGlobalScript())
	err = b.RunByBlock()
	if err != nil {
		code = ErrRunningErr
		errmap[code] = err.Error()
	}

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]
	ctx.JSON(http.StatusOK, res)
	return nil
}

func BotCreateBatch(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	req := &BotBatchCreateRequest{}
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

	if req.Name == "" {
		goto EXT
	}
	if req.Num == 0 {
		goto EXT
	}

	err = factory.Global.AddTask(req.Name, int32(req.Num))
	if err != nil {
		res.Msg = err.Error()
	}

EXT:
	res.Code = int(code)
	res.Msg = errmap[code]

	ctx.JSON(http.StatusOK, res)
	return nil
}

func BotList(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	body := &BotListResponse{}
	code := Succ

	body.Lst = factory.Global.GetBatchInfo()
	res.Code = int(code)
	res.Msg = errmap[code]
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

func DebugStep(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	req := &StepRequest{}
	body := &StepResponse{}
	code := Succ
	var b *bot.Bot

	var err error
	var ok bool
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
	body.Blackboard, body.Change, body.RuntimeErr, ok = b.GetMetadata()
	if !ok {
		code = ErrMetaData
		goto EXT
	}
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

func DebugCreate(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	body := &CreateDebugBotResponse{}
	code := Succ
	var b *bot.Bot

	name := ctx.Request().Header.Get("FileName")
	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = ErrContentRead // tmp
		fmt.Println(err.Error())
		goto EXT
	}

	b = factory.Global.CreateDebugBot(name, bts)
	if b == nil {
		fmt.Println("create bot err", name)
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

func ReqPrint() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println(c.Path(), c.Request().Method, "enter")
			err := next(c)
			fmt.Println(c.Path(), c.Request().Method, err)
			return err
		}
	}
}

func Route(e *echo.Echo) {

	e.GET("/health", func(ctx echo.Context) error {
		ctx.JSONBlob(http.StatusOK, []byte(``))
		return nil
	})

	e.POST("/file.uploadTxt", FileTextUpload)
	e.POST("/file.uploadBlob", FileBlobUpload)
	e.POST("/file.remove", FileRemove)
	e.POST("/file.list", FileGetList)
	e.POST("/file.get", FileGetBlob)
	e.POST("/file.setTags", FileSetTags)

	e.POST("/config.get", ConfigGetInfo)
	e.POST("/config.list", ConfigListInfo)
	e.POST("/config.rmv", ConfigRemove)
	e.POST("/config.upload", ConfigUpload)

	e.POST("/bot.run", BotRun)
	e.POST("/bot.batch", BotCreateBatch) // 创建一批bot
	e.POST("/bot.list", BotList)

	e.POST("/debug.create", DebugCreate) // 创建一个 edit 中的bot 实例、
	e.POST("/debug.step", DebugStep)     // 单步运行 edit 中的bot

	e.POST("/report.info", GetReport)
}
