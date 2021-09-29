package behavior

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"

	pb "github.com/golang/protobuf/ptypes/struct"
)

// Encode encodes map[string]interface{} to pb.Struct
func Encode(src map[string]interface{}, dest *pb.Struct) error {
	if dest == nil {
		return errors.New("structpb: failed to encode. dest struct is nil")
	}
	if src == nil || len(src) < 1 {
		return nil
	}
	for k, v := range src {
		var pbv pb.Value
		if err := EncodeValue(v, &pbv); err != nil {
			return err
		}

		if dest.Fields == nil {
			dest.Fields = make(map[string]*pb.Value)
		}
		dest.Fields[k] = &pbv
	}
	return nil
}

// EncodeFromStruct encodes an interface of struct to pb.Struct
func EncodeFromStruct(srcStruct interface{}, dest *pb.Struct) error {
	if dest == nil {
		return errors.New("structpb: failed to encode. dest struct is nil")
	}
	rv := reflect.ValueOf(srcStruct)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = reflect.Indirect(rv)
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("structpb: failed to encode. unexpected kind of src struct %s", rv.Kind())
	}
	rt := rv.Type()
	numFields := rt.NumField()
	if numFields < 1 {
		return nil
	}
	for i := 0; i < numFields; i++ {
		fname := rt.Field(i).Name
		if !exportedFieldName(fname) {
			continue
		}
		var pbv pb.Value
		if err := EncodeValue(rv.FieldByName(fname).Interface(), &pbv); err != nil {
			return err
		}
		if dest.Fields == nil {
			dest.Fields = make(map[string]*pb.Value, numFields)
		}
		dest.Fields[fname] = &pbv
	}
	return nil
}

// EncodeValue encodes src value to pb.Struct
func EncodeValue(src interface{}, dest *pb.Value) error {
	if dest == nil {
		return errors.New("structpb: failed to encode. dest value is nil")
	}
	switch s := src.(type) {
	case nil:
		dest.Kind = &pb.Value_NullValue{}
	case int:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *int:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case int8:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *int8:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case int16:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *int16:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case int32:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *int32:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case int64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *int64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case uint:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *uint:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case uint8:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *uint8:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case uint16:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *uint16:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case uint32:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *uint32:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case uint64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *uint64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case float32:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(s)}
	case *float32:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(*s)}
	case float64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: s}
	case *float64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: *s}
	case string:
		dest.Kind = &pb.Value_StringValue{StringValue: s}
	case *string:
		dest.Kind = &pb.Value_StringValue{StringValue: *s}
	case bool:
		dest.Kind = &pb.Value_BoolValue{BoolValue: s}
	case *bool:
		dest.Kind = &pb.Value_BoolValue{BoolValue: *s}
	default:
		if err := encodeValueReflect(src, dest); err != nil {
			return err
		}
	}
	if dest.GetKind() == nil {
		panic("pb.Value.Kind is nil")
	}
	return nil
}

func encodeValueReflect(src interface{}, dest *pb.Value) error {
	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			dest.Kind = &pb.Value_NullValue{}
			return nil
		}
		rv = reflect.Indirect(rv)
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(rv.Int())}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(rv.Uint())}
	case reflect.Float32, reflect.Float64:
		dest.Kind = &pb.Value_NumberValue{NumberValue: float64(rv.Float())}
	case reflect.String:
		dest.Kind = &pb.Value_StringValue{StringValue: rv.String()}
	case reflect.Bool:
		dest.Kind = &pb.Value_BoolValue{BoolValue: rv.Bool()}
	case reflect.Struct:
		var pbst pb.Struct
		if err := EncodeFromStruct(src, &pbst); err != nil {
			return err
		}
		dest.Kind = &pb.Value_StructValue{StructValue: &pbst}

	case reflect.Map:
		pbst := pb.Struct{}
		for iter := rv.MapRange(); iter.Next(); {
			k, v := iter.Key(), iter.Value()
			if k.Kind() != reflect.String {
				// support only string key
				continue
			}
			var pbv pb.Value
			if err := EncodeValue(v.Interface(), &pbv); err != nil {
				return err
			}
			if pbst.Fields == nil {
				pbst.Fields = make(map[string]*pb.Value, rv.Len())
			}
			pbst.Fields[k.String()] = &pbv
		}
		dest.Kind = &pb.Value_StructValue{StructValue: &pbst}

	case reflect.Array, reflect.Slice:
		pblv := pb.ListValue{}
		for i := 0; i < rv.Len(); i++ {
			var pbv pb.Value
			if err := EncodeValue(rv.Index(i).Interface(), &pbv); err != nil {
				return err
			}
			if pblv.Values == nil {
				pblv.Values = make([]*pb.Value, rv.Len())
			}
			pblv.Values[i] = &pbv
		}
		dest.Kind = &pb.Value_ListValue{ListValue: &pblv}

	default:
		return fmt.Errorf("structpb: failed to encode. unexpected src kind %s", rv.Kind())
	}
	return nil
}

func exportedFieldName(name string) bool {
	if len(name) < 1 {
		panic("structpb: name len is zero")
	}
	return unicode.IsUpper(rune(name[0]))
}
