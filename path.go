package swaggos

import (
	"github.com/go-openapi/spec"
	"strings"
)

const applicationJson = "application/json"

const (
	inForm = iota + 1
	inBody
)

const (
	InPath   = "path"
	InQuery  = "query"
	InBody   = "body"
	InHeader = "header"
	InForm   = "formData"
)

type Path struct {
	prop spec.OperationProps

	response  map[int]spec.Response
	paramDeep int

	doc *Swaggo
}

func newPath(d *Swaggo) *Path {
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
					Name:     name,
					In:       InPath,
					Required: true,
				},
			})
		}
	}
}

func (p *Path) Tag(v ...string) *Path {
	p.prop.Tags = v
	return p
}

func (p *Path) Summary(v string) *Path {
	p.prop.Summary = v
	return p
}

func (p *Path) Description(s string) *Path {
	p.prop.Description = s
	return p
}

func (p *Path) ContentType(req, resp string) {
	p.prop.Consumes = []string{req}
	p.prop.Produces = []string{resp}
}

func (p *Path) build() *spec.Operation {
	var (
		defaultResponse *spec.Response
	)
	if resp, ok := p.response[200]; ok {
		defaultResponse = &resp
	}
	p.prop.Responses = &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			Default:             defaultResponse,
			StatusCodeResponses: p.response,
		},
	}
	return &spec.Operation{
		OperationProps: p.prop,
	}
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
			Description: attribute.Description,
			Name:        name,
			In:          InForm,
			Required:    attribute.Required,
		},
	})
	return p
}

func (p *Path) FormObject(v interface{}) *Path {
	if p.paramDeep == inBody {
		panic("body and form can't be set at same time")
	}
	p.paramDeep = inForm
	ref := p.doc.buildSchema(v)
	for name, sch := range ref.SchemaProps.Properties {
		var param = spec.Parameter{
			SimpleSchema: spec.SimpleSchema{},
			ParamProps: spec.ParamProps{
				Name:   name,
				In:     InForm,
				Schema: &sch,
			},
		}
		p.prop.Parameters = append(p.prop.Parameters, param)
	}
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
			Description: attribute.Description,
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
			Description: attribute.Description,
			Name:        name,
			In:          InQuery,
			Required:    attribute.Required,
		},
	})
	return p
}

func (p *Path) QueryObject(v interface{}) *Path {
	ref := p.doc.buildSchema(v)
	for name, sch := range ref.SchemaProps.Properties {
		var param = spec.Parameter{
			SimpleSchema: spec.SimpleSchema{},
			ParamProps: spec.ParamProps{
				Name:   name,
				In:     InQuery,
				Schema: &sch,
			},
		}
		p.prop.Parameters = append(p.prop.Parameters, param)
	}
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
			Description: attribute.Description,
			Name:        name,
			In:          InHeader,
			Required:    attribute.Required,
		},
	})
	return p
}

func (p *Path) Body(v interface{}) *Path {
	if p.paramDeep == inForm {
		panic("body and form can't be set at same time")
	}
	p.paramDeep = inBody
	ref := p.doc.Define(v)
	p.prop.Parameters = append(p.prop.Parameters, spec.Parameter{
		ParamProps: spec.ParamProps{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
			In:   InBody,
			Name: InBody,
		},
	})
	return p
}

func (p *Path) JSON(v interface{}) *Path {
	ref := p.doc.Define(v)
	p.addResponse(200, ref, v)
	return p
}

func (p *Path) BadRequest(v interface{}) *Path {
	ref := p.doc.Define(v)
	p.addResponse(400, ref, v)
	return p
}

func (p *Path) ServerError(v interface{}) *Path {
	ref := p.doc.Define(v)
	p.addResponse(500, ref, v)
	return p
}

func (p *Path) Forbidden(v interface{}) *Path {
	ref := p.doc.Define(v)
	p.addResponse(403, ref, v)
	return p
}

func (p *Path) UnAuthorization(v interface{}) *Path {
	ref := p.doc.Define(v)
	p.addResponse(401, ref, v)
	return p
}

func (p *Path) addResponse(status int, ref spec.Ref, example interface{}) {
	p.response[status] = spec.Response{
		ResponseProps: spec.ResponseProps{
			Description: "json response",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
			Examples: map[string]interface{}{
				applicationJson: example,
			},
		},
	}
}
