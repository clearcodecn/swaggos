package yidoc

//
//import (
//	"github.com/go-openapi/spec"
//	"reflect"
//)
//
//type ArgumentType string
//
//const (
//	TypeString  ArgumentType = "string"
//	TypeNumber               = "number"
//	TypeInteger              = "integer"
//	TypeBoolean              = "boolean"
//	TypeArray                = "array"
//	TypeFile                 = "file"
//	TypeObject               = "object"
//)
//
//type ArgPosition string
//
//func (a ArgPosition) TagName() spec.StringOrArray {
//	switch a {
//	case PosHeader:
//		return spec.StringOrArray{"header"}
//	case PosQuery:
//		return spec.StringOrArray{"query,json"}
//	case PosBody:
//		return spec.StringOrArray{"json"}
//	case PosFormData:
//		return spec.StringOrArray{"form", "json"}
//	default:
//		return nil
//	}
//}
//
//const (
//	PosHeader   = "header"
//	PosQuery    = "query"
//	PosPath     = "path"
//	PosFormData = "formData"
//	PosBody     = "body"
//)
//
//type Arg interface {
//	Name() string
//	Type() ArgumentType
//	Description() string
//	Required() bool
//	In() ArgPosition
//	BuildRef() *spec.Schema
//}
//
//type emptyArg struct{}
//
//func (e *emptyArg) Name() string           { return "" }
//func (e *emptyArg) Type() ArgumentType     { return "" }
//func (e *emptyArg) Description() string    { return "" }
//func (e *emptyArg) Required() bool         { return false }
//func (e *emptyArg) In() ArgPosition        { return "" }
//func (e *emptyArg) BuildRef() *spec.Schema { return nil }
//
//type arrayAble struct {
//	itemType string
//}
//
//type queryArgument struct {
//	emptyArg
//	arrayAble
//	name     string
//	typ      ArgumentType
//	required bool
//	desc     string
//}
//
//func (a *queryArgument) Name() string { return a.name }
//
//func (a *queryArgument) Type() ArgumentType { return a.typ }
//
//func (a *queryArgument) Required() bool { return a.required }
//
//func (a *queryArgument) In() ArgPosition { return PosQuery }
//
//func QueryString(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &queryArgument{name: name, typ: TypeString, required: required, desc: desc}
//}
//
//func QueryNumber(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &queryArgument{name: name, typ: TypeNumber, required: required, desc: desc}
//}
//
//// query array string
//func QueryStringArray(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &queryArgument{name: name, typ: TypeArray, arrayAble: arrayAble{itemType: "string"}, required: required, desc: desc}
//}
//
//// query array string
//func QueryNumberArray(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &queryArgument{name: name, typ: TypeArray, arrayAble: arrayAble{itemType: "number"}, required: required, desc: desc}
//}
//
//// path start.
//type pathArgument struct {
//	emptyArg
//	name     string
//	typ      ArgumentType
//	required bool
//	desc     string
//}
//
//func (p *pathArgument) Type() ArgumentType { return p.typ }
//
//func (p *pathArgument) In() ArgPosition { return PosPath }
//
//func (p *pathArgument) Name() string   { return p.name }
//func (p *pathArgument) Required() bool { return p.required }
//
//func (p *pathArgument) Description() string { return p.desc }
//
//func PathString(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &pathArgument{name: name, typ: TypeString, required: required, desc: desc}
//}
//
//func PathNumber(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &pathArgument{name: name, typ: TypeString, required: required, desc: desc}
//}
//
//// form start
//type formArgument struct {
//	emptyArg
//	name     string
//	typ      ArgumentType
//	required bool
//	desc     string
//}
//
//func (f *formArgument) Type() ArgumentType { return f.typ }
//
//func (f *formArgument) In() ArgPosition { return PosFormData }
//
//func (f *formArgument) Name() string   { return f.name }
//func (f *formArgument) Required() bool { return f.required }
//
//func (f *formArgument) Description() string { return f.desc }
//
//func FormString(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &formArgument{name: name, typ: TypeString, required: required, desc: desc}
//}
//
//func FormNumber(name string, descAndRequired ...interface{}) Arg {
//	desc, required := getDescAndRequired(descAndRequired)
//	return &formArgument{name: name, typ: TypeNumber, required: required, desc: desc}
//}
//
//// body start
//type bodyArgument struct {
//	object interface{}
//}
//
//func (b *bodyArgument) Name() string {
//	return reflect.TypeOf(b.object).Name()
//}
//
//func (b *bodyArgument) Type() ArgumentType {
//	return ""
//}
//
//func (b *bodyArgument) Description() string {
//	return ""
//}
//
//func (b *bodyArgument) Required() bool {
//	return false
//}
//
//func (b *bodyArgument) In() ArgPosition {
//	return PosBody
//}
//
//func (b *bodyArgument) BuildRef() *spec.Schema {
//	return &spec.Schema{
//		SchemaProps: spec.SchemaProps{
//			Ref:         spec.Ref{},
//			Description: "",
//			Type:        b.getType(b.object),
//			Items:       nil, // for array
//			Properties:  nil, // for object
//		},
//		ExtraProps: nil,
//	}
//}
//
//func (b *bodyArgument) getType(object interface{}) spec.StringOrArray {
//	val := reflect.ValueOf(object)
//	if val.Kind() == reflect.Ptr {
//		val = val.Elem()
//	}
//	typ := val.Type()
//	var objTyp string
//	switch typ.Kind() {
//	case reflect.Struct:
//		objTyp = TypeObject
//	case reflect.Slice:
//		objTyp = TypeArray
//	default:
//		//objTyp = TypeString
//	}
//	return spec.StringOrArray{objTyp}
//}
