package swaggos

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/spec"
)

// Define defines a object or a array to swagger definitions area.
// it will find all sub-items and build them into swagger tree.
// it returns the definitions ref.
func (swaggos *Swaggos) Define(v interface{}) spec.Ref {
	schema := swaggos.buildSchema(v)
	return swaggos.addDefinition(v, schema)
}

// addDefinition add a definition to swagger definitions.
// the name will get from the given type.
// if name's name is repeated, will add package path prefix to the name until name is unique.
func (swaggos *Swaggos) addDefinition(t interface{}, v spec.Schema) spec.Ref {
	var (
		name string
		typ  reflect.Type
	)
	switch tt := t.(type) {
	case reflect.Type:
		typ = tt
	case reflect.Value:
		typ = tt.Type()
	default:
		typ = reflect.TypeOf(t)
		if typ == nil {
			typ = reflect.TypeOf(new(interface{}))
		} else {
			typ = reflect.Indirect(reflect.ValueOf(t)).Type()
		}
		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			name = fmt.Sprintf("%sArray", typ.Elem().Name())
		}
	}
	if name == "" {
		name = typ.Name()
	}
	if swaggos.typeNames == nil {
		swaggos.typeNames = make(map[reflect.Type]string)
	}
	if name, ok := swaggos.typeNames[typ]; ok {
		return definitionRef(name)
	}
	pkgPath := pkgPath(typ)
	subName := name
	i := 1
	for {
		// create a newName. like pkgName
		if _, ok := swaggos.definitions[subName]; ok {
			var prefix string
			if len(pkgPath)-i > 0 {
				prefix = pkgPath[len(pkgPath)-i]
			} else {
				i++
			}
			subName = fmt.Sprintf("%s.%s", prefix, name)
		} else {
			name = subName
			break
		}
	}
	swaggos.definitions[name] = v
	swaggos.typeNames[typ] = name
	return definitionRef(name)
}

func (swaggos *Swaggos) buildSchema(v interface{}) spec.Schema {
	typ := reflect.TypeOf(v)
	// if given nil interface{}, typ is nil, then we return a empty object schema
	if typ == nil {
		return emptyObjectSchema()
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if isBasicType(typ) {
		return basicSchema(typ)
	}
	switch typ.Kind() {
	case reflect.Array, reflect.Slice:
		elTyp := typ.Elem()
		if elTyp.Kind() == reflect.Ptr {
			elTyp = elTyp.Elem()
		}
		elVal := reflect.New(elTyp).Elem()

		// basicArray
		if isBasicType(elTyp) {
			schema := basicSchema(elTyp)
			return arraySchema(&schema)
		}

		// structArray
		if elTyp.Kind() == reflect.Struct {
			schema := swaggos.buildSchema(elVal.Interface())
			ref := swaggos.addDefinition(elVal, schema)
			return arraySchemaRef(ref)
		}
		var arraySchema = emptyArray()
		childType, childSchema := arrayProps(elTyp, &arraySchema)
		if isBasicType(childType) {
			basic := basicSchema(childType)
			childSchema.Items = &spec.SchemaOrArray{
				Schema: &basic,
			}
		} else {
			schema := swaggos.buildSchema(reflect.New(childType).Elem().Interface())
			ref := swaggos.addDefinition(childType, schema)
			childSchema.Items = refArraySchema(ref)
		}
		return arraySchema
	case reflect.Struct:
		return swaggos.buildStructSchema(v)
	case reflect.Map, reflect.Interface:
		// TODO:: handle map schema
		return emptyObjectSchema()
	}
	return emptyObjectSchema()
}

// val is struct value
func (swaggos *Swaggos) buildStructSchema(v interface{}) spec.Schema {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	// if given nil interface{}, typ is nil, then we return a empty object schema
	if typ == nil {
		return emptyObjectSchema()
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = reflect.Indirect(reflect.New(typ))
	}
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("interface{} is not struct: %T", v))
	}
	var schema spec.Schema
	schema.Properties = make(spec.SchemaProperties)
	schema.Type = spec.StringOrArray{_Object}
	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()
		if !isExport(typ.Field(i).Name) {
			continue
		}
		tg := newTags(typ.Field(i).Tag)
		if tg.ignore() {
			continue
		}
		// Anonymous field, if there are json tag on the field. then we say it's a object reference.
		// if it's basic type, add it to properties directly, else we build it.
		if typ.Field(i).Anonymous && tg.jsonTag() == "" {
			if isBasicType(fieldType) {
				prop := basicSchema(fieldType)
				schema.Properties[typ.Field(i).Name] = prop
				continue
			} else {
				// TODO:: if the field is a array type. what should we do here???
				fieldSchema := swaggos.buildSchema(field.Interface())
				for name, val := range fieldSchema.Properties {
					schema.Properties[name] = val
				}
				schema.Required = append(schema.Required, fieldSchema.Required...)
			}
			continue
		}
		fieldName := typ.Field(i).Name
		if name := tg.jsonName(); name != "" {
			fieldName = name
		}
		if tg.required() {
			schema.Required = append(schema.Required, fieldName)
		}
		var prop spec.Schema
		if isBasicType(fieldType) {
			prop = basicSchema(fieldType)
		} else {
			prop = swaggos.buildSchema(field.Interface())
		}
		prop = tg.mergeSchema(prop)
		schema.Properties[fieldName] = prop
	}
	return schema
}
