package script

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestFindOne(t *testing.T) {
	mgoMod := MgoModule{}

	L := lua.NewState()
	defer L.Close()

	err := L.DoFile("../json.lua")
	assert.Equal(t, err, nil)
	err = L.DoFile("../global.lua")
	assert.Equal(t, err, nil)

	L.PreloadModule("mgo", mgoMod.Loader)

	test1id := uuid.NewString()
	test2id := uuid.NewString()
	testdb := "testfindone"
	L.DoString(`test1id = "` + test1id + `"`)
	L.DoString(`test2id = "` + test2id + `"`)
	L.DoString(`testdb = "` + testdb + `"`)

	err = L.DoString(`
		local mgo = require("mgo")
		ret = mgo.conn("bot", "mongodb://127.0.0.1:27017")
		assert(ret == "succ", "mgo connect err " .. ret)

		ret = mgo.insert_one(testdb, {_id = test1id, msg = "aa"})
		assert(ret == "succ", "mgo insert_one err " .. ret)
		
		val, ret = mgo.find_one(testdb, {_id = test1id})
		assert(ret == "succ", "mgo find err " .. ret)

		lt = json.decode(val)
		table.print(lt)
		assert(lt._id == test1id, "find err " .. test1id)

		mgo.delete_one(testdb, {_id = test1id})
	`)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func TestFind(t *testing.T) {

	mgoMod := MgoModule{}

	L := lua.NewState()
	defer L.Close()

	err := L.DoFile("../json.lua")
	assert.Equal(t, err, nil)
	err = L.DoFile("../global.lua")
	assert.Equal(t, err, nil)

	L.PreloadModule("mgo", mgoMod.Loader)

	testdb := "testfind"
	L.DoString(`testdb = "` + testdb + `"`)

	err = L.DoString(`
		local mgo = require("mgo")
		ret = mgo.conn("bot", "mongodb://127.0.0.1:27017")
		assert(ret == "succ", "mgo connect err " .. ret)

		doc = {
			{a = 1, b = "a", c = false},
			{a = 2, b = "b", c = true},
			{a = { b = "c" }},
			{a = {"a", "b", "c", "d"}},
			{a = {1, 2, 3, 4}}
		}

		ret = mgo.insert_many(testdb, doc)
		assert(ret == "succ", "mgo insert many err " .. ret)

		val, ret = mgo.find(testdb, {})
		assert(ret == "succ", "mgo find err " .. ret)

		lt = json.decode(val)
		assert(lt[1].a == 1, "find err " .. val)

		ret = mgo.delete_many(testdb, {})
		assert(ret == "succ", "delete many err " .. ret)
	`)
	if err != nil {
		fmt.Println(err.Error())
	}

	t.Fail()
}

func TestInsertMany(t *testing.T) {
	mgoMod := MgoModule{}

	L := lua.NewState()
	defer L.Close()

	err := L.DoFile("../json.lua")
	assert.Equal(t, err, nil)
	err = L.DoFile("../global.lua")
	assert.Equal(t, err, nil)

	L.PreloadModule("mgo", mgoMod.Loader)

	testdb := "testinsertmany"
	L.DoString(`testdb = "` + testdb + `"`)

	err = L.DoString(`

		local mgo = require("mgo")
		ret = mgo.conn("bot", "mongodb://127.0.0.1:27017")
		assert(ret == "succ", "mgo connect err " .. ret)

		--ret = mgo.insert_many(testdb, {{a = 1},{a = 2}})
		--assert(ret == "succ", "mgo insert many err " .. ret)

		val, ret = mgo.find(testdb, {})
		assert(ret == "succ", "mgo find err " .. ret)

		lt = json.decode(val)
		assert(lt[1].a == 1, "insert many err " .. val)
		assert(lt[2].a == 2, "insert many err " .. val)
	`)

	if err != nil {
		fmt.Println(err.Error())
	}

}

func TestUpdateOne(t *testing.T) {
	mgoMod := MgoModule{}

	L := lua.NewState()
	defer L.Close()

	err := L.DoFile("../json.lua")
	assert.Equal(t, err, nil)
	err = L.DoFile("../global.lua")
	assert.Equal(t, err, nil)

	L.PreloadModule("mgo", mgoMod.Loader)

	testdb := "testupdateone"
	L.DoString(`testdb = "` + testdb + `"`)

	err = L.DoString(`
		local mgo = require("mgo")
		ret = mgo.conn("bot", "mongodb://127.0.0.1:27017")
		assert(ret == "succ", "mgo connect err " .. ret)

		ret = mgo.insert_one(testdb, {a = 1, b = "aa"})
		assert(ret == "succ", "mgo insert one err " .. ret)

		settable = {}
		settable["$set"] = {b = "bb"}
		mgo.update_one(testdb, {a = 1}, settable)

		val, ret = mgo.find_one(testdb, {a = 1})
		assert(ret == "succ", "mgo find err " .. ret)

		lt = json.decode(val)
		assert(lt.b == "bb", "update one err " .. val)

		ret = mgo.delete_many(testdb, {})
		assert(ret == "succ", "delete many err " .. ret)
	`)

	if err != nil {
		fmt.Println(err.Error())
	}
}
