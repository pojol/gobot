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
	case AND:
		for _, v := range e.Exprs {
			b = v.decide(meta)
			if !b {
				break
			}
		}
	case OR:
		for _, v := range e.Exprs {
			b = v.decide(meta)
			if b {
				break
			}
		}
	case EQ, NE, LT, LTE, GT, GTE:
		v := getValueWithStrcut(e.Object.Left, meta)
		b = decide_conv(v, e.Object.Right, e.Symbol)

	case NIN:
		v := getValueWithStrcut(e.Object.Left, meta)
		b = decide_nin(v, e.Object.Right)

	case IN:
		v := getValueWithStrcut(e.Object.Left, meta)
		b = decide_in(v, e.Object.Right)

	default:
		println("decide", e.Symbol)
	}

	return b
}

func (eg *ExpressionGroup) DecideWithStruct(meta interface{}) bool {
	return eg.Root.decide(meta)
}
