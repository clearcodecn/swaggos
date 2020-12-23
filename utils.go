package swaggos

import (
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
)

const (
	_String  = "string"
	_Number  = "number"
	_Integer = "integer"
	_Boolean = "boolean"
	_Array   = "array"
	_File    = "file"
	_Object  = "object"
)

const (
	_Int32    = "int32"
	_Int64    = "int64"
	_Float    = "float"
	_Double   = "double"
	_Byte     = "byte"
	_Binary   = "binary"
	_Date     = "date"
	_DateTime = "date-time"
	_Password = "password"
)

func definitionRef(name string) spec.Ref {
	return spec.MustCreateRef("#/definitions/" + name)
}

func emptyObjectSchema() spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"object"},
		},
	}
}

func emptyArray() spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:  spec.StringOrArray{_Array},
			Items: &spec.SchemaOrArray{},
		},
	}
}

func arraySchema(schema *spec.Schema) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{_Array},
			Items: &spec.SchemaOrArray{
				Schema: schema,
			},
		},
	}
}

func arraySchemaRef(ref spec.Ref) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{_Array},
			Items: &spec.SchemaOrArray{
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Ref: ref,
					},
				},
			},
		},
	}
}

func refSchema(ref spec.Ref) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: ref,
		},
	}
}

func refArraySchema(ref spec.Ref) *spec.SchemaOrArray {
	schema := refSchema(ref)
	return &spec.SchemaOrArray{
		Schema: &schema,
	}
}

func isExport(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}

// arrayProps find the deep type in a array loop,
// it returns the deep child schema pointer and deep child type
// also build tree in given schema
func arrayProps(typ reflect.Type, schema *spec.Schema) (reflect.Type, *spec.Schema) {
	var childSchema = schema
	for typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		typ = typ.Elem()
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		arraySchema := emptyArray()
		childSchema = &arraySchema
		schema.SchemaProps.Items = &spec.SchemaOrArray{
			Schema: &arraySchema,
		}
		elVal := reflect.New(typ).Elem()
		return arrayProps(elVal.Type(), childSchema)
	}
	return typ, childSchema
}

// isBasicType returns if the given type is a basic type in:
// number,string,boolean,integer.
func isBasicType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.UnsafePointer, reflect.Float64, reflect.Float32:
		return true
	default:
		return false
	}
}

func basicSchema(typ reflect.Type) spec.Schema {
	switch typ.Kind() {
	case reflect.Bool:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{_Boolean},
			},
		}
	case reflect.Int, reflect.Int64, reflect.Int8, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.UnsafePointer:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:   spec.StringOrArray{_Integer},
				Format: string(_Int64),
			},
		}
	case reflect.Float64, reflect.Float32:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:   spec.StringOrArray{_Number},
				Format: string(_Float),
			},
		}
	case reflect.String:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{_String},
			},
		}
	default:
		return emptyObjectSchema()
	}
}

func pkgPath(typ reflect.Type) []string {
	switch typ.Kind() {
	case reflect.Slice, reflect.Array, reflect.Ptr:
		typ = typ.Elem()
	}
	return strings.Split(typ.PkgPath(), "/")
}
