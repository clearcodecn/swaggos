# Swaggos

swaggos is a tool for build swagger docs for golang. it generates docs from native golang code. 

## Installation
```
    go get -u github.com/clearcodecn/swaggos
```

## Usage

### New instance
basic config set host name and api prefix
``` 
    doc := swaggos.Default()
    doc.HostInfo("www.github.com","/api/v1")
    
    // default is application/json
    doc.Produces("application/json")
    
    // default is application/json
    doc.Consumes("application/json")
```

### Authorization
swagger support `oauth2`, `basic auth`, `apiKey` and this project is full implement.
#### Oauth2
```
    var scopes = []string{"openid"}
    var tokenURL = "https://yourtokenurl"
    var authURL = "https://yourAuthURL"

    // config password flow
    doc.Oauth2Password(tokenURL,scopes)
    // access code
    doc.Oauth2AccessCode(authURL,tokenURL,scopes)
    // client 
    doc.Oauth2Client(tokenURL,scopes)
    // implicit
    doc.Oauth2Implicit(authURL,scopes)
```
#### Basic Auth
```
    doc.Basic()
```
#### Custom token
```
    // will create header param
    // access_token: your token
    doc.JWT("access_token")
```

### Common Params 
```
    // will create header param in each request
    doc.Header("name","description",true)
    doc.Query("name","description",true)
    doc.Form("name","description",true)
```

### Paths
you can change every thing in path. 

> Tips: It's better to create struct for every params. 

```
    path := doc.Get("user_information")
    // now you can access path apis
    
    path.
        Tag("tagName"). // create a tag 
        Summary("summary"). // summary the request
        Description("...."). // create description
        ContentType("application/json","text/html"). // set content type

    // path params 
    path.Form("key",swaggos.Attribute{})
  
    // form files
    path.FormFile("file",swaggos.Attribute{Required:true})
    
    // form object reference
    path.FormObject(new(User))

    // query object
    path.QueryObject(new(User))
    
    // body
    path.Body(new(User))
    // json response
    path.JSON(new(User))

    // Attribute rules: 
    type Attribute struct {
    	Model       string      `json:"model"`          // key name
    	Description string      `json:"description"`    // description 
    	Required    bool        `json:"required"`       // if it's required 
    	Type        string      `json:"type"`           // the param type
    	Example     interface{} `json:"example"`        // example value
    
    	Nullable  bool          `json:"nullable,omitempty"`     // if it's nullable
    	Format    string        `json:"format,omitempty"`       // format 
    	Title     string        `json:"title,omitempty"`        // title 
    	Default   interface{}   `json:"default,omitempty"`      // default value
    	Maximum   *float64      `json:"maximum,omitempty"`       // max num
    	Minimum   *float64      `json:"minimum,omitempty"`       // min num
    	MaxLength *int64        `json:"maxLength,omitempty"`    // max length
    	MinLength *int64        `json:"minLength,omitempty"`    // min length
    	Pattern   string        `json:"pattern,omitempty"`      // regexp pattern
    	MaxItems  *int64        `json:"maxItems,omitempty"`     // max array length
    	MinItems  *int64        `json:"minItems,omitempty"`     // min array length
    	Enum      []interface{} `json:"enum,omitempty"`         // enum values
    	Ignore    bool          `json:"ignore"`                 // if it's ignore
    	Json      string        `json:"json"`                   // key name
    }


```

### Groups
group will add the common params to each path item below the group.
```
    	g := doc.Group("/api/v2")
    	g.Get("/user") // --> /api/v2/user
        // ... other methods
        g.Form ...
        g.Query ...
        g.Header ...
```

### path Response
```
    // will provide a example response
    // 400 
    path.BadRequest(map[string]interface{
            "data": nil,
            "code": 400,
    })
    // 401
    path.UnAuthorization(v)
    // 403
    path.Forbidden(v)
    // 500 
    path.ServerError(v)
```

### Global Response
```
    doc.Response(200, new(Success))
    doc.Response(400, new(Fail))
    doc.Response(500, new(ServerError))
```

###  Object Rules
 swaggos will parse object tag to create swagger rules.follow options are supported:

```
type RuleUser struct {
	// Model for field name
	// this example field name will be m1
	ModelName string `model:"m1" json:"m2"`
	// field name will be username
	Username string `json:"username"`
	//  field name will be Password
	Password string
	//  will be ignore
	Ignored string `json:"-"`
	//  true will be required field, false or empty will be not required
	Required string `required:"true"`
	// create description
	Description string `description:"this is description"`
	// a time type
	Type string `type:"time"`
	// default is abc
	DefaultValue string `default:"abc"`
	// max value is: 100
	Max float64 `maximum:"100"`
	// min value is: 0
	Min float64 `min:"0"`
	// MaxLength is 20
	MaxLength string `maxLength:"20"`
	// MinLength is 10
	MinLength string `minLength:"10"`
	// Pattern for regexp rules
	Pattern string `pattern:"\d{0,9}"`
	// array.length must <= 3
	MaxItems []int `maxItems:"3"`
	// array.length must >= 3
	MinItem []int `minItems:"3"`
	// Enum values
	EnumValue int `enum:"a,b,c,d"`
	// ignore
	IgnoreField string `ignore:"true"`
}
```
 

### Build
```
    data,_ := doc.Build()
    fmt.Println(string(data))
    => this is the swagger schema in json format

    data,_ := doc.Yaml()
    fmt.Println(string(data))
    => yaml format
```

### Serve HTTP
```
    http.Handle("/swagger",doc)
```


## Contact Me
QQ Group: 642154119 
![wechat](./images/wechat.png) 