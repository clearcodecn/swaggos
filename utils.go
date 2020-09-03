package yidoc

import (
	"github.com/go-openapi/spec"
	"reflect"
)

const (
	String  = "string"
	Number  = "number"
	Integer = "integer"
	Boolean = "boolean"
	Array   = "array"
	File    = "file"
	Object  = "object"
)

const (
	Int32    = "int32"
	Int64    = "int64"
	Float    = "float"
	Double   = "double"
	Byte     = "byte"
	Binary   = "binary"
	Date     = "date"
	DateTime = "date-time"
	Password = "password"
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
			Type:  spec.StringOrArray{Array},
			Items: &spec.SchemaOrArray{},
		},
	}
}

func arraySchema(schema *spec.Schema) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{Array},
			Items: &spec.SchemaOrArray{
				Schema: schema,
			},
		},
	}
}

func arraySchemaRef(ref spec.Ref) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{Array},
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
				Type: spec.StringOrArray{Boolean},
			},
		}
	case reflect.Int, reflect.Int64, reflect.Int8, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.UnsafePointer:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:   spec.StringOrArray{Integer},
				Format: string(Int64),
			},
		}
	case reflect.Float64, reflect.Float32:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:   spec.StringOrArray{Number},
				Format: string(Float),
			},
		}
	case reflect.String:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{String},
			},
		}
	default:
		return emptyObjectSchema()
	}
}
