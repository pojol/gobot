package script

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pojol/gobot/driver/script/book"
	lua "github.com/yuin/gopher-lua"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqbyt, _ := ioutil.ReadAll(req.Body)
		var byt []byte

		fmt.Println("http server recv ", req.RequestURI)

		if req.RequestURI == "/book" {
			var msg book.AddressBook
			err := proto.Unmarshal(reqbyt, &msg)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println("book", msg)

		} else if req.RequestURI == "/unmarshal" {
			var err error
			phones := []*book.Person_PhoneNumber{}
			phones = append(phones, &book.Person_PhoneNumber{
				Number: "666",
				Type:   book.Person_HOME,
			})
			person := []*book.Person{}
			person = append(person, &book.Person{
				Name:   "Joy",
				Id:     222,
				Email:  "joy@outlook.com",
				Phones: phones,
			})

			body := book.AddressBook{
				People: person,
			}

			byt, err = proto.Marshal(&body)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		w.Write(byt)
	}))
	defer ts.Close()

	m.Run()
}

func TestProtobufEncode(t *testing.T) {
	protomod := ProtoModule{}
	httpmod := NewHttpModule(&http.Client{})
	L := lua.NewState()
	defer L.Close()

	L.DoFile("./json.lua")
	L.DoFile("./global.lua")
	L.PreloadModule("http", httpmod.Loader)
	L.PreloadModule("proto", protomod.Loader)
	L.SetGlobal("url", lua.LString(ts.URL))

	err := L.DoString(`
		local proto = require("proto")
		local http = require("http")

		local person = {
			name = "joy",
			id = 111,
			email = "joy@outlook.com",
			phones = { {number = "555", type = 2} }
		}

		local book = {
			people = { person }
		}

		local parm = {
			timeout = 1000,
			headers = {},
		}
		
		posturl = url .. "/book"
		print("post : " .. posturl)

		parm.body, errmsg = proto.marshal("AddressBook", json.encode(book))
		if errmsg ~= "" then
			print(errmsg)
		end
		http.post(posturl, parm)

	`)
	if err != nil {
		fmt.Println(err.Error())
	}

}

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

func TestProtobufDecode(t *testing.T) {
	protomod := ProtoModule{}
	httpmod := NewHttpModule(&http.Client{})
	L := lua.NewState()
	defer L.Close()

	L.DoFile("./json.lua")
	L.DoFile("./global.lua")
	L.PreloadModule("http", httpmod.Loader)
	L.PreloadModule("proto", protomod.Loader)
	L.SetGlobal("url", lua.LString(ts.URL))

	err := L.DoString(`

		local proto = require("proto")
		local http = require("http")

		res, err = http.post(url .. "/unmarshal", {})
		print(res, err)

		body, err = proto.unmarshal("AddressBook", res["body"])
		print(err)

		print("lua table", json.encode(body))
	`)
	if err != nil {
		fmt.Println(err.Error())
	}
}
