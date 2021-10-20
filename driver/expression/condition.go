package expression

import (
	"errors"
)

var (
	ErrParse        = errors.New("expression parse err")
	ErrEmptyValue   = errors.New("object empty value")
	ErrNotMatchRule = errors.New("parse value can't find match rule")
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
