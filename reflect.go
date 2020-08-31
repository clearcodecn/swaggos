package yidoc

import (
	"github.com/go-openapi/spec"
	. "reflect"
	"strconv"
	"strings"
)

const (
	dataTypeString  = "string"
	dataTypeNumber  = "number"
	dataTypeBoolean = "boolean"
	dataTypeInteger = "integer"
	dataTypeFile    = "file"
	dataTypeArray   = "array"
	dataTypeObject  = "object"
)

// 1. name 可以为空的几种情况:
//		* in = 'body' => 具体的参数，根据type来看，
//			* 如果 type = 'Object', 那么 参数就在 Object 字段上面，Object 字段是递归类型
//			* 如果 type = 'array',  那么 具体的参数在 items 上面，根据items即可得出需要的字段类型.
// 			* 如果 解析的 v 是基本类型(没有名字和变量的话,)
type Object struct {
	Typ          Type
	Val          Value
	StructField  StructField
	ObjectFields []*Object
	ArrayFields  interface{}

	IsNull   bool
	DataType string
}

func parseObject(v interface{}) *Object {
	obj := new(Object)
	val := Indirect(ValueOf(v))
	obj.Val = val
	obj.Typ = val.Type()

	if typ := getBasicType(val); typ != nil {
		return typ
	}
	switch val.Kind() {
	case Map, Interface:
		obj.IsNull = true
		obj.DataType = dataTypeObject
	case Array, Slice:
		obj.DataType = dataTypeArray
		elem := New(val.Type().Elem()).Elem()
		obj.ArrayFields = parseArrayFields(elem)
	case Struct:
		obj.DataType = dataTypeObject
		obj.ObjectFields = parseStructFields(val)
	default:
		obj = nil
	}
	return obj
}

func getBasicType(val Value) (*Object, ) {
	var obj = new(Object)
	typ := val.Type()
	obj.Val = val
	obj.Typ = typ
	switch typ.Kind() {
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		obj.DataType = dataTypeInteger
	case Float32, Float64:
		obj.DataType = dataTypeNumber
	case Bool:
		obj.DataType = dataTypeBoolean
	case String:
		obj.DataType = dataTypeString
	default:
		obj = nil
	}
	return obj
}

// must return an array Object
func parseArrayFields(val Value) []*Object {
	var objects []*Object
	typ := val.Type()
	if basic := getBasicType(val); basic != nil {
		objects = append(objects, basic)
		return objects
	}
	switch typ.Kind() {
	case Map, Interface:
		objects = append(objects, &Object{
			Typ:      typ,
			Val:      val,
			IsNull:   true,
			DataType: "Object",
		})
	case Struct:
		objects = append(objects, parseStructFields(val)...)
	case Array, Slice:
		objs := parseArrayFields(New(typ.Elem()).Elem())
		objects = append(objects, &Object{
			Typ:         typ,
			Val:         val,
			DataType:    "array",
			ArrayFields: objs,
		})
	}
	return objects
}

func parseStructFields(v Value) []*Object {
	typ := Indirect(v).Type()
	var objects []*Object
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldVal := v.Field(i)
		if fieldType.Anonymous {
			if basic := getBasicType(fieldVal); basic != nil {
				basic.StructField = fieldType
				objects = append(objects, basic)
			} else {
				switch fieldVal.Kind() {
				case Map, Interface:
					objects = append(objects, &Object{
						Typ:      fieldVal.Type(),
						Val:      fieldVal,
						IsNull:   true,
						DataType: "Object",
					})
				case Array, Slice:
					objects = append(objects, parseArrayFields(New(fieldVal.Type().Elem()).Elem())...)
				case Struct:
					objects = append(objects, parseStructFields(fieldVal)...)
				default:
				}
			}
			continue
		}

		if basic := getBasicType(fieldVal); basic != nil {
			basic.StructField = fieldType
			objects = append(objects, basic)
		} else {
			switch fieldVal.Kind() {
			case Map, Interface:
				objects = append(objects, &Object{
					Typ:         fieldVal.Type(),
					Val:         fieldVal,
					IsNull:      true,
					DataType:    dataTypeObject,
					StructField: fieldType,
				})
			case Array, Slice:
				t := fieldVal.Type().Elem()
				if t.Kind() == Ptr {
					t = t.Elem()
				}
				obj := &Object{
					Typ:         fieldVal.Type(),
					Val:         fieldVal,
					DataType:    dataTypeArray,
					ArrayFields: parseArrayFields(New(t).Elem()),
					StructField: fieldType,
				}
				objects = append(objects, obj)
			case Struct:
				obj := &Object{
					Typ:          fieldVal.Type(),
					Val:          fieldVal,
					DataType:     dataTypeObject,
					ObjectFields: parseStructFields(fieldVal),
					StructField:  fieldType,
				}
				objects = append(objects, obj)
			default:
			}
		}
	}

	return objects
}

func (o *Object) buildSchema(doc *YiDoc) spec.SchemaProps {
	prop := spec.SchemaProps{
		ID: strconv.FormatInt(nextId(), 10),
	}
	if o.IsNull {
		return prop
	}
	switch o.DataType {
	case dataTypeInteger, dataTypeString, dataTypeBoolean, dataTypeNumber, dataTypeFile:
		prop.Type = spec.StringOrArray{o.DataType}
	case dataTypeArray:
		switch fieldType := o.ArrayFields.(type) {
		case *Object:
			prop.Items = &spec.SchemaOrArray{
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{fieldType.DataType},
					},
				},
			}
		case []*Object:
			var schemas []spec.Schema
			for _, ft := range fieldType {
				props := ft.buildSchema(doc)
				schemas = append(schemas, spec.Schema{
					SchemaProps: props,
				})
			}
			prop.Items = &spec.SchemaOrArray{
				Schemas: schemas,
			}
		}
	case dataTypeObject:
		prop.Properties = make(spec.SchemaProperties)
		for _, obj := range o.ObjectFields {
			name := getName(obj.Typ, obj.StructField)
			schema := spec.Schema{}
			if obj.DataType == dataTypeObject {
				prop := o.buildSchema(doc)
				ref := doc.addModel(o.Typ, prop)
				schema.SchemaProps.Ref = spec.MustCreateRef(ref)
				prop.Properties[name] = schema
			} else if obj.DataType == dataTypeArray {
				schema.Type = spec.StringOrArray{obj.DataType}
				prop = obj.buildSchema(doc)
			} else {
				schema.Type = spec.StringOrArray{obj.DataType}
			}
		}
	}
	return prop
}

func getName(typ Type, field StructField) string {
	name := field.Name
	jsonName := getTagName(field.Tag, "json")
	if jsonName != "" {
		name = jsonName
	}
	return name
}

func getTagName(tag StructTag, name string) string {
	data := tag.Get(name)
	return strings.Split(data, ",")[0]
}
