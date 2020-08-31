package v2

import (
	"fmt"
	"github.com/go-openapi/spec"
	"reflect"
)

type YiDoc struct {
	definitions spec.Definitions
	paths       map[string]map[string]*Path

	securityDefinitions spec.SecurityDefinitions
	security            []map[string][]string
	info                spec.InfoProps
	host                string
	basePath            string
}

func (y *YiDoc) Define(name string, v interface{}) spec.Ref {
	if y.definitions == nil {
		y.definitions = make(spec.Definitions)
	}
	schema := y.buildSchema(v)
	return y.addDefine(name, schema)
}

func (y *YiDoc) addDefine(name string, v spec.Schema) spec.Ref {
	if y.defExist(name) {
		panic(fmt.Errorf("repeated: %s definition", name))
	}
	y.definitions[name] = v
	return spec.MustCreateRef("#/definitions/" + name)
}

func (y *YiDoc) defExist(name string) bool {
	if y.definitions == nil {
		y.definitions = make(spec.Definitions)
	}
	// def 都么有，那么package就不可能重复
	if _, ok := y.definitions[name]; !ok {
		return false
	}
	return false
}

func (y *YiDoc) buildSchema(v interface{}) spec.Schema {
	typ := reflect.TypeOf(v)
	if typ == nil {
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
			},
		}
	}
	if typ.Kind() == reflect.Ptr {
		v = reflect.New(reflect.TypeOf(v).Elem()).Interface()
	}
	val := reflect.Indirect(reflect.ValueOf(v))
	typ = val.Type()

	if isBasicType(typ) {
		return getBasicSchema(typ)
	}

	switch typ.Kind() {
	case reflect.Array, reflect.Slice:
		elVal := reflect.New(typ.Elem()).Elem()
		elTyp := elVal.Type()
		if isBasicType(elTyp) {
			schema := getBasicSchema(elTyp)
			return spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &schema,
					},
				},
			}
		}
		if elTyp.Kind() == reflect.Struct {
			schema := y.buildSchema(elVal.Interface())
			return spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &schema,
					},
				},
			}
		}

		if elTyp.Kind() == reflect.Slice || elTyp.Kind() == reflect.Array {
			val := reflect.New(elTyp.Elem()).Elem()
			schema := y.buildSchema(val.Interface())
			return spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &schema,
					},
				},
			}
		}

		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"array"},
				Items: &spec.SchemaOrArray{
					Schema: &spec.Schema{},
				},
			},
		}
	case reflect.Struct:
		return y.buildStructSchema(v)
	case reflect.Map, reflect.Interface:
		return spec.Schema{}
	}
	return spec.Schema{}
}

// val is struct value
func (y *YiDoc) buildStructSchema(v interface{}) spec.Schema {
	val := reflect.Indirect(reflect.ValueOf(v))
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("interface{} is not struct: %T", v))
	}
	var schema spec.Schema
	schema.Properties = make(spec.SchemaProperties)
	schema.Type = spec.StringOrArray{"object"}
	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()
		if !isExport(typ.Field(i).Name) {
			continue
		}
		// 匿名字段.
		// 检查是否有json tag，如果没有，则合并这些字段到struct上面.
		if typ.Field(i).Anonymous && typ.Field(i).Tag.Get("json") == "" {
			if isBasicType(fieldType) {
				prop := getBasicSchema(fieldType)
				schema.Properties[typ.Field(i).Name] = prop
				continue
			} else {
				// TODO:: merge attributes in one.
				fieldSchema := y.buildSchema(field.Interface())
				for name, val := range fieldSchema.Properties {
					schema.Properties[name] = val
				}
				schema.Required = append(schema.Required, fieldSchema.Required...)
			}
			continue
		}
		tg := newTags(typ.Field(i).Tag)
		if tg.ignore() {
			continue
		}
		name := tg.name()
		if name == "" {
			name = typ.Field(i).Name
		}
		required := tg.required()
		if required {
			schema.Required = append(schema.Required, name)
		}
		if isBasicType(fieldType) {
			prop := getBasicSchema(fieldType)
			schema.Properties[name] = prop
			continue
		} else {
			fieldSchema := y.buildSchema(field.Interface())
			ref := y.addDefine(getFieldTypeName(typ.Field(i).Type), fieldSchema)
			schema.Properties[name] = spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			}
		}

	}
	return schema
}

func isBasicType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.UnsafePointer, reflect.Float64, reflect.Float32:
		return true
	default:
		return false
	}
}

func getBasicSchema(typ reflect.Type) spec.Schema {
	switch typ.Kind() {
	case reflect.Bool:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"boolean"},
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
		return spec.Schema{}
	}
}

func getFieldTypeName(typ reflect.Type) string {
	return typ.String()
}

func isExport(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}
