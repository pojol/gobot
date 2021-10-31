package main

import (
	"encoding/json"

	"github.com/pojol/gobot/driver/plugins"
	"github.com/pojol/gobot/driver/serialization"
)

type jsonp struct {
}

func (jm *jsonp) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jm *jsonp) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

var Plugin = plugins.Plugin{
	Name: "jsonparse",
	Type: "json",
	CreateFunc: func() serialization.ISerialization {

		jm := &jsonp{}

		return jm
	},
}
