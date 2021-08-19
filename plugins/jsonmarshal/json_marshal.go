package main

import (
	"fmt"

	"github.com/pojol/apibot/marshal"
	"github.com/pojol/apibot/plugins"
)

type jmarshal struct {
}

func (jm *jmarshal) Marshal(v interface{}) ([]byte, error) {

	fmt.Println("json marshal", v)

	return []byte(""), nil
}

var Plugin = plugins.Plugin{
	Name: "jsonmarshal",
	Type: "json",
	CreateFunc: func() marshal.IMarshal {

		jm := &jmarshal{}

		return jm
	},
}
