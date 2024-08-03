package typings

import (
	"reflect"
)

func IsMap(input interface{}) bool {
	if input == nil {
		return false
	}

	val := reflect.ValueOf(input)

	return val.Kind() == reflect.Map
}

func IsPointerToStruct(v interface{}) bool {
	if v == nil {
		return false
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Kind() == reflect.Struct
	}
	return false
}

func IsStruct(v interface{}) bool {
	if v == nil {
		return false
	}
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Struct
}

func IsZeroValue(i interface{}) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return v.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.String:
		return v.String() == ""
	default:
		return false
	}
}
