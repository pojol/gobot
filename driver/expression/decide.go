package expression

import (
	"reflect"
)

func decide_in(dst interface{}, exprval interface{}) bool {

	v := reflect.ValueOf(dst)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return false
		}

		for i := 0; i < v.Len(); i++ {
			if decide_conv(v.Index(i).Interface(), exprval, EQ) {
				return true
			}
		}

	}

	return false
}

func decide_nin(dst interface{}, exprval interface{}) bool {
	v := reflect.ValueOf(dst)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return true
		}

		for i := 0; i < v.Len(); i++ {
			if decide_conv(v.Index(i).Interface(), exprval, EQ) {
				return false
			}
		}
		return true
	}

	return false
}

func decide_eq(dst interface{}, exprval interface{}) bool {
	return dst == exprval
}

func decide_ne(dst interface{}, exprval interface{}) bool {
	return dst != exprval
}

func decide_lt_int64(dst, exprval int64) bool {
	return dst < exprval
}

func decide_lt_float64(dst, exprval float64) bool {
	return dst < exprval
}

func decide_lte_int64(dst, exprval int64) bool {
	return dst <= exprval
}

func decide_lte_float64(dst, exprval float64) bool {
	return dst <= exprval
}

func decide_gt_int64(dst, exprval int64) bool {
	return dst > exprval
}

func decide_gt_float64(dst, exprval float64) bool {
	return dst > exprval
}

func decide_gte_int64(dst, exprval int64) bool {
	return dst >= exprval
}

func decide_gte_float64(dst, exprval float64) bool {
	return dst >= exprval
}

func decide_conv_symbol_float64(dst, exprval float64, symbol string) bool {

	switch symbol {
	case EQ:
		return decide_eq(dst, exprval)
	case NE:
		return decide_ne(dst, exprval)
	case GT:
		return decide_gt_float64(dst, exprval)
	case GTE:
		return decide_gte_float64(dst, exprval)
	case LT:
		return decide_lt_float64(dst, exprval)
	case LTE:
		return decide_lte_float64(dst, exprval)
	}

	return false
}

func decide_conv_symbol_int64(dst, exprval int64, symbol string) bool {

	switch symbol {
	case EQ:
		return decide_eq(dst, exprval)
	case NE:
		return decide_ne(dst, exprval)
	case GT:
		return decide_gt_int64(dst, exprval)
	case GTE:
		return decide_gte_int64(dst, exprval)
	case LT:
		return decide_lt_int64(dst, exprval)
	case LTE:
		return decide_lte_int64(dst, exprval)
	}

	return false
}

func decide_conv_symbol_other(dst, exprval interface{}, symbol string) bool {

	switch symbol {
	case EQ:
		return decide_eq(dst, exprval)
	case NE:
		return decide_ne(dst, exprval)
	}

	return false
}

func decide_conv(dst interface{}, exprval interface{}, symbol string) bool {
	var ie int64
	var fe float64

	et := reflect.TypeOf(exprval)

	switch reflect.TypeOf(dst).Kind() {
	case reflect.Float64:
		fd, _ := dst.(float64)
		if et.Kind() == reflect.Float64 {
			fe, _ = exprval.(float64)
			return decide_conv_symbol_float64(fd, fe, symbol)
		} else if et.Kind() == reflect.Int64 {
			ie, _ = exprval.(int64)
			return decide_conv_symbol_float64(fd, float64(ie), symbol)
		}
	case reflect.Float32:
		fd, _ := dst.(float32)
		if et.Kind() == reflect.Float64 {
			fe, _ = exprval.(float64)
			return decide_conv_symbol_float64(float64(fd), fe, symbol)
		} else if et.Kind() == reflect.Int64 {
			ie, _ = exprval.(int64)
			return decide_conv_symbol_float64(float64(fd), float64(ie), symbol)
		}
		return float64(fd) < fe
	case reflect.Int16:
		id, _ := dst.(int16)
		if et.Kind() == reflect.Float64 {
			fe, _ = exprval.(float64)
			return decide_conv_symbol_int64(int64(id), int64(fe), symbol)
		} else if et.Kind() == reflect.Int64 {
			ie, _ = exprval.(int64)
			return decide_conv_symbol_int64(int64(id), int64(ie), symbol)
		}
	case reflect.Int32:
		id, _ := dst.(int32)
		if et.Kind() == reflect.Float64 {
			fe, _ = exprval.(float64)
			return decide_conv_symbol_int64(int64(id), int64(fe), symbol)
		} else if et.Kind() == reflect.Int64 {
			ie, _ = exprval.(int64)
			return decide_conv_symbol_int64(int64(id), int64(ie), symbol)
		}
	case reflect.Int:
		id, _ := dst.(int)
		if et.Kind() == reflect.Float64 {
			fe, _ = exprval.(float64)
			return decide_conv_symbol_int64(int64(id), int64(fe), symbol)
		} else if et.Kind() == reflect.Int64 {
			ie, _ = exprval.(int64)
			return decide_conv_symbol_int64(int64(id), int64(ie), symbol)
		}
	case reflect.Int64:
		id, _ := dst.(int64)
		if et.Kind() == reflect.Float64 {
			fe, _ = exprval.(float64)
			return decide_conv_symbol_int64(id, int64(fe), symbol)
		} else if et.Kind() == reflect.Int64 {
			ie, _ = exprval.(int64)
			return decide_conv_symbol_int64(id, int64(ie), symbol)
		}
	case reflect.String, reflect.Bool:
		return decide_conv_symbol_other(dst, exprval, symbol)
	}

	return false
}
