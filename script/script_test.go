package script

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/pojol/gobot-driver/behavior"
	"github.com/pojol/gobot-driver/gpb"
	lua "github.com/yuin/gopher-lua"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqbyt, _ := ioutil.ReadAll(req.Body)
		var byt []byte

		fmt.Println("http server recv ", req.RequestURI)

		if req.RequestURI == "/item" {

			var msg gpb.PItem
			err := proto.Unmarshal(reqbyt, &msg)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println("item", msg)

		} else if req.RequestURI == "/card" {
			var msg gpb.PCard
			err := proto.Unmarshal(reqbyt, &msg)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println("card", msg)

		} else if req.RequestURI == "/unmarshal" {
			body := gpb.PItem{
				Id:  1001,
				Num: 100,
				Lv:  1,
			}
			byt, _ = proto.Marshal(&body)
		}

		w.Write(byt)
	}))
	defer ts.Close()

	m.Run()
}

func TestProtobufEncode(t *testing.T) {
	protomod := behavior.ProtoModule{}
	httpmod := behavior.NewHttpModule(&http.Client{})
	L := lua.NewState()
	defer L.Close()

	L.DoFile("./json.lua")
	L.DoFile("./global.lua")
	L.PreloadModule("cli", httpmod.Loader)
	L.PreloadModule("proto", protomod.Loader)
	L.SetGlobal("url", lua.LString(ts.URL))

	err := L.DoString(`
		local proto = require("proto")
		local cli = require("cli")

		local item = {
			Id = 101,
			Num = 10,
			InsId = "hello,item",
			Lv = 1,
		}

		local parm = {
			timeout = 1000,
			headers = {},
		}
		
		posturl = url .. "/item"
		print("post : " .. posturl)

		parm.body = proto.marshal("PItem", json.encode(item))
		cli.post(posturl, parm)

	`)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = L.DoString(`
		local proto = require("proto")
		local cli = require("cli")

		local item = {
			Id = 101,
			Num = 10,
			InsId = "hello,item",
			Lv = 1,
		}

		local card = {
			Id = "201",
			Timeout = 1000,
			Items = { item }
		}

		local parm = {
			timeout = 1000,
			headers = {},
		}

		posturl = url .. "/card"
		print("post : " .. posturl)

		parm.body = proto.marshal("PCard", json.encode(card))
		cli.post(posturl, parm)

	`)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestProtobufDecode(t *testing.T) {
	protomod := behavior.ProtoModule{}
	httpmod := behavior.NewHttpModule(&http.Client{})
	L := lua.NewState()
	defer L.Close()

	L.DoFile("./json.lua")
	L.DoFile("./global.lua")
	L.PreloadModule("cli", httpmod.Loader)
	L.PreloadModule("proto", protomod.Loader)
	L.SetGlobal("url", lua.LString(ts.URL))

	err := L.DoString(`

		local proto = require("proto")
		local cli = require("cli")

		res, err = cli.post(url .. "/unmarshal", {})
		print(err)

		body, err = proto.unmarshal("PItem", res["body"])
		print(err)

		print("lua table", json.encode(body))
	`)
	if err != nil {
		fmt.Println(err.Error())
	}
}
