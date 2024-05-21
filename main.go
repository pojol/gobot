package main

import (
	_ "net/http/pprof"

	"github.com/pojol/gobot/driver"
)

func main() {
	driver.Run()
}
