package expression

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

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
	Heros   []string
	Ranks   []int32
	Market  *MetaMarket // 这里需要使用指针值，不然没办法使用反射向下传递
}

func TestParse(t *testing.T) {

	// 匹配成功的表达式列表
	succlst := []string{
		`$eq : { Token : '' }`,
		`$eq : {Market.ID : 'aabb'}`,
		`$gte :{Gold:100}`,
		`$lt : {Gold:100}`,
		`$lte : {Gold : 100}`,
		`$and : [$ne:{Token:''}, $gt:{Gold:101}}] `,
		`$or : [$gt:{Gold: 100}, $gt:{Diamond:100}]`,
		`$and : [ $eq:{Token:'aabb'}, $or:[$gt:{Diamond:100}, $gt:{Ticket:100}] ]`,
		`$eq : { Token : '中文' }`,
		`$in : { Heros : 'b' }`,
		`$in : {Ranks : 10}`,
		`$nin : {Heros : 'aa'}`,
	}

	// 匹配失败的表达式列表
	faillst := []string{}

	var metalst []map[string]interface{}
	// 被匹配的数据源
	orglst := []string{
		`{"Token":""}`,
		`{"Token":"", "Market":{ "ID": "aabb"}}`,
		`{"Gold":100}`,
		`{"Gold":99}`,
		`{"Gold":100}`,
		`{"Token":"aabb", "Gold":200}`,
		`{"Token":"aabb", "Gold":200}`,
		`{"Token":"aabb", "Gold":200, "Diamond":150}`,
		`{"Token":"中文"}`,
		`{"Heros":["a", "b", "c"]}`,
		`{"Ranks":[1,3,5,7,10]}`,
		`{"Heros":["a", "b", "c"]}`,
	}

	//assert.Equal(t, len(succlst), len(orglst))
	//assert.Equal(t, len(faillst), len(orglst))

	for _, v := range orglst {
		m := make(map[string]interface{})
		json.Unmarshal([]byte(v), &m)
		metalst = append(metalst, m)
	}

	meta := make(map[string]interface{})

	for k, v := range succlst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)

		mergo.MergeWithOverwrite(&meta, metalst[k])
		fmt.Println("equal", v)
		assert.Equal(t, true, eg.DecideWithMap(meta))
	}

	for k, v := range faillst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)

		mergo.MergeWithOverwrite(&meta, metalst[k])
		assert.Equal(t, false, eg.DecideWithMap(meta))
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
			Gold:    90,
			Diamond: 150,
		},
		{
			Gold:    100,
			Token:   "aabb",
			Diamond: 150,
		},
		{
			Token: "aabb",
			Gold:  110,
		},
		{
			Gold: 200,
		},
		{
			Token:   "aabb",
			Gold:    200,
			Diamond: 150,
		},
		{
			Token: "中文",
		},
		{
			Heros: []string{"a", "b", "c", "d"},
		},
		{
			Ranks: []int32{1, 3, 5, 7, 9, 10},
		},
		{
			Heros: []string{"a", "b", "c", "d"},
		},
	}

	for k, v := range succlst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)

		fmt.Println(k, metas[k])
		assert.Equal(t, true, eg.DecideWithStruct(&metas[k]))
	}

}

func BenchmarkParseMap(b *testing.B) {

	// 匹配成功的表达式列表
	succlst := []string{
		`$eq : { Token : '' }`,
		`$eq : {Market.ID : 'aabb'}`,
		`$gte :{Gold:100}`,
		`$lt : {Gold:100}`,
		`$lte : {Gold : 100}`,
		`$and : [$ne:{Token:''}, $gt:{Gold:101}}] `,
		`$or : [$gt:{Gold: 100}, $gt:{Diamond:100}]`,
		`$and : [ $eq:{Token:'aabb'}, $or:[$gt:{Diamond:100}, $gt:{Ticket:100}] ]`,
		`$eq : { Token : '中文' }`,
		`$in : { Heros : 'b' }`,
		`$in : {Ranks : 10}`,
	}

	var metalst []map[string]interface{}
	// 被匹配的数据源
	orglst := []string{
		`{"Token":""}`,
		`{"Token":"", "Market":{ "ID": "aabb"}}`,
		`{"Gold":100}`,
		`{"Gold":99}`,
		`{"Gold":100}`,
		`{"Token":"aabb", "Gold":200}`,
		`{"Token":"aabb", "Gold":200}`,
		`{"Token":"aabb", "Gold":200, "Diamond":150}`,
		`{"Token":"中文"}`,
		`{"Heros":["a", "b", "c"]}`,
		`{"Ranks":[1,3,5,7,10]}`,
	}

	for _, v := range orglst {
		m := make(map[string]interface{})
		json.Unmarshal([]byte(v), &m)
		metalst = append(metalst, m)
	}

	meta := make(map[string]interface{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		for k, v := range succlst {
			eg, _ := Parse(v)
			mergo.MergeWithOverwrite(&meta, metalst[k])
			eg.DecideWithMap(meta)
		}

	}

}

func BenchmarkParseStruct(b *testing.B) {
	// 匹配成功的表达式列表
	succlst := []string{
		`$eq : { Token : '' }`,
		`$eq : {Market.ID : 'aabb'}`,
		`$gte :{Gold:100}`,
		`$lt : {Gold:100}`,
		`$lte : {Gold : 100}`,
		`$and : [$ne:{Token:''}, $gt:{Gold:101}}] `,
		`$or : [$gt:{Gold: 100}, $gt:{Diamond:100}]`,
		`$and : [ $eq:{Token:'aabb'}, $or:[$gt:{Diamond:100}, $gt:{Ticket:100}] ]`,
		`$eq : { Token : '中文' }`,
		`$in : { Heros : 'b' }`,
		`$in : {Ranks : 10}`,
	}

	rand.Seed(time.Now().UnixNano())

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
			Gold:    90,
			Diamond: 150,
		},
		{
			Gold:    100,
			Token:   "aabb",
			Diamond: 150,
		},
		{
			Token: "aabb",
			Gold:  110,
		},
		{
			Gold: 200,
		},
		{
			Token:   "aabb",
			Gold:    200,
			Diamond: 150,
		},
		{
			Token: "中文",
		},
		{
			Heros: []string{"a", "b", "c", "d"},
		},
		{
			Ranks: []int32{1, 3, 5, 7, 9, 10},
		},
	}

	var eg []*ExpressionGroup

	for _, v := range succlst {
		e, _ := Parse(v)
		eg = append(eg, e)
	}

	eglen := len(eg)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ei := rand.Intn(eglen)
		eg[ei].DecideWithStruct(&metas[ei])
	}

}
