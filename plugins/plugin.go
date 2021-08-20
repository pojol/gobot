// Package plugin 用于加载插件
package plugins

import (
	"errors"
	"fmt"
	"plugin"

	"github.com/pojol/apibot/serialization"
)

type Plugin struct {
	Name string
	Type string

	CreateFunc func() serialization.ISerialization
}

var plugins map[string]serialization.ISerialization

func Get(name string) serialization.ISerialization {
	if _, ok := plugins[name]; ok {
		return plugins[name]
	}

	return nil
}

// Load loads a plugin created with `go build -buildmode=plugin`
func Load(path string) error {

	fmt.Println("load", path)

	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	s, err := p.Lookup("Plugin")
	if err != nil {
		return err
	}
	pl, ok := s.(*Plugin)
	if !ok {
		return errors.New("could not cast Plugin object")
	}
	switch pl.Type {
	case "json":
		plugins[pl.Name] = pl.CreateFunc()
	}

	return nil
}

func init() {
	plugins = make(map[string]serialization.ISerialization)
}
