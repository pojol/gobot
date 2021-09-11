package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// from https://github.com/yuin/gluamapper

// Option is a configuration that is used to create a new mapper.
type Option struct {
	// Function to convert a lua table key to Go's one. This defaults to "ToUpperCamelCase".
	NameFunc func(string) string

	// Returns error if unused keys exist.
	ErrorUnused bool

	// A struct tag name for lua table keys . This defaults to "gluamapper"
	TagName string
}

var camelre = regexp.MustCompile(`_([a-z])`)

// ToUpperCamelCase is an Option.NameFunc that converts strings from snake case to upper camel case.
func ToUpperCamelCase(s string) string {
	return strings.ToUpper(string(s[0])) + camelre.ReplaceAllStringFunc(s[1:len(s)], func(s string) string { return strings.ToUpper(s[1:len(s)]) })
}

// Map maps the lua table to the given struct pointer.
func Table2Map(tbl *lua.LTable) (map[interface{}]interface{}, error) {
	opt := Option{
		NameFunc: ToUpperCamelCase,
		TagName:  "gluamapper",
	}
	mp, ok := ToGoValue(tbl, opt).(map[interface{}]interface{})
	if !ok {
		return mp, errors.New("arguments #1 must be a table, but got an array")
	}

	return mp, nil
}

// ToGoValue converts the given LValue to a Go object.
func ToGoValue(lv lua.LValue, opt Option) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[interface{}]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(ToGoValue(key, opt))
				ret[opt.NameFunc(keystr)] = ToGoValue(value, opt)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i), opt))
			}
			return ret
		}
	default:
		return v
	}
}
