package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type guestRes struct {
	Token string
}

type accInfoReq struct {
	Token string
}

type accInfoRes struct {
	Token   string
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
	Token string
	Heros []heroInfo
}

type lvupReq struct {
	Token  string
	HeroID string
}

type lvupRes struct {
	Token string
	Heros []heroInfo
}

type MockAcc struct {
	Token   string
	Heros   []heroInfo
	Diamond int32
	Gold    int32
}

type MockServer struct {
	srv *httptest.Server
}

var accmap map[string]*MockAcc
var acclock sync.RWMutex

var strlst = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}

func createAcc() *MockAcc {

	token := ""
	for i := 0; i < 6; i++ {
		token += strlst[rand.Intn(len(strlst))]
	}

	acc := &MockAcc{
		Token:   token,
		Diamond: 50,
		Gold:    100,
		Heros: []heroInfo{
			{ID: "joy", Lv: 1},
			{ID: "pojoy", Lv: 2},
		},
	}

	accmap[token] = acc

	return acc
}

func routeGuest(in []byte) []byte {
	acc := createAcc()
	byt, _ := json.Marshal(guestRes{
		Token: acc.Token,
	})

	return byt
}

func routeAccInfo(in []byte) ([]byte, bool) {
	accinfo := accInfoReq{}
	byt := []byte("{}")
	ret := true

	err := json.Unmarshal(in, &accinfo)
	if err != nil {
		byt = []byte(err.Error())
		ret = false
		goto ext
	}

	if _, ok := accmap[accinfo.Token]; ok {
		byt, _ = json.Marshal(accInfoRes{
			Token:   accinfo.Token,
			Diamond: accmap[accinfo.Token].Diamond,
			Gold:    accmap[accinfo.Token].Gold,
		})
	} else {
		fmt.Println("can't find acc", accinfo.Token)
		ret = false
	}
ext:
	return byt, ret
}

func routeHeroInfo(in []byte) ([]byte, bool) {
	byt := []byte("{}")
	ret := true

	heroinfo := heroInfoReq{}
	err := json.Unmarshal(in, &heroinfo)
	if err != nil {
		byt = []byte(err.Error())
		ret = false
		goto ext
	}

	if _, ok := accmap[heroinfo.Token]; ok {
		byt, _ = json.Marshal(heroInfoRes{
			Token: heroinfo.Token,
			Heros: accmap[heroinfo.Token].Heros,
		})
	} else {
		fmt.Println("can't find acc", heroinfo.Token)
		ret = false
	}
ext:
	return byt, ret
}

func routeHeroLvup(in []byte) ([]byte, bool) {
	byt := []byte("{}")
	ret := true

	lvup := lvupReq{}
	err := json.Unmarshal(in, &lvup)
	if err != nil {
		byt = []byte(err.Error())
		ret = false
		goto ext
	}

	if _, ok := accmap[lvup.Token]; ok {

		flag := false

		for k := range accmap[lvup.Token].Heros {
			if accmap[lvup.Token].Heros[k].ID == lvup.HeroID {
				accmap[lvup.Token].Heros[k].Lv++
				flag = true
				break
			}
		}

		if !flag {
			fmt.Println("can't find hero token", lvup.Token, "hero", lvup.HeroID)
			ret = false
			goto ext
		}

		byt, _ = json.Marshal(lvupRes{
			Token: lvup.Token,
			Heros: accmap[lvup.Token].Heros,
		})

	} else {
		fmt.Println("can't find acc", lvup.Token)
		ret = false
	}

ext:
	return byt, ret
}

func mockRoute(w http.ResponseWriter, req *http.Request) {

	reqbyt, _ := ioutil.ReadAll(req.Body)
	var byt []byte
	var ok bool

	//fmt.Println("http server recv ", req.RequestURI)
	acclock.Lock()
	defer acclock.Unlock()

	if req.RequestURI == "/login/guest" {
		ok = true
		byt = routeGuest(reqbyt)
	} else if req.RequestURI == "/base/acc.info" {
		byt, ok = routeAccInfo(reqbyt)
	} else if req.RequestURI == "/base/hero.info" {
		byt, ok = routeHeroInfo(reqbyt)
	} else if req.RequestURI == "/base/hero.lvup" {
		byt, ok = routeHeroLvup(reqbyt)
	}

	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Println("write", ok)
	w.Write(byt)
}

func NewServer() *MockServer {

	accmap = make(map[string]*MockAcc)
	rand.Seed(time.Now().UnixNano())

	ms := &MockServer{}
	// mock server
	ms.srv = httptest.NewServer(http.HandlerFunc(mockRoute))

	return ms
}

func (ms *MockServer) Url() string {
	return ms.srv.URL
}

func (ms *MockServer) Reset(token string) {
	delete(accmap, token)
}

func (ms *MockServer) Close() {
	ms.srv.Close()
}
