package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/pojol/apibot/assert"
	"github.com/pojol/apibot/behavior"
)

var srv *httptest.Server

func TestMain(m *testing.M) {

	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		parm := GetAccountInfoParam{}

		body, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(body, &parm)

		w.Write([]byte(parm.Token))
	}))
	defer srv.Close()

	os.Exit(m.Run())
}

type Metadata struct {
	Val string
}

type GetAccountInfoParam struct {
	Token string
}
type GetAccountPack struct {
}

func (p *GetAccountPack) Marshal(meta interface{}, param interface{}) []byte {

	byt, err := json.Marshal(&param)
	if err != nil {
		fmt.Println(err.Error())
	}

	return byt
}

func (p *GetAccountPack) Unmarshal(meta interface{}, body []byte, header http.Header) {
	mp := meta.(*Metadata)
	mp.Val = string(body)
}

type GetAccountAssert struct {
}

func (p *GetAccountAssert) Do(meta interface{}) error {
	mp := meta.(*Metadata)
	return assert.Equal(mp.Val, "aabbcc", reflect.TypeOf(*p).Name())
}

var compose = `
{
	"behavior": "post",
	"url": "",
	"name": "",
	"script" : "GetAccountPack",
	"param" : {
		"Token" : "aabbcc" 
	}
}
`

var structmap map[string]interface{}

type composeInfo struct {
	Behavior string      `json:"behavior"`
	Url      string      `json:"url"`
	Name     string      `json:"name"`
	Script   string      `json:"script"`
	Param    interface{} `json:"param"`
}

func TestBot(t *testing.T) {

	md := Metadata{}
	b := New("", &md)

	structmap = make(map[string]interface{})
	structmap["GetAccountPack"] = &GetAccountPack{}

	info := &composeInfo{}

	err := json.Unmarshal([]byte(compose), info)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

	if info.Behavior == "post" {

		pack := reflect.New(reflect.TypeOf(structmap[info.Script]).Elem())

		b.Post(&behavior.HTTPPost{
			URL:   srv.URL,
			Meta:  b.metadata,
			Param: info.Param,
			Api:   pack.Interface(),
		})
	}

	b.Run()

}
