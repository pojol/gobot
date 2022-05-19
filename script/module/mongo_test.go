package script

import (
	"fmt"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestFindOne(t *testing.T) {
	mgoMod := MgoModule{}

	L := lua.NewState()
	defer L.Close()

	L.DoFile("./json.lua")
	L.DoFile("./global.lua")

	L.PreloadModule("mgo", mgoMod.Loader)
	err := L.DoString(`
		local mgo = require("mgo")
		errmsg = mgo.conn("bot", "mongodb://127.0.0.1:27017")
		print(errmsg)

		-- errmsg = mgo.insert_one("test", {msg = "aa"})
		-- errmsg = mgo.insert_one("test", {msg = "bb"})
		-- print(errmsg)

		val, errmsg = mgo.find("test", {})
		print(val)
		if errmsg == "succ" then
			table.print(json.decode(val))
		else 
			print(errmsg)
		end
	`)
	if err != nil {
		fmt.Println(err.Error())
	}

	t.Fail()
}

func TestFind(t *testing.T) {

	mgoMod := MgoModule{}

	L := lua.NewState()
	defer L.Close()

	L.DoFile("./json.lua")
	L.DoFile("./global.lua")

	L.PreloadModule("mgo", mgoMod.Loader)
	err := L.DoString(`
		local mgo = require("mgo")
		errmsg = mgo.conn("bot", "mongodb://127.0.0.1:27017")
		print(errmsg)

		-- errmsg = mgo.insert_one("test", {msg = "aa"})
		-- errmsg = mgo.insert_one("test", {msg = "bb"})
		-- print(errmsg)

		val, errmsg = mgo.find("test", {})
		print(val)
		if errmsg == "succ" then
			table.print(json.decode(val))
		else 
			print(errmsg)
		end
	`)
	if err != nil {
		fmt.Println(err.Error())
	}

	t.Fail()
}
