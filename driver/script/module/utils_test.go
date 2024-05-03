package script

import (
	"math/rand"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func TestUtilsModule(t *testing.T) {
	utilsMod := UtilsModule{}
	rand.Seed(time.Now().UnixNano())

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("utils", utilsMod.Loader)
	L.DoFile("./global.lua")

	L.DoString(`
		local utils = require("utils")
		
		print("uuid", utils.uuid())
		print("random", utils.random(100))

		meta = {
			Token = "",
			Info = "",      -- debug log [info]
			Err = "",       -- debug log [err]
			Warn = "",      -- debug log [warn]
		}

		table.print(meta)
	`)
}
