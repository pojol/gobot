package expression

import (
	"fmt"
	"strconv"
	"strings"
)

// isValidNumber reports whether s is a valid JSON number literal.
func isValidNumber(s string) bool {
	// This function implements the JSON numbers grammar.
	// See https://tools.ietf.org/html/rfc7159#section-6
	// and https://www.json.org/img/number.png

	if s == "" {
		return false
	}

	// Optional -
	if s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case s[0] == '0':
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			if s == "" {
				return false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// Make sure we are at the end.
	return s == ""
}

func (eg *ExpressionGroup) parse_symbol(symbol string, val string, e *Expression) {

	switch symbol {
	case OR, AND:
		eg.ParseAndOr(val, e)
	case GTE, GT, LTE, LT, EQ, NE, IN, NIN:
		eg.ParseObject(symbol, val, e)
	}

}

func (eg *ExpressionGroup) ParseAndOr(val string, e *Expression) {

	content := val[1 : len(val)-1]

	for i := 0; i < len(content); {
		subcontent := content[i:]
		et := Expression{}

		idx := strings.Index(subcontent, ":")
		symbol := subcontent[:idx]

		switch symbol {
		case AND, OR:
			et.Symbol = symbol
			tailIdx := strings.Index(subcontent, "]")
			right := subcontent[idx+1 : tailIdx+1]

			i += tailIdx + 2
			eg.parse_symbol(symbol, right, &et)
		case GTE, GT, LTE, LT, EQ, NE, IN, NIN:

			tailIdx := strings.Index(subcontent, "}")
			right := subcontent[idx+1 : tailIdx+1]

			i += tailIdx + 2

			eg.parse_symbol(symbol, right, &et)
		default:
			fmt.Println("parse unknown symbol", symbol)
			goto ext
		}

		e.Exprs = append(e.Exprs, et)
	}
ext:
}

// 依据对字符串的解析，动态创建值并赋予 right
func (eg *ExpressionGroup) ParseObject(symbol, val string, e *Expression) {
	e.Symbol = symbol
	obj := val[1 : len(val)-1] // 去掉 {}

	objArr := strings.Split(obj, ":")
	e.Object.Left = objArr[0]

	vlen := len(objArr[1])
	// int, float32, string, true, false
	if vlen == 0 {
		panic(ErrEmptyValue)
	}

	if vlen >= 2 {
		if objArr[1][0] == '\'' && objArr[1][vlen-1] == '\'' {
			if vlen == 2 {
				e.Object.Right = ""
			} else {
				e.Object.Right = objArr[1][1 : vlen-1]
			}

			return
		}
	}

	if objArr[1] == "true" {
		e.Object.Right = true
		return
	}

	if objArr[1] == "false" {
		e.Object.Right = false
		return
	}

	if isValidNumber(objArr[1]) {

		if strings.Contains(objArr[1], ".") {
			fv, err := strconv.ParseFloat(objArr[1], 64)
			if err != nil {
				panic(err)
			}
			e.Object.Right = fv
		} else {
			iv, err := strconv.ParseInt(objArr[1], 10, 64)
			if err != nil {
				panic(err)
			}
			e.Object.Right = int32(iv)
		}

		return
	}

	panic(ErrNotMatchRule)

}

// 后面还要加上语法的校验, 在 parse阶段遇到问题 直接 panic
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

	return &eg, nil
}
