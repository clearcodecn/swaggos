package yidoc

import (
	"fmt"
	"reflect"
)

type fieldType interface {
	Basic() ArgumentType
}

func getFieldType(v reflect.Value) (fieldType, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()

	var ft fieldType
	switch typ.Kind() {
	case reflect.Bool:
		ft = newBoolean()
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.UnsafePointer, reflect.Uintptr:
		ft = newInteger(fieldFormats[v.Kind()])
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		ft = newNumber(fieldFormats[v.Kind()])
	case reflect.Slice:
		ft = newArray(reflect.New(typ.Elem()).Elem())
	case reflect.Struct:
		ft = newObject(reflect.New(typ).Elem())
	default:
		return nil, fmt.Errorf("un support type: %s", typ.Kind())
	}
	return ft, nil
}

var fieldFormats = map[reflect.Kind]string{
	reflect.Bool:    "boolean",
	reflect.Int:     "int32",
	reflect.Int8:    "int32",
	reflect.Int16:   "int32",
	reflect.Int32:   "int32",
	reflect.Int64:   "int64",
	reflect.Uint:    "int32",
	reflect.Uint8:   "int32",
	reflect.Uint16:  "int32",
	reflect.Uint32:  "int32",
	reflect.Uint64:  "int64",
	reflect.Uintptr: "int32",
	reflect.Float32: "float",
	reflect.Float64: "float",
}

type integer struct{ format string }

func (i *integer) Basic() ArgumentType { return TypeInteger }

func newInteger(format string) fieldType { return &integer{format: format} }

type number struct{ format string }

func (i *number) Basic() ArgumentType { return TypeNumber }

func newNumber(format string) fieldType { return &integer{format: format} }

type boolean struct{ format string }

func (i *boolean) Basic() ArgumentType { return TypeBoolean }

func newBoolean() fieldType { return &boolean{format: "boolean"} }

type array struct{ itemValue reflect.Value }

func (a *array) Basic() ArgumentType { return TypeArray }

func newArray(v reflect.Value) fieldType { return &array{itemValue: v} }

type object struct{ value reflect.Value }

func (o *object) Basic() ArgumentType     { return TypeObject }
func newObject(v reflect.Value) fieldType { return &object{value: v} }
