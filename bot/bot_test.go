package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pojol/apibot/plugins"
	"github.com/stretchr/testify/assert"
)

var srv *httptest.Server

type guestRes struct {
	Token string
}

type infoRes struct {
	Diamond int32
	Gold    int32
}

func TestMain(m *testing.M) {

	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ioutil.ReadAll(req.Body)
		var byt []byte

		fmt.Println("http server recv ", req.RequestURI)

		if req.RequestURI == "/login/guest" {
			byt, _ = json.Marshal(guestRes{
				Token: "abcd",
			})
		} else if req.RequestURI == "/base/acc.info" {
			byt, _ = json.Marshal(infoRes{
				Diamond: 100,
				Gold:    100,
			})
		}

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
    "id":"b36fabfd-dd9a-4d24-941c-69f64233a589",
    "ty":"RootNode",
    "pos":{
        "x":0,
        "y":0
    },
    "children":[
        {
            "id":"7872b200-52ef-40f8-8059-d019fca99501",
            "ty":"LoopNode",
            "pos":{
                "x":-5,
                "y":47
            },
            "children":[
                {
                    "id":"0059758d-cba6-4718-a98f-19bcad43f975",
                    "ty":"SelectorNode",
                    "pos":{
                        "x":-15,
                        "y":126
                    },
                    "children":[
                        {
                            "id":"31121c03-787d-44f0-88ec-81a440700c61",
                            "ty":"ConditionNode",
                            "pos":{
                                "x":-20,
                                "y":179
                            },
                            "children":[
                                {
                                    "id":"8fe159f5-fcb5-4106-9b3c-ed7f950cd547",
                                    "ty":"HTTPActionNode",
                                    "pos":{
                                        "x":-50,
                                        "y":245
                                    },
                                    "children":[

                                    ],
                                    "api":"/login/guest",
                                    "parm":{

                                    }
                                }
                            ],
                            "expr":"$eq:{Token:''}"
                        },
                        {
                            "id":"5fe427f2-a19e-40c9-a17a-10dd994134f5",
                            "ty":"ConditionNode",
                            "pos":{
                                "x":50,
                                "y":179
                            },
                            "children":[
                                {
                                    "id":"1b2569a6-bb1b-4b88-8866-7d93cf552739",
                                    "ty":"HTTPActionNode",
                                    "pos":{
                                        "x":55,
                                        "y":245
                                    },
                                    "children":[

                                    ],
                                    "api":"/base/acc.info",
                                    "parm":{
                                        "token":"meta.token"
                                    }
                                }
                            ],
                            "expr":"$ne:{Token:''}"
                        }
                    ]
                }
            ],
            "loop":3
        }
    ]
}`

func TestBot(t *testing.T) {

	err := plugins.Load("../plugins/json/json.so")
	assert.Equal(t, err, nil)

	b, _ := NewWithBehaviorFile([]byte(compose), srv.URL)
	b.Run()
}

func TestStep(t *testing.T) {
	err := plugins.Load("../plugins/json/json.so")
	assert.Equal(t, err, nil)

	b, _ := NewWithBehaviorFile([]byte(compose), srv.URL)
	for i := 0; i < 30; i++ {
		b.RunStep()
	}
}
