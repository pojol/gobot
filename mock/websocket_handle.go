package mock

import (
	"math/rand"
	"sync"

	proto "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	sessionmap = make(map[string]*websocket.Conn)
	sessionMu  sync.RWMutex
)

func getSession(id string) *websocket.Conn {
	sessionMu.RLock()
	defer sessionMu.RUnlock()
	return sessionmap[id]
}

func setSession(id string, ws *websocket.Conn) {
	sessionMu.Lock()
	sessionmap[id] = ws
	sessionMu.Unlock()
}

func delSession(id string) {
	sessionMu.Lock()
	delete(sessionmap, id)
	sessionMu.Unlock()
}

func wsHelloHandle(conn *websocket.Conn) error {
	var dict = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	var msg string

	for i := 0; i < 3; i++ {
		msg += dict[rand.Intn(len(dict)-1)]
	}

	res := HelloRes{
		Message: msg,
	}

	byt, _ := proto.Marshal(&res)

	return conn.WriteMessage(Hello, byt)
}

func wsHeroInfoHandle(conn *websocket.Conn, msgBody []byte) error {

	req := &GetHeroInfoReq{}
	err := proto.Unmarshal(msgBody, req)
	if err != nil {
		return err
	}

	res := &GetHeroInfoRes{
		HeroInfo: &Hero{},
	}

	acc, err := getAccInfo(req.SessionID)
	if err != nil {
		return err
	}

	for _, v := range acc.Heros {
		if v.ID == req.HeroID {
			res.HeroInfo.ID = v.ID
			res.HeroInfo.Lv = v.Lv
		}
	}

	byt, _ := proto.Marshal(res)

	return conn.WriteMessage(HeroInfo, byt)
}

func wsHeroLvupHandle(conn *websocket.Conn, msgBody []byte) error {

	req := &HeroLvupReq{}
	err := proto.Unmarshal(msgBody, req)
	if err != nil {
		return err
	}

	acc, err := getAccInfo(req.SessionID)
	if err != nil {
		return err
	}

	for k := range acc.Heros {
		if acc.Heros[k].ID == req.HeroID {
			acc.Heros[k].Lv++
			break
		}
	}

	// response
	var heros []*Hero
	for _, v := range acc.Heros {
		heros = append(heros, &Hero{
			ID: v.ID,
			Lv: v.Lv,
		})
	}

	res := HeroLvupRes{
		AccInfo: &Acc{
			Heros:   heros,
			Diamond: acc.Diamond,
			Gold:    acc.Gold,
		},
	}

	byt, _ := proto.Marshal(&res)

	return conn.WriteMessage(HeroLvup, byt)
}
