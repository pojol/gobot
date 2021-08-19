// Package plugin 用于加载插件
package plugins

import (
	"errors"
	"plugin"

	"github.com/pojol/apibot/marshal"
)

type Plugin struct {
	Name string
	Type string

	CreateFunc func() marshal.IMarshal
}

var plugins map[string]marshal.IMarshal

func Get(name string) marshal.IMarshal {
	if _, ok := plugins[name]; ok {
		return plugins[name]
	}

	return nil
}

// Load loads a plugin created with `go build -buildmode=plugin`
func Load(path string) error {

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
	plugins = make(map[string]marshal.IMarshal)
}
