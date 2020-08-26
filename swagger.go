package ydoc

type Swagger struct {
	Swagger             string                           `json:"swagger"`
	Info                Info                             `json:"info"`
	Host                string                           `json:"host"`
	BasePath            string                           `json:"basePath"`
	Schemes             []string                         `json:"schemes"` // "http", "https", "ws", "wss".
	Consumes            []string                         `json:"consumes"`
	Produces            []string                         `json:"produces"`
	Paths               map[string]map[string]PathObject `json:"paths"`
	Definitions         interface{}                      `json:"definitions"`
	Parameters          interface{}                      `json:"parameters"`
	Response            interface{}                      `json:"response"`
	SecurityDefinitions interface{}                      `json:"securityDefinitions"`
	Security            interface{}                      `json:"security"`
	Tags                []string                         `json:"tags"`
	ExteernalDocs       interface{}                      `json:"exteernalDocs"`
}

type Info struct {
	Title          string  `json:"title"` // required
	Description    string  `json:"description"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"` // api version
}

type Contact struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PathObject struct {
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	//ExternalDocs []string `json:"externalDocs"`
	OperationId string      `json:"operationId"` // must unique
	Consumes    []string    `json:"consumes"`
	Produces    []string    `json:"produces"`
	Parameters  Parameter   `json:"parameters"`
	Response    Response    `json:"response"`
	Schemes     []string    `json:"schemes"`
	Deprecated  bool        `json:"deprecated"`
	Security    interface{} `json:"security"`
}

type Parameter interface{}

type Response struct {
}
