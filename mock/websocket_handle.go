package mock

import (
	"bytes"
	"encoding/binary"
	"math/rand"

	proto "github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func wsGuestHandle(conn *websocket.Conn) error {
	var heros []*Hero

	acc := createAcc(uuid.NewString())
	for _, v := range acc.Heros {
		heros = append(heros, &Hero{
			ID: v.ID,
			Lv: v.Lv,
		})
	}

	res := LoginGuestRes{
		AccInfo: &Acc{
			Heros:   heros,
			Diamond: acc.Diamond,
			Gold:    acc.Gold,
		},
		SessionID: acc.SessionID,
	}

	byt, _ := proto.Marshal(&res)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(LoginGuest))
	binary.Write(buf, binary.LittleEndian, byt)
	return conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
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
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(Hello))
	binary.Write(buf, binary.LittleEndian, byt)
	return conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
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
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(HeroInfo))
	binary.Write(buf, binary.LittleEndian, byt)
	return conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
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

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(HeroLvup))
	binary.Write(buf, binary.LittleEndian, byt)
	return conn.WriteMessage(websocket.BinaryMessage, byt)
}
