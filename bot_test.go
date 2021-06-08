package bot

import (
	"encoding/json"
	"fmt"
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
		w.Write([]byte("123"))
	}))
	defer srv.Close()

	os.Exit(m.Run())
}

type Metadata struct {
	Val string
}

type APIGetAccountInfo struct {
	Meta *Metadata

	Token string
}

func (p *APIGetAccountInfo) Marshal() []byte {

	byt, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}

	return byt
}

func (p *APIGetAccountInfo) Unmarshal(body []byte, header http.Header) {

	p.Meta.Val = string(body)

}

func (p *APIGetAccountInfo) Assert() error {
	return assert.Equal(p.Meta.Val, "123", reflect.TypeOf(*p).Name())
}

func TestBot(t *testing.T) {

	md := Metadata{}
	b := New("", &md)

	b.Post(&behavior.HTTPPost{
		URL: srv.URL,
		Api: &APIGetAccountInfo{
			Meta:  &md,
			Token: "111",
		},
	})

	b.Run()

}
