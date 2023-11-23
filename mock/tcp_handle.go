package mock

import (
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
)

const (
	LoginGuest = 1001 // 模拟登录
	Hello      = 1002 // 模拟获取用户信息
	HeroInfo   = 1003 // 模拟获取角色信息
	HeroLvup   = 1004 // 模拟修改角色等级
)

func tcpRouteGuestHandle(conn *net.TCPConn, msgBody []byte) error {
	var heros []*Hero
	rand.Seed(time.Now().UnixNano())

	f, _ := conn.File()
	acc := createAcc(strconv.Itoa(int(f.Fd())))
	for _, v := range acc.Heros {
		heros = append(heros, &Hero{
			ID: v.ID,
			Lv: v.Lv,
		})
	}

	res := LoginGuestRes{
		AccInfo: &Acc{
			Token:   acc.Token,
			Heros:   heros,
			Diamond: acc.Diamond,
			Gold:    acc.Gold,
		},
	}

	byt, _ := proto.Marshal(&res)
	writeMsg(conn, LoginGuest, []byte{}, byt)

	return nil
}

func tcpHelloHandle(conn *net.TCPConn, msgBody []byte) error {
	var dict = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	var msg string

	for i := 0; i < 3; i++ {
		msg += dict[rand.Intn(len(dict)-1)]
	}

	res := HelloRes{
		Message: msg,
	}

	byt, _ := proto.Marshal(&res)
	writeMsg(conn, Hello, []byte{}, byt)

	return nil
}

func tcpHeroInfoHandle(conn *net.TCPConn, msgBody []byte) error {

	f, _ := conn.File()
	acc, err := getAccInfo(strconv.Itoa(int(f.Fd())))
	if err != nil {
		return err
	}

	req := &GetHeroInfoReq{}
	err = proto.Unmarshal(msgBody, req)
	if err != nil {
		return err
	}

	res := &GetHeroInfoRes{}

	for _, v := range acc.Heros {
		if v.ID == req.HeroID {
			res.HeroInfo.ID = v.ID
			res.HeroInfo.Lv = v.Lv
		}
	}

	byt, _ := proto.Marshal(res)
	writeMsg(conn, HeroInfo, []byte{}, byt)

	return nil
}

func tcpHeroLvupHandle(conn *net.TCPConn, msgBody []byte) error {

	f, _ := conn.File()
	acc, err := getAccInfo(strconv.Itoa(int(f.Fd())))
	if err != nil {
		return err
	}

	req := &HeroLvupReq{}
	err = proto.Unmarshal(msgBody, req)
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
			Token:   acc.Token,
			Heros:   heros,
			Diamond: acc.Diamond,
			Gold:    acc.Gold,
		},
	}

	byt, _ := proto.Marshal(&res)
	writeMsg(conn, HeroLvup, []byte{}, byt)

	return nil
}
