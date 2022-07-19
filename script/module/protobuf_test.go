package script

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/pojol/gobot/script/book"
	lua "github.com/yuin/gopher-lua"
)

func mockServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

	return ts
}

func TestProtobufEncode(t *testing.T) {
	protomod := ProtoModule{}
	httpmod := NewHttpModule()
	L := lua.NewState()
	defer L.Close()

	ts := mockServer()
	defer ts.Close()

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

func TestProtobufDecode(t *testing.T) {
	protomod := ProtoModule{}
	httpmod := NewHttpModule()
	L := lua.NewState()
	defer L.Close()

	ts := mockServer()
	defer ts.Close()

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
