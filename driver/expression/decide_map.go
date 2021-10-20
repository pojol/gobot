package expression

import (
	"strings"
)

func getValueWithMap(key string, m map[string]interface{}) interface{} {
	var parent, child string
	var inner bool
	idx := strings.Index(key, ".")
	if idx != -1 {
		parent = key[:idx]
		child = key[idx+1:]
		inner = true
	}

	if inner {
		for k := range m {
			if k == parent {
				return getValueWithMap(child, m[k].(map[string]interface{}))
			}
		}
	} else {
		if _, ok := m[key]; ok {
			return m[key]
		}
	}

	return nil
}

func (e *Expression) decidemap(m map[string]interface{}) bool {

	var b bool

	switch e.Symbol {
	case AND:
		for _, v := range e.Exprs {
			b = v.decidemap(m)
			if !b {
				break
			}
		}
	case OR:
		for _, v := range e.Exprs {
			b = v.decidemap(m)
			if b {
				break
			}
		}
	case EQ, NE, LT, LTE, GT, GTE:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_conv(v, e.Object.Right, e.Symbol)
	case NIN:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_nin(v, e.Object.Right)
	case IN:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_in(v, e.Object.Right)
	default:
		println("decide unknown symbol", e.Symbol)
	}

	return b
}

func (eg *ExpressionGroup) DecideWithMap(m map[string]interface{}) bool {
	return eg.Root.decidemap(m)
}
