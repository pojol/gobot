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
		api := APIGetAccountInfo{}

		body, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(body, &api)

		w.Write([]byte(api.Token))
	}))
	defer srv.Close()

	os.Exit(m.Run())
}

type Metadata struct {
	Val string
}

type APIGetAccountInfo struct {
	Token string
}

func (p *APIGetAccountInfo) Marshal(meta interface{}) []byte {

	byt, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}

	return byt
}

func (p *APIGetAccountInfo) Unmarshal(meta interface{}, body []byte, header http.Header) {
	mp := meta.(*Metadata)
	mp.Val = string(body)
}

func (p *APIGetAccountInfo) Assert(meta interface{}) error {
	mp := meta.(*Metadata)
	return assert.Equal(mp.Val, "111", reflect.TypeOf(*p).Name())
}

func TestBot(t *testing.T) {

	md := Metadata{}
	b := New("", &md)

	b.Post(&behavior.HTTPPost{
		URL:  srv.URL,
		Meta: &md,
		Api: &APIGetAccountInfo{
			Token: "111",
		},
	})

	b.Run()

}
