package behavior

import "testing"

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

func BenchmarkNew(b *testing.B) {

	byt := []byte(compose)

	for i := 0; i < b.N; i++ {
		New(byt)
	}

}
