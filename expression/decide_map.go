package expression

import (
	"fmt"
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
	case EQ:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_eq(v, e.Object.Right)
	case NE:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_ne(v, e.Object.Right)
	case LT:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_lt(v, e.Object.Right)
	case LTE:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_lte(v, e.Object.Right)
	case GT:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_gt(v, e.Object.Right)
	case GTE:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_gte(v, e.Object.Right)
	case NIN:
		v := getValueWithMap(e.Object.Left, m)
		b = decide_nin(v, e.Object.Right)
	case IN:
		v := getValueWithMap(e.Object.Left, m)
		fmt.Println("map in", v)
		b = decide_in(v, e.Object.Right)
	default:
		println("decide unknown symbol", e.Symbol)
	}

	return b
}

func (eg *ExpressionGroup) DecideWithMap(m map[string]interface{}) bool {
	return eg.Root.decidemap(m)
}
