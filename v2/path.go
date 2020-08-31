package v2

import (
	"github.com/go-openapi/spec"
	"reflect"
	"strings"
)

const applicationJson = "application/json"

const (
	inPath = iota + 1
	inQuery
	inForm
	inBody
	inHeader
)

type AttributeType string

const (
	String  = "string"
	Number  = "number"
	Integer = "integer"
	Boolean = "boolean"
	Array   = "array"
	File    = "file"
)

const (
	InPath   = "path"
	InQuery  = "query"
	InBody   = "body"
	InHeader = "header"
	InForm   = "formData"
)

type Format string

const (
	Int32    Format = "int32"
	Int64           = "int64"
	Float           = "float"
	Double          = "double"
	Byte            = "byte"
	Binary          = "binary"
	Date            = "date"
	DateTime        = "date-time"
	Password        = "password"
)

type Path struct {
	prop spec.OperationProps

	response  map[int]spec.Response
	paramDeep int

	doc *YiDoc
}

func newPath(d *YiDoc) *Path {
	path := new(Path)
	path.response = make(map[int]spec.Response)
	path.paramDeep = 0
	path.doc = d
	return path
}

func (p *Path) parsePath(path string) {
	arr := strings.Split(path, "/")
	for _, a := range arr {
		if a == "" {
			continue
		}
		if a[0] == '{' && a[len(a)-1] == '}' {
			name := a[1 : len(a)-1]
			if name == "" {
				continue
			}
			p.prop.Parameters = append(p.prop.Parameters, spec.Parameter{
				SimpleSchema: spec.SimpleSchema{
					Type: String,
				},
				ParamProps: spec.ParamProps{
					Name: name,
					In:   InPath,
				},
			})
		}
	}
}

type Attribute struct {
	Desc     string
	Required bool
	Type     AttributeType
	Format   Format
	Default  interface{}
	Example  interface{}
}

func (p *Path) Form(name string, attribute Attribute) *Path {
	if p.paramDeep == inBody {
		panic("body and form can't be set at same time")
	}
	p.paramDeep = inForm
	p.prop.Parameters = append(p.prop.Parameters, spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type:    string(attribute.Type),
			Format:  string(attribute.Format),
			Default: attribute.Default,
			Example: attribute.Example,
		},
		ParamProps: spec.ParamProps{
			Description: attribute.Desc,
			Name:        name,
			In:          InForm,
			Required:    attribute.Required,
		},
	})
	return p
}

func (p *Path) FormFile(name string, attribute Attribute) *Path {
	if p.paramDeep == inBody {
		panic("body and form can't be set at same time")
	}
	p.paramDeep = inForm
	p.prop.Parameters = append(p.prop.Parameters, spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: File,
		},
		ParamProps: spec.ParamProps{
			Description: attribute.Desc,
			Name:        name,
			In:          InForm,
			Required:    attribute.Required,
		},
	})
	return p
}

func (p *Path) Query(name string, attribute Attribute) *Path {
	p.prop.Parameters = append(p.prop.Parameters, spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type:    string(attribute.Type),
			Format:  string(attribute.Format),
			Default: attribute.Default,
			Example: attribute.Example,
		},
		ParamProps: spec.ParamProps{
			Description: attribute.Desc,
			Name:        name,
			In:          InQuery,
			Required:    attribute.Required,
		},
	})
	return p
}

func (p *Path) Header(name string, attribute Attribute) *Path {
	p.prop.Parameters = append(p.prop.Parameters, spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type:    string(attribute.Type),
			Format:  string(attribute.Format),
			Default: attribute.Default,
			Example: attribute.Example,
		},
		ParamProps: spec.ParamProps{
			Description: attribute.Desc,
			Name:        name,
			In:          InHeader,
			Required:    attribute.Required,
		},
	})
	return p
}

func (p *Path) Body(v interface{}, names ...string) {
	if p.paramDeep == inForm {
		panic("body and form can't be set at same time")
	}
	p.paramDeep = inBody
	refName := reflect.TypeOf(v).Name()
	if len(names) > 0 {
		if len(names) > 0 {
			refName = names[0]
		}
	}
	ref := p.doc.Define(refName, v)
	p.prop.Parameters = append(p.prop.Parameters, spec.Parameter{
		ParamProps: spec.ParamProps{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
		},
	})
}

func (p *Path) JSON(v interface{}, names ...string) {
	refName := reflect.TypeOf(v).Name()
	if len(names) > 0 {
		if len(names) > 0 {
			refName = names[0]
		}
	}
	ref := p.doc.Define(refName, v)
	resp := spec.Response{
		ResponseProps: spec.ResponseProps{
			Description: "json response",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
			Examples: map[string]interface{}{
				applicationJson: v,
			},
		},
	}
	p.response[200] = resp
	//p.prop.Responses = &spec.Responses{
	//	ResponsesProps: spec.ResponsesProps{
	//		Default:             &resp,
	//		StatusCodeResponses: p.response,
	//	},
	//}
}

func (p *Path) BadRequest(v interface{}, names ...string) {
	refName := reflect.TypeOf(v).Name()
	if len(names) > 0 {
		if len(names) > 0 {
			refName = names[0]
		}
	}
	ref := p.doc.Define(refName, v)
	p.response[400] = spec.Response{
		ResponseProps: spec.ResponseProps{
			Description: "json response",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
			Examples: map[string]interface{}{
				applicationJson: v,
			},
		},
	}
}

func (p *Path) ServerError(v interface{}, names ...string) {
	refName := reflect.TypeOf(v).Name()
	if len(names) > 0 {
		if len(names) > 0 {
			refName = names[0]
		}
	}
	ref := p.doc.Define(refName, v)
	p.response[500] = spec.Response{
		ResponseProps: spec.ResponseProps{
			Description: "json response",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
			Examples: map[string]interface{}{
				applicationJson: v,
			},
		},
	}
}

func (p *Path) Forbidden(v interface{}, names ...string) {
	refName := reflect.TypeOf(v).Name()
	if len(names) > 0 {
		if len(names) > 0 {
			refName = names[0]
		}
	}
	ref := p.doc.Define(refName, v)
	p.response[403] = spec.Response{
		ResponseProps: spec.ResponseProps{
			Description: "json response",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
			Examples: map[string]interface{}{
				applicationJson: v,
			},
		},
	}
}

func (p *Path) Tag(v ...string) {
	p.prop.Tags = v
}

func (p *Path) Summary(v string) {
	p.prop.Summary = v
}

func (p *Path) ContentType(req, resp string) {
	p.prop.Consumes = []string{req}
	p.prop.Produces = []string{resp}
}
