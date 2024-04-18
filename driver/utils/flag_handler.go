package utils

import (
	"flag"
	"fmt"
)

type Flag struct {
	Help bool

	Cluster      bool
	ScriptPath   string
	OpenHttpMock int
	OpenTcpMock  int
	OpenWSMock   int
	OpenAPIPort  int
	DBType       string
}

func InitFlag() *Flag {
	flag.BoolVar(&f.Help, "h", false, "this help")

	flag.BoolVar(&f.Cluster, "cluster", false, "open cluster mode")
	flag.StringVar(&f.DBType, "db", "sqlite", "The database type, defaulting to sqlite, can also be configured to mysql.")
	flag.IntVar(&f.OpenHttpMock, "httpmock", 0, "open http mock server")
	flag.IntVar(&f.OpenTcpMock, "tcpmock", 0, "open tcp mock server")
	flag.IntVar(&f.OpenWSMock, "websocketmock", 0, "open websocket mock server")
	flag.IntVar(&f.OpenAPIPort, "apiport", 0, "open api server port, defautl is 8888")
	flag.StringVar(&f.ScriptPath, "script_path", "script/", "Path to bot script")

	return f
}

var f = &Flag{}

func ShowUseage() bool {
	if f.Help {
		flag.Usage()

		fmt.Println("\nexample: ./main -db sqlite -httpmock 6666  // use sqlite as database, open http mock server on port 6666")
		return true
	}
	return false
}
