package doc

import (
	"github.com/go-openapi/spec"
	"reflect"
)

const (
	tagQuery = "query"
	tagForm  = "form"
	tagJson  = "json" // body

	typeNumber = "number"
	typeBool   = "bool"
	typeString = "string"
	typeArray  = "array"
	typeObject = "object"
)

func parseObject(val reflect.Value) []BasicParam {
	var res []BasicParam
	var param BasicParam
	typ := val.Type()

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		param.Type = typeNumber
	case reflect.Bool:
		param.Type = typeBool
	case reflect.String:
		param.Type = typeString
	case reflect.Slice:
		param.Type = typeArray
		el := typ.Elem()
		var t reflect.Type
		switch el.Kind() {
		case reflect.Ptr, reflect.Slice:
			t = el.Elem()
		case reflect.Struct, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool:
			t = el
		default:
			t = reflect.TypeOf(nil)
		}
		data := parseObject(reflect.New(t).Elem())
		param.Items = data
	case reflect.Struct:
		param.Type = typeObject
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldVal := val.Field(i)
			if !fieldVal.CanSet() {
				continue
			}
			if field.Anonymous {
				if fieldVal.Kind() == reflect.Ptr {
					fieldVal = reflect.New(field.Type.Elem()).Elem()
				}
				data := parseObject(fieldVal)
				res = append(res, data...)
				continue
			}
			val,ok := field.Tag.Lookup("json")
			if ok {

			}
		}
	}
}

func getProps(val reflect.Value) spec.SchemaProps {
	typ := val.Type()
	prop := spec.SchemaProps{}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		prop.Type = spec.StringOrArray{"number"}
	case reflect.Bool:
		prop.Type = spec.StringOrArray{"boolean"}
	case reflect.String:
		prop.Type = spec.StringOrArray{"string"}
	case reflect.Slice:
		prop.Type = spec.StringOrArray{"array"}
		el := typ.Elem()
		var t reflect.Type
		switch el.Kind() {
		case reflect.Ptr, reflect.Slice:
			t = el.Elem()
		case reflect.Struct, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool:
			t = el
		default:
			t = reflect.TypeOf(nil)
		}
		props := getProps(reflect.New(t).Elem())
		prop.Items = &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: props,
			},
		}
	case reflect.Struct:
		prop.Type = spec.StringOrArray{"object"}
		prop.Properties = make(map[string]spec.Schema)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldVal := val.Field(i)
			if !fieldVal.CanSet() {
				continue
			}
			if field.Anonymous {
				if fieldVal.Kind() == reflect.Ptr {
					fieldVal = reflect.New(field.Type.Elem()).Elem()
				}
				props := getProps(fieldVal)
				for k, v := range props.Properties {
					prop.Properties[k] = v
				}
				continue
			}
			key := field.Tag.Get("json")
			key, required := keyName(key, field.Name)
			if key == "" {
				continue
			}
			prop.Properties[key] = spec.Schema{
				SchemaProps: getProps(fieldVal),
			}
			if required {
				prop.Required = append(prop.Required, key)
			}
		}
	default:
		prop.Type = spec.StringOrArray{"object"}
	}
	return prop
}
