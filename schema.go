package yidoc

import (
	"fmt"
	"github.com/go-openapi/spec"
	"reflect"
	"strings"
)

type schemaBuilder struct {
	schemas map[string]*spec.SchemaProps
}

// Build return refString
func (sb *schemaBuilder) Build(v reflect.Value, pos ArgPosition) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()

	ft, err := getFieldType(v)
	if err != nil {
		return ""
	}
	switch ft.Basic() {
	case TypeObject:
		name := sb.buildFromStruct(v, pos)
		return name
	case TypeArray:
		id := sb.newId(typ)
		_, err := sb.parseArray(v, pos)
		if err != nil {
			return ""
		}
		return id
	}

	switch typ.Kind() {
	case reflect.Struct:
		sb.buildFromStruct(v, pos)
	}
}

func (sb *schemaBuilder) buildFromStruct(v reflect.Value, pos ArgPosition) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()
	prop := &spec.SchemaProps{
		ID: sb.newId(typ),
	}
	defer func() {
		sb.schemas[prop.ID] = prop
	}()
	if typ.NumField() == 0 {
		return prop.ID
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := v.Field(i)
		tags := pos.TagName()
		name, ignore := mustParseTag(field, tags)
		if ignore {
			continue
		}
		if !fieldVal.CanSet() {
			continue
		}
		if field.Anonymous {
			valProp := spec.SchemaProps{}
			any, err := sb.buildFromAnonymous(fieldVal, pos)
			if err != nil {
				continue
			}
			ft, err := getFieldType(fieldVal)
			switch any.typ {
			case TypeObject:
				for k, v := range any.properties {
					prop.Properties[k] = v
				}
				continue
			case TypeNumber:
				valProp.Type = spec.StringOrArray{"number"}
				valProp.Format = ft.(*number).format
			case TypeInteger:
				valProp.Type = spec.StringOrArray{"number"}
				valProp.Format = ft.(*integer).format
			case TypeBoolean:
				valProp.Type = spec.StringOrArray{"boolean"}
				valProp.Format = "boolean"
			case TypeArray:
				valProp = any.item
			}
			prop.Properties[name] = spec.Schema{SchemaProps: valProp}
			continue
		}
		ft, err := getFieldType(fieldVal)
		if err != nil {
			continue
		}
		valProp := spec.SchemaProps{}
		switch ft.Basic() {
		case TypeNumber:
			valProp.Type = spec.StringOrArray{"number"}
			valProp.Format = ft.(*number).format
		case TypeInteger:
			valProp.Type = spec.StringOrArray{"integer"}
			valProp.Format = ft.(*integer).format
		case TypeBoolean:
			valProp.Type = spec.StringOrArray{"boolean"}
			valProp.Format = "boolean"
		case TypeObject:
			valProp.Ref = spec.MustCreateRef(sb.buildFromStruct(reflect.New(field.Type).Elem(), pos))
			// create ref.
		case TypeArray:
			valProp, err = sb.parseArray(fieldVal, pos)
			if err != nil {
				continue
			}
		}
		prop.Properties[name] = spec.Schema{SchemaProps: valProp}
	}
	return prop.ID
}

func (sb *schemaBuilder) parseArray(val reflect.Value, pos ArgPosition) (spec.SchemaProps, error) {
	var prop = spec.SchemaProps{}
	prop.Type = spec.StringOrArray{"array"}
	item := reflect.New(val.Type().Elem()).Elem()
	itemFieldType, err := getFieldType(item)
	if err != nil {
		return spec.SchemaProps{}, err
	}
	switch itemFieldType.Basic() {
	case TypeNumber:
		prop.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"number"}}}}
	case TypeInteger:
		prop.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}}}
	case TypeBoolean:
		prop.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"boolean"}}}}
	case TypeArray:
		sp, err := sb.parseArray(item, pos)
		if err != nil {
			return spec.SchemaProps{}, err
		}
		prop.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: sp}}
	case TypeObject:
		ref := sb.buildFromStruct(item, pos)
		prop.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: spec.SchemaProps{Ref: spec.MustCreateRef(ref)}}}
	}
	return prop, nil
}

func mustParseTag(field reflect.StructField, names []string) (name string, ignore bool) {
	for _, n := range names {
		name = field.Tag.Get(n)
		if name == "" {
			continue
		}
		if name == "-" {
			ignore = true
		}
		return name, ignore
	}
	return field.Name, false
}

type anonymous struct {
	typ ArgumentType

	// for string
	item spec.SchemaProps

	// for struct
	properties spec.SchemaProperties
}

func (sb *schemaBuilder) buildFromAnonymous(v reflect.Value, pos ArgPosition) (*anonymous, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	ft, err := getFieldType(v)
	if err != nil {
		return nil, err
	}
	any := new(anonymous)
	any.typ = ft.Basic()

	switch ft.Basic() {
	case TypeObject:
		id := sb.buildFromStruct(v, pos)
		any.properties = sb.schemas[id].Properties
	case TypeArray:
		sp, err := sb.parseArray(v, pos)
		if err != nil {
			return nil, err
		}
		any.item = sp
	}

	return any, nil
}

func (sb *schemaBuilder) newId(typ reflect.Type) string {
	if _, ok := sb.schemas[typ.Name()]; !ok {
		return typ.Name()
	}

	pkg := strings.Split(typ.PkgPath(), "/")
	prefix := pkg[len(pkg)-1]

	return strings.NewReplacer(".", "", "_", "").Replace(fmt.Sprintf("%s%s", prefix, typ.Name()))
}
