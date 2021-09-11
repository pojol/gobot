package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pojol/apibot/behavior"
	"github.com/pojol/apibot/mock"
	"github.com/pojol/apibot/utils"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

var srv *httptest.Server

type guestRes struct {
	Token string
}

type infoRes struct {
	Diamond int32
	Gold    int32
}

func TestMain(m *testing.M) {

	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ioutil.ReadAll(req.Body)
		var byt []byte

		fmt.Println("http server recv ", req.RequestURI)

		if req.RequestURI == "/login/guest" {
			byt, _ = json.Marshal(guestRes{
				Token: "abcd",
			})
		} else if req.RequestURI == "/base/acc.info" {
			byt, _ = json.Marshal(infoRes{
				Diamond: 100,
				Gold:    100,
			})
		}

		w.Write(byt)
	}))
	defer srv.Close()

	os.Exit(m.Run())
}

type Metadata struct {
	Val string
}

var compose = `
<behavior>
  <id>bb49e4e6-a89d-419b-9517-d2fb9f9d6d11</id>
  <ty>RootNode</ty>
  <pos>
    <x>0</x>
    <y>0</y>
  </pos>
  <children>
    <id>e0d9de38-a927-4ebf-a88b-787942362564</id>
    <ty>LoopNode</ty>
    <pos>
      <x>-5</x>
      <y>47</y>
    </pos>
    <children>
      <id>5e9b233c-5e9d-4706-996c-da868e7af15a</id>
      <ty>SelectorNode</ty>
      <pos>
        <x>-15</x>
        <y>126</y>
      </pos>
      <children>
        <id>190affa9-a509-4de9-a23c-e3519af0a8ea</id>
        <ty>ConditionNode</ty>
        <pos>
          <x>-20</x>
          <y>179</y>
        </pos>
        <children>
          <id>67d3595e-2c76-4534-8734-4365434ae974</id>
          <ty>HTTPActionNode</ty>
          <pos>
            <x>-50</x>
            <y>245</y>
          </pos>
          <code>-- http request parm
local parm = {
  body = {},    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local cli = require(&#34;cli&#34;)

--[[ 
  Write to http client call
  like &gt;&gt; cli.post(&#34;url&#34;,&#34;api&#34;, parm)
]]--
function execute()
    url = mock .. "/login/guest"
    res, errmsg = cli.post(url, parm)
    print(url, errmsg)
    if errmsg == nil then
        merge(meta, json.decode(res["body"]))
    end

    table.print(meta)
end
</code>
        </children>
        <code>
-- Write expression to return true or false
function execute()

    return meta.Token == &#34;&#34;

end</code>
      </children>
      <children>
        <id>1caa5e2f-8c81-48e3-ba90-0c8229e60fad</id>
        <ty>ConditionNode</ty>
        <pos>
          <x>50</x>
          <y>179</y>
        </pos>
        <children>
          <id>8dde4c5e-4676-42d0-9a59-854c003244a6</id>
          <ty>HTTPActionNode</ty>
          <pos>
            <x>55</x>
            <y>245</y>
          </pos>
          <children>
            <id>1de3a303-e821-4b76-9e4a-889f3910321f</id>
            <ty>SequenceNode</ty>
            <pos>
              <x>0</x>
              <y>314</y>
            </pos>
            <children>
              <id>1ab45b16-9bfb-478f-b384-852ee9440e53</id>
              <ty>HTTPActionNode</ty>
              <pos>
                <x>-35</x>
                <y>382</y>
              </pos>
              <code>
-- http request parm
local parm = {
  body = {
      Token: &#34;abcd&#34;,
  },    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local cli = require(&#34;cli&#34;)

--[[ 
  Write to http client call
  like &gt;&gt; cli.post(&#34;url&#34;,&#34;api&#34;, parm)
]]--
function execute()

    cli.post(&#34;htts://127.0.0.1:8888/base/hero.info&#34;, parm)

end
</code>
            </children>
            <children>
              <id>4410b7fd-1b76-4d16-a269-5c3fef3cdc52</id>
              <ty>WaitNode</ty>
              <pos>
                <x>55</x>
                <y>382</y>
              </pos>
              <wait>100</wait>
            </children>
            <children>
              <id>6dd0dd90-c11c-44dc-a6b8-22ac82d149fe</id>
              <ty>LoopNode</ty>
              <pos>
                <x>105</x>
                <y>377</y>
              </pos>
              <children>
                <id>97952050-95c1-4ca2-be75-d131818ec6f7</id>
                <ty>HTTPActionNode</ty>
                <pos>
                  <x>105</x>
                  <y>455</y>
                </pos>
                <code>
-- http request parm
local parm = {
  body = {
      Token = &#34;abcd&#34;,
      HeroID = &#34;joy&#34;,
  },    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local cli = require(&#34;cli&#34;)

--[[ 
  Write to http client call
  like &gt;&gt; cli.post(&#34;url&#34;,&#34;api&#34;, parm)
]]--
function execute()

    cli.post(&#34;htts://127.0.0.1:8888/base/hero.lvup&#34;, parm)

end
</code>
              </children>
              <loop>5</loop>
            </children>
          </children>
          <code>-- http request parm
local parm = {
  body = {
      Token = &#34;abcd&#34;,
  },    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local cli = require(&#34;cli&#34;)

--[[
  Write to http client call
  like &gt;&gt; cli.post(&#34;url&#34;,&#34;api&#34;, parm)
]]--
function execute()

    cli.post(&#34;htts://127.0.0.1:8888/base/acc.info&#34;, parm)

end
</code>
        </children>
        <code>
-- Write expression to return true or false
function execute()

    return meta.Token ~= &#34;&#34;

end</code>
      </children>
    </children>
    <loop>5</loop>
  </children>
</behavior>

`

func TestLoad(t *testing.T) {

	var tree *behavior.Tree
	var bot *Bot

	srv := mock.NewServer()

	tree, err := behavior.New([]byte(compose))
	assert.Equal(t, err, nil)

	bot = NewWithBehaviorTree(tree, srv.Url())
	for i := 0; i < 10; i++ {
		bot.RunStep()
	}
}

func TestStep(t *testing.T) {
	/*
		err := plugins.Load("../plugins/json/json.so")
		assert.Equal(t, err, nil)

		srv := mock.NewServer()

		var tree *behavior.Tree
		var bot *Bot

		tree, err = behavior.New([]byte(compose))
		assert.Equal(t, err, nil)

		bot = NewWithBehaviorTree(tree, srv.Url())

		for i := 0; i < 30; i++ {
			bot.RunStep()
		}
	*/
}

var luastruct = `
meta = {
    name = "Michel",
    age  = 31,
    fly = false,
}

function condition()
    return meta.name == "joy"
end 
`

func TestScript(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(luastruct); err != nil {
		panic(err)
	}
	if err := L.DoString(`meta.name="joy"`); err != nil {
		panic(err)
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("condition"),
		NRet:    1,
		Protect: true,
	}); err != nil {
		panic(err)
	}
	ret := L.Get(-1) // returned value
	fmt.Println("condition ret", ret)
	L.Pop(1) // remove received value

	meta, err := utils.Table2Map(L.GetGlobal("meta").(*lua.LTable))
	if err != nil {
		panic(err)
	}

	fmt.Println(meta)
}
