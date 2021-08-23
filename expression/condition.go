package expression

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrParse = errors.New("expression parse err")
)

type ICondition interface {
	Validator(dst interface{}, org interface{}) bool
}

type ExprVal struct {
	Left  string
	Right interface{}
}

type Expression struct {
	Symbol string
	Exprs  []Expression
	Object ExprVal
}

type ExpressionGroup struct {
	Root Expression
}

func decide_eq(dst interface{}, org interface{}) bool {
	return dst == org
}

func decide_ne(dst interface{}, org interface{}) bool {
	return dst != org
}

func decide_gt(dst interface{}, org interface{}) bool {
	a, _ := dst.(int32)
	b, _ := org.(int32)
	return a > b
}

func getValue(key string, org interface{}) interface{} {

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
					return getValue(child, v.Field(i).Interface())
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
	case "$and", "$or":
		for _, v := range e.Exprs {
			b = v.decide(meta)
			if !b {
				break
			}
		}
	case "$eq":
		v := getValue(e.Object.Left, meta)
		b = decide_eq(v, e.Object.Right)

	case "$ne":
		v := getValue(e.Object.Left, meta)
		b = decide_ne(v, e.Object.Right)

	case "$gt":
		v := getValue(e.Object.Left, meta)
		b = decide_gt(v, e.Object.Right)
	default:
		println("decide", e.Symbol)
	}

	return b
}

func (eg *ExpressionGroup) Decide(meta interface{}) bool {

	return eg.Root.decide(meta)

}

func (eg *ExpressionGroup) parse_symbol(symbol string, val string, e *Expression) {

	switch symbol {
	case "$and", "$or":
		eg.ParseAndOr(val, e)
	case "$eq", "$gt", "$ne":
		eg.ParseObject(symbol, val, e)
	}

}

func (eg *ExpressionGroup) ParseAndOr(val string, e *Expression) {

	content := val[1 : len(val)-1]

	for i := 0; i < len(content); {
		subcontent := content[i:]
		//fmt.Println("sub", subcontent)
		et := Expression{}

		idx := strings.Index(subcontent, ":")
		symbol := subcontent[:idx]
		//fmt.Println("sub symbol", left)

		switch symbol {
		case "$and", "$or":
			et.Symbol = symbol
			tailIdx := strings.Index(subcontent, "]")
			right := subcontent[idx+1 : tailIdx+1]

			i += tailIdx + 2
			eg.parse_symbol(symbol, right, &et)
		case "$eq", "$gt", "$ne":

			tailIdx := strings.Index(subcontent, "}")
			right := subcontent[idx+1 : tailIdx+1]
			//fmt.Println("sub right", right)

			i += tailIdx + 2

			//fmt.Println(i, tailIdx+2, len(content))
			eg.parse_symbol(symbol, right, &et)
		default:
			fmt.Println("??")
			goto ext
		}

		e.Exprs = append(e.Exprs, et)
	}
ext:
}

func (eg *ExpressionGroup) ParseObject(symbol, val string, e *Expression) {
	e.Symbol = symbol
	obj := val[1 : len(val)-1]

	objArr := strings.Split(obj, ":")
	e.Object.Left = objArr[0]

	// number, string, true, false

	nv := strings.Replace(objArr[1], "\"", "", -1)
	e.Object.Right = nv
}

func printTree(e *Expression) {

	switch e.Symbol {
	case "$and", "$or":
		fmt.Println(e.Symbol)
		for _, v := range e.Exprs {
			printTree(&v)
		}
	case "$eq", "$gt", "$ne":
		fmt.Println(e.Symbol, "left", e.Object.Left, "right", e.Object.Right)
	default:
		println("!!", e.Symbol)
	}

}

// 后面还要加上语法的校验
func Parse(s string) (*ExpressionGroup, error) {

	eg := ExpressionGroup{
		Root: Expression{},
	}

	s = strings.Replace(s, " ", "", -1)
	if strings.Count(s, ":") < 2 {
		return nil, ErrParse
	}

	idx := strings.Index(s, ":")
	left, right := s[:idx], s[idx+1:]
	eg.Root.Symbol = left

	eg.parse_symbol(left, right, &eg.Root)
	//printTree(&eg.Root)

	return &eg, nil
}
