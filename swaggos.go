package swaggos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/go-openapi/spec"
)

// Swaggos is the base builder
type Swaggos struct {
	definitions spec.Definitions
	paths       map[string]map[string]*Path

	securityDefinitions spec.SecurityDefinitions
	security            []map[string][]string
	info                spec.InfoProps
	host                string
	basePath            string
	params              map[string]spec.Parameter
	schemas             []string
	extend              *spec.ExternalDocumentation

	typeNames map[reflect.Type]string
	produces  []string
	consumes  []string
	response  map[int]spec.Response
}

// NewSwaggo returns a new Swaggos instanence
func NewSwaggo(option ...Option) *Swaggos {
	doc := new(Swaggos)
	doc.definitions = make(spec.Definitions)
	doc.paths = make(map[string]map[string]*Path)
	doc.params = make(map[string]spec.Parameter)
	for _, o := range option {
		o(doc)
	}
	return doc
}

// Default create a default swaggo instanence.
func Default() *Swaggos {
	return NewSwaggo(DefaultOptions()...)
}

// Get add a get path operation.
func (swaggos *Swaggos) Get(path string) *Path { return swaggos.addPath(http.MethodGet, path) }

// Post add a post path operation.
func (swaggos *Swaggos) Post(path string) *Path { return swaggos.addPath(http.MethodPost, path) }

// Put add a put operation
func (swaggos *Swaggos) Put(path string) *Path { return swaggos.addPath(http.MethodPut, path) }

// Patch add a patch operation.
func (swaggos *Swaggos) Patch(path string) *Path { return swaggos.addPath(http.MethodPatch, path) }

// Options add a options operation
func (swaggos *Swaggos) Options(path string) *Path { return swaggos.addPath(http.MethodOptions, path) }

// Delete add a delete operation
func (swaggos *Swaggos) Delete(path string) *Path { return swaggos.addPath(http.MethodDelete, path) }

func (swaggos *Swaggos) addPath(method string, path string) *Path {
	path = "/" + strings.Trim(path, "/")
	path = "/" + strings.TrimLeft(strings.TrimPrefix(path, swaggos.basePath), "/")
	if _, ok := swaggos.paths[path]; !ok {
		swaggos.paths[path] = make(map[string]*Path)
	}
	if _, ok := swaggos.paths[path][method]; ok {
		panic(fmt.Errorf("repeated method&path: %s %s", method, path))
	}
	p := newPath(swaggos)
	swaggos.paths[path][method] = p
	p.parsePath(path)
	return p
}

// JWT create a jwt header
func (swaggos *Swaggos) JWT(keyName string) *Swaggos {
	def := spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Description: "jwt token",
			Type:        "apiKey",
			Name:        keyName,
			In:          "header",
		},
	}
	swaggos.addAuth(keyName, &def, map[string][]string{
		keyName: []string{},
	})
	return swaggos
}

// BasicAuth set basic auth support
func (swaggos *Swaggos) BasicAuth() *Swaggos {
	def := spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Description: "basic auth",
			Type:        "basic",
		},
	}
	swaggos.addAuth(`basicAuth`, &def, map[string][]string{
		`basicAuth`: []string{},
	})
	return swaggos
}

// Header add a custom header
func (swaggos *Swaggos) Header(name string, desc string, required bool) {
	if swaggos.params == nil {
		swaggos.params = make(map[string]spec.Parameter)
	}
	if _, ok := swaggos.params[name]; ok {
		panic(fmt.Errorf("repeated header param: %s", name))
	}
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: _String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          _InHeader,
			Required:    required,
		},
	}
	swaggos.params[name] = param
}

func (swaggos *Swaggos) addAuth(key string, schema *spec.SecurityScheme, security map[string][]string) *Swaggos {
	if swaggos.securityDefinitions == nil {
		swaggos.securityDefinitions = make(map[string]*spec.SecurityScheme)
	}
	if swaggos.security == nil {
		swaggos.security = make([]map[string][]string, 0)
	}
	swaggos.securityDefinitions[key] = schema
	if security == nil {
		security = make(map[string][]string)
	}
	swaggos.security = append(swaggos.security, security)
	return swaggos
}

// HostInfo add host info to documents
func (swaggos *Swaggos) HostInfo(host string, basePath string) *Swaggos {
	swaggos.info = spec.InfoProps{
		Version: "2.0",
		Title:   fmt.Sprintf("document of %s", host),
	}
	swaggos.host = strings.TrimRight(host, "/")
	swaggos.basePath = "/" + strings.Trim(basePath, "/")
	return swaggos
}

// Produces create global produces header
func (swaggos *Swaggos) Produces(s ...string) {
	swaggos.produces = append(swaggos.produces, s...)
}

// Consumes create global consumes header
func (swaggos *Swaggos) Consumes(s ...string) {
	swaggos.consumes = append(swaggos.consumes, s...)
}

// Build return json data of swagger doc
func (swaggos *Swaggos) Build() ([]byte, error) {
	swag := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Consumes:            swaggos.consumes,
			Produces:            swaggos.produces,
			Swagger:             "2.0",
			Info:                &spec.Info{InfoProps: swaggos.info},
			Host:                swaggos.host,
			BasePath:            swaggos.basePath,
			Paths:               swaggos.buildPaths(),
			Definitions:         swaggos.definitions,
			SecurityDefinitions: swaggos.securityDefinitions,
			Security:            swaggos.security,
			Parameters:          swaggos.params,
			Schemes:             swaggos.schemas,
			ExternalDocs:        swaggos.extend,
		},
	}
	data, err := json.Marshal(swag)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Yaml build yaml data of swagger doc
func (swaggos *Swaggos) Yaml() ([]byte, error) {
	data, err := swaggos.Build()
	if err != nil {
		return nil, err
	}
	var i interface{}
	err = yaml.Unmarshal(data, &i)
	if err != nil {
		return nil, err
	}
	data, err = yaml.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (swaggos *Swaggos) buildPaths() *spec.Paths {
	paths := &spec.Paths{Paths: map[string]spec.PathItem{}}
	for path, items := range swaggos.paths {
		pi := spec.PathItem{
			PathItemProps: spec.PathItemProps{},
		}
		for method, item := range items {
			var operate = item.build()
			switch method {
			case http.MethodGet:
				pi.Get = operate
			case http.MethodPut:
				pi.Put = operate
			case http.MethodPost:
				pi.Post = operate
			case http.MethodPatch:
				pi.Patch = operate
			case http.MethodDelete:
				pi.Delete = operate
			case http.MethodOptions:
				pi.Options = operate
			}
			for name := range swaggos.params {
				pi.PathItemProps.Parameters = append(pi.PathItemProps.Parameters, spec.Parameter{
					Refable: spec.Refable{
						Ref: spec.MustCreateRef("#/parameters/" + name),
					},
				})
			}
		}
		paths.Paths[path] = pi
	}
	return paths
}

// Extend extend the swagger docs.
func (swaggos *Swaggos) Extend(url string, desc string) {
	swaggos.extend = &spec.ExternalDocumentation{
		Description: desc,
		URL:         url,
	}
}

// Query add query param to the group
func (swaggos *Swaggos) Query(name string, desc string, required bool) *Swaggos {
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: _String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          _InQuery,
			Required:    required,
		},
	}
	swaggos.params[name] = param
	return swaggos
}

// Form add form param to the group
func (swaggos *Swaggos) Form(name string, desc string, required bool) *Swaggos {
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: _String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          _InForm,
			Required:    required,
		},
	}
	swaggos.params[name] = param
	return swaggos
}

// ServeHTTP export a http handler for serve document.
func (swaggos *Swaggos) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data []byte
	format := r.URL.Query().Get("format")
	var contentType = "application/json"
	switch format {
	case "yaml":
		data, _ = swaggos.Yaml()
		contentType = "application/yaml"
	default:
		data, _ = swaggos.Build()
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

// Response setup response object and build a example
func (swaggos *Swaggos) Response(status int, v interface{}) *Swaggos {
	ref := swaggos.Define(v)
	if swaggos.response == nil {
		swaggos.response = make(map[int]spec.Response)
	}
	swaggos.response[status] = spec.Response{
		ResponseProps: spec.ResponseProps{
			Description: "json response",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: ref,
				},
			},
			Examples: map[string]interface{}{
				applicationJSON: v,
			},
		},
	}
	return swaggos
}
