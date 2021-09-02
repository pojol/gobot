package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type guestRes struct {
	Token string
}

type accInfoReq struct {
	Token string
}

type accInfoRes struct {
	Diamond int32
	Gold    int32
}

type heroInfoReq struct {
	Token string
}

type heroInfo struct {
	ID string
	Lv int32
}

type heroInfoRes struct {
	Heros []heroInfo
}

type lvupReq struct {
	Token  string
	HeroID string
}

type lvupRes struct {
	Heros []heroInfo
}

const ACCTOKEN = "abcd"

var heros []heroInfo

type MockServer struct {
	srv *httptest.Server
}

func NewServer() *MockServer {

	ms := &MockServer{}
	// mock server
	ms.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		reqbyt, _ := ioutil.ReadAll(req.Body)
		var byt []byte

		fmt.Println("http server recv ", req.RequestURI)

		if req.RequestURI == "/login/guest" {
			byt, _ = json.Marshal(guestRes{
				Token: ACCTOKEN,
			})
		} else if req.RequestURI == "/base/acc.info" {
			accinfo := accInfoReq{}
			err := json.Unmarshal(reqbyt, &accinfo)
			if err == nil && accinfo.Token == ACCTOKEN {
				byt, _ = json.Marshal(accInfoRes{
					Diamond: 100,
					Gold:    100,
				})
			} else {
				byt = []byte("error")
			}

		} else if req.RequestURI == "/base/hero.info" {

			heroinfo := heroInfoReq{}
			err := json.Unmarshal(reqbyt, &heroinfo)
			if err == nil && heroinfo.Token == ACCTOKEN {
				byt, _ = json.Marshal(heroInfoRes{
					Heros: heros,
				})
			} else {
				byt = []byte("error")
			}

		} else if req.RequestURI == "/base/hero.lvup" {
			lvup := lvupReq{}
			err := json.Unmarshal(reqbyt, &lvup)
			if err == nil && lvup.Token == ACCTOKEN {
				for k := range heros {
					if heros[k].ID == lvup.HeroID {
						heros[k].Lv++
						byt, _ = json.Marshal(lvupRes{
							Heros: heros,
						})
						break
					}
				}
			} else {
				byt = []byte("error")
			}
		}

		w.Write(byt)
	}))

	return ms
}

func (ms *MockServer) Url() string {
	return ms.srv.URL
}

func (ms *MockServer) Reset() {
	heros = heros[:0]
	heros = append(heros, heroInfo{ID: "pojol", Lv: 1})
	heros = append(heros, heroInfo{ID: "joy", Lv: 1})
}

func (ms *MockServer) Close() {
	ms.srv.Close()
}
