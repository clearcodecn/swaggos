package yidoc

//import (
//	"encoding/json"
//	"gopkg.in/yaml.v2"
//	"reflect"
//	"strings"
//)
//
//const (
//	TypeString = "string"
//	TypeNumber = "number"
//	TypeFile   = "file"
//	TypeBool   = "bool"
//	TypeObject = "object"
//	TypeArray  = "array"
//)
//
//type Type string
//
//type Item struct {
//	Name        string `json:"name" yaml:"name"`
//	Type        Type   `json:"type" yaml:"type"`
//	Required    bool   `json:"required" yaml:"required"`
//	Description string `json:"description" yaml:"description"`
//
//	Item   *Item   `json:"items" yaml:"items"` // when type is array, items will show this field.
//	Object []*Item `json:"object" yaml:"objects"`
//}
//
//type Ref string
//
//type Path struct {
//	Path                string         `json:"path" yaml:"path"`
//	Method              string         `json:"method" yaml:"method"`
//	Headers             []Item         `json:"headers" yaml:"headers"`
//	Body                Ref            `json:"body" yaml:"body"`
//	Query               Ref            `json:"query" yaml:"query"`
//	Form                Ref            `json:"form" yaml:"form"`
//	FormData            Ref            `json:"formData" yaml:"formData"`
//	Response            map[string]Ref `json:"response" yaml:"response"`
//	ContentType         string         `json:"contentType" yaml:"contentType"`
//	ResponseContentType string         `json:"responseContentType" yaml:"responseContentType"`
//}
//
//type YiDoc struct {
//	BaseURL   string              `json:"baseURL" yaml:"baseURL"`
//	Version   string              `json:"version" yaml:"version"`
//	Headers   []Item              `json:"headers" yaml:"headers"`
//	Responses map[int]Ref         `json:"responses" yaml:"responses"`
//	Models    []map[string][]Item `json:"models" yaml:"models"`
//	Paths     []Path              `json:"paths" yaml:"paths"`
//}
//
//func (y *YiDoc) buildJson() ([]byte, error) {
//	return json.MarshalIndent(y, "", "\t")
//}
//
//func (y *YiDoc) buildYaml() ([]byte, error) {
//	return yaml.Marshal(y)
//}
//
//func (y *YiDoc) Response(code int, v interface{}) Response {
//	var resp Response
//
//	val := reflect.ValueOf(v)
//	switch val.Kind() {
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
//		resp = &basicItemResponse{
//			item: &Item{Type: TypeNumber},
//		}
//	case reflect.Bool:
//		resp = &basicItemResponse{
//			item: &Item{Type: TypeBool},
//		}
//	case reflect.String:
//		resp = &basicItemResponse{
//			item: &Item{Type: TypeString},
//		}
//	}
//	if resp != nil {
//		return resp
//	}
//	var name string
//	name, resp = y.buildRef(reflect.ValueOf(v))
//}
//
//func (y *YiDoc) buildRef(val reflect.Value) (string, Response) {
//	var (
//		refName string
//		resp    Response
//		items   []*Item
//	)
//	typ := val.Type()
//	refName = typ.Name()
//	switch typ.Kind() {
//	case reflect.Slice:
//		el := typ.Elem()
//		var t reflect.Type
//		switch el.Kind() {
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
//			resp = &basicItemResponse{
//				item: &Item{Type: TypeNumber},
//			}
//		case reflect.Bool:
//			resp = &basicItemResponse{
//				item: &Item{Type: TypeBool},
//			}
//		case reflect.String:
//			resp = &basicItemResponse{
//				item: &Item{Type: TypeString},
//			}
//		}
//		if resp != nil {
//			return "", resp
//		}
//		_, resp = y.buildRef(reflect.New(t).Elem())
//		if resp.IsObject() {
//			resp = &arrayResponse{item: resp.Object()[0]}
//		}
//	case reflect.Struct:
//		resp = &objectResponse{items: items}
//		for i := 0; i < typ.NumField(); i++ {
//			field := typ.Field(i)
//			fieldVal := val.Field(i)
//			if fieldVal.Kind() == reflect.Ptr {
//				fieldVal = reflect.New(field.Type.Elem()).Elem()
//			}
//			if !fieldVal.CanSet() {
//				continue
//			}
//			var it = new(Item)
//			if field.Anonymous {
//				switch field.Type.Kind() {
//				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
//					it.Type = TypeNumber
//				case reflect.Bool:
//					it.Type = TypeBool
//				case reflect.String:
//					it.Type = TypeString
//				default:
//					_, childItem := y.buildRef(fieldVal)
//					switch {
//					case childItem.IsArray():
//						it.Type = TypeArray
//						it.Item = childItem.ArrayElement()
//					case childItem.IsObject():
//						items = append(items, childItem.Object()...)
//						continue
//					default:
//						it.Item = childItem.Item()
//					}
//				}
//			}
//			it.Description = getTrimTag(field, "desc")
//			it.Required = getTrimTag(field, "required") != ""
//
//			queryTag := getTrimTag(field, "query")
//			fileTag := getTrimTag(field, "file")
//			formTag := getTrimTag(field, "form")
//			jsonTag := getTrimTag(field, "json")
//			switch {
//			case queryTag != "":
//				it.Name = queryTag
//			case fileTag != "":
//				it.Name = fileTag
//				it.Type = TypeFile
//			case formTag != "":
//				it.Name = formTag
//			case jsonTag != "":
//				it.Name = jsonTag
//			default:
//				it.Name = field.Name
//			}
//			if it.Type == "" {
//				it.Type = getItemTypeByKind(field.Type.Kind())
//			}
//			if (it.Type == TypeObject || it.Type == TypeArray) && !field.Anonymous {
//				_, childItem := y.buildRef(fieldVal)
//				switch {
//				case childItem.IsArray():
//					it.Type = TypeArray
//					it.Item = childItem.ArrayElement()
//				case childItem.IsObject():
//					items = append(items, childItem.Object()...)
//					continue
//				default:
//					it.Item = childItem.Item()
//				}
//			}
//		}
//	}
//
//	return refName, resp
//}
//
//func getTrimTag(field reflect.StructField, tagName string) string {
//	return strings.Split(strings.TrimSpace(field.Tag.Get(tagName)), " ")[0]
//}
//
//func getItemTypeByKind(kind reflect.Kind) string {
//	switch kind {
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
//		return TypeNumber
//	case reflect.Bool:
//		return TypeBool
//	case reflect.String:
//		return TypeString
//	case reflect.Slice:
//		return TypeArray
//	case reflect.Struct:
//		return TypeObject
//	}
//	return ""
//}
