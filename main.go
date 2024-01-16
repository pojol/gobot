package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	_ "net/http/pprof"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pojol/braid-go"
	"github.com/pojol/braid-go/components"
	"github.com/pojol/braid-go/components/discoverk8s"
	"github.com/pojol/braid-go/module/meta"
	"github.com/pojol/gobot/constant"
	"github.com/pojol/gobot/factory"
	"github.com/pojol/gobot/mock"
	"github.com/pojol/gobot/server"
	"github.com/redis/go-redis/v9"
)

var (
	help bool

	dbmode       bool
	cluster      bool
	scriptPath   string
	openHttpMock bool
	openTcpMock  bool
)

const (
	// Version of gobot driver
	Version = "v0.4.0"

	banner = `
              __              __      
             /\ \            /\ \__   
   __     ___\ \ \____    ___\ \ ,_\  
 /'_ '\  / __'\ \ '__'\  / __'\ \ \  
/\ \L\ \/\ \L\ \ \ \L\ \/\ \L\ \ \ \_ 
\ \____ \ \____/\ \_,__/\ \____/\ \__\
 \/___L\ \/___/  \/___/  \/___/  \/__/
   /\____/                            
   \_/__/             %s                

`
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.BoolVar(&cluster, "cluster", false, "open cluster mode")
	flag.BoolVar(&dbmode, "no_database", false, "Run in local mode")
	flag.BoolVar(&openHttpMock, "httpmock", false, "open http mock server")
	flag.BoolVar(&openTcpMock, "tcpmock", false, "open tcp mock server")
	flag.StringVar(&scriptPath, "script_path", "script/", "Path to bot script")
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println("panic:", string(buf[:n]))
		}
	}()

	initFlag()
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	_, err := factory.Create(
		factory.WithNoDatabase(dbmode),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("open cluster mode", cluster)
	if cluster {
		constant.SetClusterState(true)

		redis_addr := os.Getenv("REDIS_ADDR")

		b, _ := braid.NewService(
			"bot",
			os.Getenv("POD_NAME"),
			&components.DefaultDirector{
				Opts: &components.DirectorOpts{
					RedisCliOpts: &redis.Options{
						Addr: redis_addr,
					},
					DiscoverOpts: []discoverk8s.Option{
						discoverk8s.WithNamespace("bot"),
						discoverk8s.WithSelectorTag("bot"),
					},
				},
			},
		)

		b.Init()
		b.Run()
		defer b.Close()

		statechan, err := braid.Topic(meta.TopicElectionChangeState).Sub(context.TODO(), "election"+uuid.NewString())
		if err != nil {
			panic(err)
		}
		defer statechan.Close()

		statechan.Arrived(func(msg *meta.Message) error {
			smsg := meta.DecodeStateChangeMsg(msg)
			constant.SetServerState(smsg.State)
			return nil
		})
	}

	fmt.Println("open http mock", openHttpMock)
	if openHttpMock {
		ms := mock.NewHttpServer()
		go ms.Start(":6666")
		defer ms.Close()
	}

	fmt.Println("open tcp mock", openTcpMock)
	if openTcpMock {
		tcpls := mock.StarTCPServer(":6667")
		defer tcpls.Close()
	}

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	// 查看有没有未完成的队列
	factory.Global.CheckTaskHistory()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	server.Route(e)
	e.Start(":8888")
	// Stop the service gracefully.
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
