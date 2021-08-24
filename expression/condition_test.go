package expression

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
)

/*
	列表 [  ]
	键值 {  }
	表达式 $symbol : {}
*/

type MetaMarket struct {
	ID string
}

type Meta struct {
	Token   string
	Diamond int32
	Gold    int32
	Ticket  int32
	Market  *MetaMarket
}

func TestParse(t *testing.T) {

	lst := []string{
		`$eq : { Token : "" }`,
		`$eq : {Market.ID : "aabb"}`,
		`$and : [$ne:{Token:""}, $gt:{Gold:100}}] `,
		`$or : [$gt:{Gold: 100}, $gt:{Diamond:100}]`,
		`$and : [ $eq:{Token:"aabb"}, $or:[$gt:{Diamond:100}, $gt:{Ticket:100}] ]`,
	}

	var metalst []map[string]interface{}

	// 模拟服务返回，映射到 Meta 结构
	m1 := make(map[string]interface{})
	json.Unmarshal([]byte(`{"Token":""}`), &m1)
	metalst = append(metalst, m1)

	m2 := make(map[string]interface{})
	json.Unmarshal([]byte(`{"Token":"", "Market":{ "ID": "aabb"}}`), &m2)
	metalst = append(metalst, m2)

	m3 := make(map[string]interface{})
	json.Unmarshal([]byte(`{"Token":"aabb", "Gold":200}`), &m3)
	metalst = append(metalst, m3)
	metalst = append(metalst, m3)

	m4 := make(map[string]interface{})
	json.Unmarshal([]byte(`{"Token":"aabb", "Gold":200, "Diamond":150}`), &m3)
	metalst = append(metalst, m4)

	meta := make(map[string]interface{})

	for k, v := range lst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)

		mergo.MergeWithOverwrite(&meta, metalst[k])
		fmt.Println(k, meta)
		assert.Equal(t, true, eg.DecideMap(meta))
	}

	metas := []Meta{
		{
			Token: "",
		},
		{
			Market: &MetaMarket{ID: "aabb"},
		},
		{ //$and : [$ne:{meta.Token:""}, $gt:{meta.Gold:100}}]
			Token: "aabb",
			Gold:  200,
		},
		{
			Diamond: 150,
		},
		{ //$and : [ $eq:{meta.Token:"aabb"}, $or:[$gt:{meta.Diamond:100}, $gt:{meta.Ticket:100}] ]
			Token:   "aabb",
			Diamond: 150,
		},
	}

	for k, v := range lst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)

		assert.Equal(t, true, eg.Decide(metas[k]))
	}

}
