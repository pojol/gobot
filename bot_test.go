package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pojol/apibot/assert"
	"github.com/pojol/apibot/plugins"
)

var srv *httptest.Server

type res struct {
	Token string
}

func TestMain(m *testing.M) {

	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		body, _ := ioutil.ReadAll(req.Body)
		fmt.Println("http server recv ", body)

		byt, _ := json.Marshal(res{
			Token: "abcd",
		})

		w.Write(byt)
	}))
	defer srv.Close()

	os.Exit(m.Run())
}

type Metadata struct {
	Val string
}

var compose = `
{
    "id":"e41d33fa-4658-49ef-ab00-db2919de2dbc",
    "ty":"RootNode",
    "pos":{
        "x":0,
        "y":0
    },
    "children":[
        {
            "id":"65c0cb5c-8cf3-427d-a6a2-5ee636627fef",
            "ty":"SelectorNode",
            "pos":{
                "x":-15,
                "y":67
            },
            "children":[
                {
                    "id":"f758e386-c521-4e50-8be7-7300975b354c",
                    "ty":"ConditionNode",
                    "pos":{
                        "x":-20,
                        "y":130
                    },
                    "children":[
                        {
                            "id":"a7609b7b-064d-4a0e-b166-d89f20fcabcd",
                            "ty":"HTTPActionNode",
                            "pos":{
                                "x":-35,
                                "y":192
                            },
                            "children":[

                            ],
                            "api":"/login/guest",
                            "parm":{

                            }
                        }
                    ],
                    "script":{
                        "$eq":{
                            "meta.Token":""
                        }
                    }
                },
                {
                    "id":"24f27096-453f-4035-819b-d5dbfca1471b",
                    "ty":"ConditionNode",
                    "pos":{
                        "x":50,
                        "y":130
                    },
                    "children":[
                        {
                            "id":"ea54ba5d-366c-4248-a12e-2fda71122514",
                            "ty":"HTTPActionNode",
                            "pos":{
                                "x":35,
                                "y":192
                            },
                            "children":[

                            ],
                            "api":"/base/acc.info",
                            "parm":{
                                "Token":"meta.Token"
                            }
                        }
                    ],
                    "script":{
                        "$ne":{
                            "meta.Token":""
                        }
                    }
                }
            ]
        }
    ]
}
`

func TestBot(t *testing.T) {

	err := plugins.Load("plugins/json/json.so")
	assert.Equal(t, err, nil)

	md := Metadata{}
	b, _ := NewWithBehaviorFile([]byte(compose), srv.URL, &md)

	b.Run()
}
