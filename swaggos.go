package swaggos

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-openapi/spec"
	"net/http"
	"reflect"
	"strings"
)

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
func (y *Swaggos) Get(path string) *Path { return y.addPath(http.MethodGet, path) }

// Post add a post path operation.
func (y *Swaggos) Post(path string) *Path { return y.addPath(http.MethodPost, path) }

// Put add a put operation
func (y *Swaggos) Put(path string) *Path { return y.addPath(http.MethodPut, path) }

// Patch add a patch operation.
func (y *Swaggos) Patch(path string) *Path { return y.addPath(http.MethodPatch, path) }

// Options add a options operation
func (y *Swaggos) Options(path string) *Path { return y.addPath(http.MethodOptions, path) }

// Delete add a delete operation
func (y *Swaggos) Delete(path string) *Path { return y.addPath(http.MethodDelete, path) }

func (y *Swaggos) addPath(method string, path string) *Path {
	path = "/" + strings.Trim(path, "/")
	path = "/" + strings.TrimLeft(strings.TrimPrefix(path, y.basePath), "/")
	if _, ok := y.paths[path]; !ok {
		y.paths[path] = make(map[string]*Path)
	}
	if _, ok := y.paths[path][method]; ok {
		panic(fmt.Errorf("repeated method&path: %s %s", method, path))
	}
	p := newPath(y)
	y.paths[path][method] = p
	p.parsePath(path)
	return p
}

// JWT create a jwt header
func (y *Swaggos) JWT(keyName string) *Swaggos {
	def := spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Description: "jwt token",
			Type:        "apiKey",
			Name:        keyName,
			In:          "header",
		},
	}
	y.addAuth(keyName, &def, map[string][]string{
		keyName: []string{},
	})
	return y
}

func (y *Swaggos) Basic() *Swaggos {
	def := spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Description: "basic auth",
			Type:        "basic",
		},
	}
	y.addAuth(`basicAuth`, &def, map[string][]string{
		`basicAuth`: []string{},
	})
	return y
}

// Header add a custom header
func (y *Swaggos) Header(name string, desc string, required bool) {
	if y.params == nil {
		y.params = make(map[string]spec.Parameter)
	}
	if _, ok := y.params[name]; ok {
		panic(fmt.Errorf("repeated header param: %s", name))
	}
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          InHeader,
			Required:    required,
		},
	}
	y.params[name] = param
}

func (y *Swaggos) addAuth(key string, schema *spec.SecurityScheme, security map[string][]string) *Swaggos {
	if y.securityDefinitions == nil {
		y.securityDefinitions = make(map[string]*spec.SecurityScheme)
	}
	if y.security == nil {
		y.security = make([]map[string][]string, 0)
	}
	y.securityDefinitions[key] = schema
	if security == nil {
		security = make(map[string][]string)
	}
	y.security = append(y.security, security)
	return y
}

// HostInfo add host info to documents
func (y *Swaggos) HostInfo(host string, basePath string) *Swaggos {
	y.info = spec.InfoProps{
		Version: "2.0",
		Title:   fmt.Sprintf("document of %s", host),
	}
	y.host = strings.TrimRight(host, "/")
	y.basePath = "/" + strings.Trim(basePath, "/")
	return y
}

// Produces create global produces header
func (y *Swaggos) Produces(s ...string) {
	y.produces = append(y.produces, s...)
}

// Consumes create global consumes header
func (y *Swaggos) Consumes(s ...string) {
	y.consumes = append(y.consumes, s...)
}

// Build return json schema of swagger doc
func (y *Swaggos) Build() ([]byte, error) {
	swag := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Consumes:            y.consumes,
			Produces:            y.produces,
			Swagger:             "2.0",
			Info:                &spec.Info{InfoProps: y.info},
			Host:                y.host,
			BasePath:            y.basePath,
			Paths:               y.buildPaths(),
			Definitions:         y.definitions,
			SecurityDefinitions: y.securityDefinitions,
			Security:            y.security,
			Parameters:          y.params,
			Schemes:             y.schemas,
			ExternalDocs:        y.extend,
		},
	}
	data, err := json.Marshal(swag)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (y *Swaggos) Yaml() ([]byte, error) {
	data, err := y.Build()
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

func (y *Swaggos) buildPaths() *spec.Paths {
	paths := &spec.Paths{Paths: map[string]spec.PathItem{}}
	for path, items := range y.paths {
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
			for name := range y.params {
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
func (y *Swaggos) Extend(url string, desc string) {
	y.extend = &spec.ExternalDocumentation{
		Description: desc,
		URL:         url,
	}
}

// Query add query param to the group
func (y *Swaggos) Query(name string, desc string, required bool) *Swaggos {
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          InQuery,
			Required:    required,
		},
	}
	y.params[name] = param
	return y
}

// Form add form param to the group
func (y *Swaggos) Form(name string, desc string, required bool) *Swaggos {
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          InForm,
			Required:    required,
		},
	}
	y.params[name] = param
	return y
}

func (y *Swaggos) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data []byte
	format := r.URL.Query().Get("format")
	var contentType = "application/json"
	switch format {
	case "yaml":
		data, _ = y.Yaml()
		contentType = "application/yaml"
	default:
		data, _ = y.Build()
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
