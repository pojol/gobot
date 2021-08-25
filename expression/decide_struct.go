package expression

import (
	"reflect"
	"strings"
)

func getValueWithStrcut(key string, org interface{}) interface{} {

	t := reflect.TypeOf(org).Elem()
	v := reflect.ValueOf(org).Elem()

	var parent, child string
	var inner bool
	idx := strings.Index(key, ".")
	if idx != -1 {
		parent = key[:idx]
		child = key[idx+1:]
		inner = true
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() { // 检测是否可导出字段

			fieldName := t.Field(i).Name

			if inner {
				if fieldName == parent {
					return getValueWithStrcut(child, v.Field(i).Interface())
				}

			} else {

				if fieldName == key {
					return v.Field(i).Interface()
				}
			}

		}
	}

	return nil
}

func (e *Expression) decide(meta interface{}) bool {

	var b bool

	switch e.Symbol {
	case "$and":
		for _, v := range e.Exprs {
			b = v.decide(meta)
			if !b {
				break
			}
		}
	case "$or":
		for _, v := range e.Exprs {
			b = v.decide(meta)
			if b {
				break
			}
		}
	case "$eq":
		v := getValueWithStrcut(e.Object.Left, meta)
		b = decide_eq(v, e.Object.Right)

	case "$ne":
		v := getValueWithStrcut(e.Object.Left, meta)
		b = decide_ne(v, e.Object.Right)

	case "$gt":
		v := getValueWithStrcut(e.Object.Left, meta)
		b = decide_gt(v, e.Object.Right)
	default:
		println("decide", e.Symbol)
	}

	return b
}

func (eg *ExpressionGroup) DecideWithStruct(meta interface{}) bool {
	return eg.Root.decide(meta)
}
