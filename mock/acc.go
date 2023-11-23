package mock

import (
	"fmt"
	"sync"
)

type heroInfo struct {
	ID string
	Lv int32
}

type MockAcc struct {
	Token   string
	Heros   []heroInfo
	Diamond int32
	Gold    int32
}

var accmap = sync.Map{}

func createAcc(token string) *MockAcc {

	acc := &MockAcc{
		Token:   token,
		Diamond: 50,
		Gold:    100,
		Heros: []heroInfo{
			{ID: "joy", Lv: 1},
			{ID: "pojoy", Lv: 2},
		},
	}

	accmap.Store(acc.Token, acc)
	return acc
}

func getAccInfo(token string) (*MockAcc, error) {
	mapval, ok := accmap.Load(token)
	if !ok {
		return nil, fmt.Errorf("token not found %v", token)
	}

	accPtr, _ := mapval.(*MockAcc)
	return accPtr, nil
}

func setHeroLv(token string, heroID string) (*MockAcc, error) {
	mapval, ok := accmap.Load(token)
	if !ok {
		return nil, fmt.Errorf("token not found %v", token)
	}

	accPtr, _ := mapval.(*MockAcc)

	for k := range accPtr.Heros {
		if accPtr.Heros[k].ID == heroID {
			accPtr.Heros[k].Lv++
			break
		}
	}

	return accPtr, nil
}
