package expression

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
)

/*
	值
		* true
		* false
		* float32
		* int
		* string  这里要考虑到中文
	键值 { : }
	表达式 $symbol : {}
	列表 [  ]	只能包含表达式
*/

type MetaMarket struct {
	ID string
}

type Meta struct {
	Token   string
	Diamond int32
	Gold    int32
	Ticket  int32
	Market  *MetaMarket // 这里需要使用指针值，不然没办法使用反射向下传递
}

func TestParse(t *testing.T) {

	lst := []string{
		`$eq : { Token : '' }`,
		`$eq : {Market.ID : 'aabb'}`,
		`$and : [$ne:{Token:''}, $gt:{Gold:101}}] `,
		`$or : [$gt:{Gold: 100}, $gt:{Diamond:100}]`,
		`$and : [ $eq:{Token:'aabb'}, $or:[$gt:{Diamond:100}, $gt:{Ticket:100}] ]`,
		`$eq : { Token : '中文' }`,
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

	m5 := make(map[string]interface{})
	json.Unmarshal([]byte(`{"Token":"中文"}`), &m5)
	metalst = append(metalst, m5)

	meta := make(map[string]interface{})

	for k, v := range lst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)

		mergo.MergeWithOverwrite(&meta, metalst[k])
		assert.Equal(t, true, eg.DecideWithMap(meta))
	}

	metas := []Meta{
		{
			Token: "",
		},
		{
			Market: &MetaMarket{ID: "aabb"},
		},
		{
			Token: "aabb",
			Gold:  200,
		},
		{
			Diamond: 150,
		},
		{
			Token:   "aabb",
			Diamond: 150,
		},
		{
			Token: "中文",
		},
	}

	for k, v := range lst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)

		fmt.Println(k, metas[k])
		assert.Equal(t, true, eg.DecideWithStruct(&metas[k]))
	}

}
