package expression

import (
	"fmt"
	"reflect"
)

func in_slice(dst []interface{}, org interface{}) bool {
	for _, v := range dst {
		fmt.Println(v, org, reflect.TypeOf(v), reflect.TypeOf(org))
		if v == org {
			return true
		}
	}
	return false
}

func decide_eq(dst interface{}, org interface{}) bool {
	return dst == org
}

func decide_ne(dst interface{}, org interface{}) bool {
	return dst != org
}

func decide_nin(dst interface{}, org interface{}) bool {
	return false
}

func decide_in(dst interface{}, org interface{}) bool {

	v := reflect.ValueOf(dst)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return false
		}

		return in_slice(dst.([]interface{}), org)

	}

	return false
}

func decide_lt(dst interface{}, org interface{}) bool {

	var a, b float64
	switch reflect.TypeOf(org).Kind() {
	case reflect.Int:
		b = float64(org.(int))
	case reflect.Int32:
		b = float64(org.(int32))
	case reflect.Float32:
		b = float64(org.(float32))
	case reflect.Float64:
		b = org.(float64)
	}

	switch reflect.TypeOf(dst).Kind() {
	case reflect.Int:
		a = float64(dst.(int))
	case reflect.Int32:
		a = float64(dst.(int32))
	case reflect.Float32:
		a = float64(dst.(float32))
	case reflect.Float64:
		a = dst.(float64)
	}

	return a < b
}

func decide_lte(dst interface{}, org interface{}) bool {

	var a, b float64
	switch reflect.TypeOf(org).Kind() {
	case reflect.Int:
		b = float64(org.(int))
	case reflect.Int32:
		b = float64(org.(int32))
	case reflect.Float32:
		b = float64(org.(float32))
	case reflect.Float64:
		b = org.(float64)
	}

	switch reflect.TypeOf(dst).Kind() {
	case reflect.Int:
		a = float64(dst.(int))
	case reflect.Int32:
		a = float64(dst.(int32))
	case reflect.Float32:
		a = float64(dst.(float32))
	case reflect.Float64:
		a = dst.(float64)
	}

	return a <= b
}

func decide_gt(dst interface{}, org interface{}) bool {

	var a, b float64
	switch reflect.TypeOf(org).Kind() {
	case reflect.Int:
		b = float64(org.(int))
	case reflect.Int32:
		b = float64(org.(int32))
	case reflect.Float32:
		b = float64(org.(float32))
	case reflect.Float64:
		b = org.(float64)
	}

	switch reflect.TypeOf(dst).Kind() {
	case reflect.Int:
		a = float64(dst.(int))
	case reflect.Int32:
		a = float64(dst.(int32))
	case reflect.Float32:
		a = float64(dst.(float32))
	case reflect.Float64:
		a = dst.(float64)
	}

	return a > b
}

func decide_gte(dst interface{}, org interface{}) bool {

	var a, b float64
	switch reflect.TypeOf(org).Kind() {
	case reflect.Int:
		b = float64(org.(int))
	case reflect.Int32:
		b = float64(org.(int32))
	case reflect.Float32:
		b = float64(org.(float32))
	case reflect.Float64:
		b = org.(float64)
	}

	switch reflect.TypeOf(dst).Kind() {
	case reflect.Int:
		a = float64(dst.(int))
	case reflect.Int32:
		a = float64(dst.(int32))
	case reflect.Float32:
		a = float64(dst.(float32))
	case reflect.Float64:
		a = dst.(float64)
	}

	return a >= b
}
