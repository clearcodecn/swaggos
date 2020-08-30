package yidoc

import (
	"bytes"
	"encoding/json"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	_id int64
)

func nextId() int64 {
	return atomic.AddInt64(&_id, 1)
}

type YiDoc struct {
	consumes            []string
	produces            []string
	schemes             []string
	swagger             string
	info                *spec.Info
	host                string
	basePath            string
	paths               *spec.Paths
	definitions         spec.Definitions
	parameters          map[string]spec.Parameter
	responses           map[string]spec.Response
	securityDefinitions spec.SecurityDefinitions
	security            []map[string][]string
	tags                []spec.Tag
	externalDocs        *spec.ExternalDocumentation

	o   sync.Once
	doc *spec.Swagger

	packageDef map[string]map[string]struct{} // package -> type
}

func NewDoc(opts ...Option) *YiDoc {
	doc := new(YiDoc)
	doc.packageDef = make(map[string]map[string]struct{})
	doc.definitions = make(map[string]spec.Schema)
	for _, o := range opts {
		o(doc)
	}
	return doc
}

func (y *YiDoc) Build() ([]byte, error) {
	y.buildOnce()
	return json.Marshal(y.doc)
}

func (y *YiDoc) BuildYaml() ([]byte, error) {
	y.buildOnce()
	var buf = bytes.NewBuffer(nil)
	err := yaml.NewEncoder(buf).Encode(y.doc)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (y *YiDoc) buildOnce() {
	y.o.Do(func() {
		y.build()
	})
}

func (y *YiDoc) build() {
	if y.doc == nil {
		return
	}
}

func (y *YiDoc) Model(v ...interface{}) {
	for _, m := range v {
		y.buildModel(m)
	}
}

func (y *YiDoc) buildModel(model interface{}) {
	typ := reflect.TypeOf(model)
	name := typ.Name()
	if !y.modelExist(typ) {
		y.packageDef[typ.PkgPath()][name] = struct{}{}
	}

	y.packageDef[typ.PkgPath()][name] = struct{}{}

	object := parseObject(model)
	schema := object.buildSchema(y)
}

func (y *YiDoc) newId(typ reflect.Type) string {

}

func (y *YiDoc) modelExist(typ reflect.Type) bool {
	if _, ok := y.packageDef[typ.PkgPath()]; !ok {
		y.packageDef[typ.PkgPath()] = make(map[string]struct{})
	}
	// already exist
	if _, ok := y.packageDef[typ.PkgPath()][typ.Name()]; ok {
		return true
	}
	return false
}

func (y *YiDoc) addModel(typ reflect.Type, props spec.SchemaProps) string {
	name := typ.Name()
	if !y.modelExist(typ) {
		y.packageDef[typ.PkgPath()][name] = struct{}{}
	}
	y.definitions[name] = spec.Schema{
		SchemaProps: props,
	}
	return name
}

/*
import (
	"github.com/go-openapi/spec"
)

type YiDoc struct {
	swagger *spec.Swagger

	paths map[string]map[string]*Path
}

func New() *YiDoc {
	return &YiDoc{
		swagger: &spec.Swagger{},
	}
}

type Path struct {
	op *spec.Operation
}

func (y *YiDoc) AddPath(method string, path string) *Path {
	if y.swagger.Paths == nil {
		y.swagger.Paths = &spec.Paths{
			Paths: map[string]spec.PathItem{},
		}
	}
	if _, ok := y.paths[path]; !ok {
		y.paths[path] = make(map[string]*Path)
	}
	ph := &Path{
		op: &spec.Operation{},
	}
	y.paths[path][method] = ph

	return ph
}

func (p *Path) Query(arg Arg, args ...Arg) *Path {
	//p.op.AddParam(&spec.Parameter{
	//	ParamProps: spec.ParamProps{
	//		Description:     "",
	//		Name:            "",
	//		In:              "",
	//		Required:        false,
	//		Schema:          nil,
	//		AllowEmptyValue: false,
	//	},
	//})

	return nil
}
*/
