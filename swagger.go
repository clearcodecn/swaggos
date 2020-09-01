package yidoc

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/spec"
	"net/http"
	"strings"
)

func NewYiDoc() *YiDoc {
	doc := new(YiDoc)
	doc.definitions = make(spec.Definitions)
	doc.paths = make(map[string]map[string]*Path)
	return doc
}

func (y *YiDoc) Get(path string) *Path     { return y.addPath(http.MethodGet, path) }
func (y *YiDoc) Post(path string) *Path    { return y.addPath(http.MethodPost, path) }
func (y *YiDoc) Put(path string) *Path     { return y.addPath(http.MethodPut, path) }
func (y *YiDoc) Patch(path string) *Path   { return y.addPath(http.MethodPatch, path) }
func (y *YiDoc) Options(path string) *Path { return y.addPath(http.MethodOptions, path) }
func (y *YiDoc) Delete(path string) *Path  { return y.addPath(http.MethodDelete, path) }

func (y *YiDoc) trimPath(path string) string {
	if len(path) == 0 {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return strings.TrimPrefix(path, y.basePath)
}

func (y *YiDoc) addPath(method string, path string) *Path {
	path = y.trimPath(path)
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

func (y *YiDoc) JWT(keyName string) *YiDoc {
	def := spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Description: "jwt token",
			Type:        "apiKey",
			Name:        keyName,
			In:          "header",
		},
	}
	y.addAuth(keyName, &def, nil)
	return y
}

func (y *YiDoc) Oauth2(tokenURL string, scopes []string, permits []string) *YiDoc {
	oauth2 := spec.OAuth2Password(tokenURL)
	if len(scopes) == 0 {
		scopes = []string{"openid"}
	}
	for _, scope := range scopes {
		oauth2.AddScope(scope, "")
	}
	y.addAuth("Oauth2", oauth2, map[string][]string{
		"Oauth2": permits,
	})
	return y
}

func (y *YiDoc) Header(name string, desc string, required bool) {
	if y.params == nil {
		y.params = make(map[string]spec.Parameter)
	}
	if _, ok := y.params[name]; ok {
		panic(fmt.Errorf("repeated header param: %s", name))
	}
	param := spec.Parameter{
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          InHeader,
			Required:    required,
		},
	}
	y.params[name] = param
}

func (y *YiDoc) addAuth(key string, schema *spec.SecurityScheme, security map[string][]string) {
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
}

func (y *YiDoc) HostInfo(host string, basePath string, info spec.InfoProps) *YiDoc {
	y.info = info
	y.host = host
	y.basePath = basePath
	return y
}

func (y *YiDoc) Build() ([]byte, error) {
	swag := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Consumes:            []string{applicationJson},
			Produces:            []string{applicationJson},
			Swagger:             "2.0",
			Info:                &spec.Info{InfoProps: y.info},
			Host:                y.host,
			BasePath:            y.basePath,
			Paths:               y.buildPaths(),
			Definitions:         y.definitions,
			SecurityDefinitions: y.securityDefinitions,
			Security:            y.security,
			Parameters:          y.params,
		},
	}
	data, err := json.MarshalIndent(swag, "", " ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (y *YiDoc) buildPaths() *spec.Paths {
	paths := &spec.Paths{Paths: map[string]spec.PathItem{}}
	for path, items := range y.paths {
		for method, item := range items {
			var operate = item.build()
			pi := spec.PathItem{
				PathItemProps: spec.PathItemProps{},
			}
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
			paths.Paths[path] = pi
		}
	}
	return paths
}
