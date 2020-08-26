package yidoc

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
	p.op.AddParam(&spec.Parameter{
		ParamProps: spec.ParamProps{
			Description:     "",
			Name:            "",
			In:              "",
			Required:        false,
			Schema:          nil,
			AllowEmptyValue: false,
		},
	})
}
