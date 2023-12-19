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
	SessionID string // 临时处理
	Heros     []heroInfo
	Diamond   int32
	Gold      int32
}

var accmap = sync.Map{}

func createAcc(id string) *MockAcc {

	acc := &MockAcc{
		SessionID: id,
		Diamond:   50,
		Gold:      100,
		Heros: []heroInfo{
			{ID: "joy", Lv: 1},
			{ID: "pojoy", Lv: 2},
		},
	}

	accmap.Store(acc.SessionID, acc)
	return acc
}

func getAccInfo(id string) (*MockAcc, error) {
	mapval, ok := accmap.Load(id)
	if !ok {
		return nil, fmt.Errorf("sessionid not found %v", id)
	}

	accPtr, _ := mapval.(*MockAcc)
	return accPtr, nil
}

func setHeroLv(id string, heroID string) (*MockAcc, error) {
	mapval, ok := accmap.Load(id)
	if !ok {
		return nil, fmt.Errorf("sessionid not found %v", id)
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
