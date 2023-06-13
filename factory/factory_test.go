package factory

import (
	"testing"
	"time"

	"github.com/pojol/gobot/database"
	"github.com/pojol/gobot/mock"
)

var compose = `
<behavior>
  <id>20913145-5f7e-4b0c-babc-4e94e7c4d6ad</id>
  <ty>RootNode</ty>
  <pos>
    <x>0</x>
    <y>0</y>
  </pos>
  <children>
    <id>6258e521-4d7a-4427-a467-d102daf6ab9e</id>
    <ty>LoopNode</ty>
    <pos>
      <x>-5</x>
      <y>66</y>
    </pos>
    <children>
      <id>1291d5c2-5964-4b98-82d1-bb106f0e9c57</id>
      <ty>SelectorNode</ty>
      <pos>
        <x>-15</x>
        <y>125</y>
      </pos>
      <children>
        <id>743186ab-655b-47c8-a986-44ae49adf33b</id>
        <ty>ConditionNode</ty>
        <pos>
          <x>-70</x>
          <y>171</y>
        </pos>
        <children>
          <id>e3e32962-edcb-4dc1-add2-f934ff8bb87e</id>
          <ty>ActionNode</ty>
          <pos>
            <x>-85</x>
            <y>222</y>
          </pos>
          <code>
local parm = {
  body = {
      Token = meta.Token
  },    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local url = &#34;http://127.0.0.1:7777/login/guest&#34;
local http = require(&#34;http&#34;)

function execute()

  -- http post request
  res, errmsg = http.post(url, parm)
  if errmsg == nil then
    body = json.decode(res[&#34;body&#34;])
    merge(meta, body.Body)
  end
end
</code>
        </children>
        <code>

-- Write expression to return true or false
function execute()

    return meta.Token == &#34;&#34;

end
</code>
      </children>
      <children>
        <id>262cf484-5075-4971-885b-f9f74c9b1e92</id>
        <ty>ConditionNode</ty>
        <pos>
          <x>95</x>
          <y>171</y>
        </pos>
        <children>
          <id>01f48591-3b97-46b9-a86c-df10a4d009c6</id>
          <ty>ActionNode</ty>
          <pos>
            <x>80</x>
            <y>222</y>
          </pos>
          <children>
            <id>f7bbb512-d239-4d5a-881b-7c445a47abc7</id>
            <ty>SequenceNode</ty>
            <pos>
              <x>70</x>
              <y>302</y>
            </pos>
            <children>
              <id>2796bf91-a0f8-4be1-92fc-b164500b7cf0</id>
              <ty>ActionNode</ty>
              <pos>
                <x>15</x>
                <y>347</y>
              </pos>
              <code>
local parm = {
  body = {
      Token = meta.Token
  },    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local url = &#34;http://127.0.0.1:7777/base/hero.info&#34;
local http = require(&#34;http&#34;)

function execute()

  -- http post request
  res, errmsg = http.post(url, parm)
  if errmsg == nil then
    body = json.decode(res[&#34;body&#34;])
    merge(meta, body.Body)
  end
end
</code>
            </children>
            <children>
              <id>b620b162-f5d5-41b4-ab5f-65436792d4b4</id>
              <ty>WaitNode</ty>
              <pos>
                <x>100</x>
                <y>347</y>
              </pos>
              <wait>100</wait>
            </children>
            <children>
              <id>2eda6249-b555-42cb-b6b1-accce15c4f34</id>
              <ty>LoopNode</ty>
              <pos>
                <x>161</x>
                <y>342</y>
              </pos>
              <children>
                <id>3c6a6691-be4d-42b3-a909-451ab741309d</id>
                <ty>ActionNode</ty>
                <pos>
                  <x>166</x>
                  <y>419</y>
                </pos>
                <code>
local parm = {
  body = {
      Token = meta.Token,
      HeroID = &#34;joy&#34;
  },    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local url = &#34;http://127.0.0.1:7777/base/hero.lvup&#34;
local http = require(&#34;http&#34;)

function execute()

  -- http post request
  res, errmsg = http.post(url, parm)
  if errmsg == nil then
    body = json.decode(res[&#34;body&#34;])
    merge(meta, body.Body)
  end
end
</code>
              </children>
              <loop>2</loop>
            </children>
          </children>
          <code>
local parm = {
  body = {
      Token = meta.Token
  },    -- request body
  timeout = &#34;10s&#34;,
  headers = {},
}

local url = &#34;http://127.0.0.1:7777/base/acc.info&#34;
local http = require(&#34;http&#34;)

function execute()

  -- http post request
  res, errmsg = http.post(url, parm)
  if errmsg == nil then
    body = json.decode(res[&#34;body&#34;])
    merge(meta, body.Body)
  end
end
</code>
        </children>
        <code>

-- Write expression to return true or false
function execute()

    return meta.Token ~= &#34;&#34;

end
</code>
      </children>
    </children>
    <loop>3</loop>
  </children>
</behavior>

`

func TestLoop(t *testing.T) {

	ms := mock.NewServer()
	go ms.Start(":7777")

	defer ms.Close()

	f, err := Create(WithScriptPath("../script/"))
	if err != nil {
		panic(err)
	}

	database.GetBehavior().Upset("behavior", []byte(compose))

	for i := 0; i < 10; i++ {
		f.AddBatch("behavior", 0, 10)
	}

	time.Sleep(time.Second * 2)
}
