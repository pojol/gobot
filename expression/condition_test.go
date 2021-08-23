package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	列表 [  ]
	键值 {  }
	表达式 $symbol : {}
*/

type Market struct {
	ID string
}

type Meta struct {
	Token   string
	Diamond int32
	Gold    int32
	Ticket  int32
	Market  *Market
}

func TestParse(t *testing.T) {

	lst := []string{
		`$eq : { Token : "" }`,
		`$eq : {Market.ID : "aabb"}`,
		`$and : [$ne:{Token:""}, $gt:{Gold:100}}] `,
		// `$or`
		//`$and : [ $eq:{meta.token:"aabb"}, $or:[$gt:{meta.diamond:100}, $gt:{meta.ticket:100}] ]`,
	}

	metas := []Meta{
		{ // $eq : { meta.Token : "" }
			Token: "",
		},
		{
			Market: &Market{ID: "aabb"},
		},
		{ //$and : [$ne:{meta.Token:""}, $gt:{meta.Gold:100}}]
			Token: "aabb",
			Gold:  200,
		},
		{ //$and : [ $eq:{meta.Token:"aabb"}, $or:[$gt:{meta.Diamond:100}, $gt:{meta.Ticket:100}] ]

		},
	}

	for k, v := range lst {
		eg, err := Parse(v)
		assert.Equal(t, err, nil)
		assert.Equal(t, true, eg.Decide(&metas[k]))
	}
}
