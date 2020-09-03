# Swaggos

swaggos is a tool for build swagger docs for golang. it generates docs from native golang code. And it wraps popular web frameworks easy to use!

### Installation
```
    go get -u github.com/clearcodecn/swaggos
```

### Example
> for now it only support gin wrappers. and will support more web framework. 
```
package main

import (
	"github.com/clearcodecn/swaggos"
	"github.com/clearcodecn/swaggos/ginwrapper"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
)

type User struct {
	Username     string `json:"username" required:"true"`
	Password     string `json:"password" required:"true" description:"密码" example:"123456" maxLength:"20" minLength:"6" pattern:"[a-zA-Z0-9]{6,20}"`
	Sex          int    `json:"sex" required:"false" default:"1" example:"1" format:"int64"`
	HeadImageURL string `json:"headImageUrl"`

	History string `json:"-"` // ignore
}

func main() {
	g := ginwrapper.Default()
	doc := g.Doc()
	g.Gin().Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})
	doc.JWT("Authorization")
	doc.HostInfo("https://localhost:8080/", "/api/v1")
	group := g.Group("/api/v1")
	{
		group.GET("/users", listUsers).
			Query("order", swaggos.DescRequired("排序", false)).
			Query("q", swaggos.DescRequired("名称迷糊查询", false)).
			JSON([]User{})

		group.POST("/user/create", createUser).
			Body(new(User)).JSON(gin.H{"id": 1})

		group.DELETE("/user/*id", deleteUser).
			JSON(gin.H{"id": 1})

		group.PUT("/user/update", createUser).
			Body(new(User)).JSON(new(User))
	}
	g.ServeDoc()
	g.Gin().Run(":8888")
}

func listUsers(ctx *gin.Context)  {}
func createUser(ctx *gin.Context) {}
func deleteUser(ctx *gin.Context) {}

```
example will generate ui: [click here to see image](./images/ui.png)

### Usage

#### add header
```
     doc.Header("name","description",true)
    => generate a required header with key name
```

#### add jwt 
```
    doc.JWT("Authorization")
    => ui will create authorization in request headers.  
```

#### Oauth2
```
    scopes:= []string{"openid"}
    doc.Oauth2("http://path/to/oauth/token/url",scopes,scopes)
    => ui will create a oauth2 password credentials client
```

### add HostInfo
```
    doc.HostInfo("yourhost.com","/your/api/prefix")
```

### add Produces
```
    doc.Produces("application/json")
```

### add Consumes
```
    doc.Consumes("application/json")
```

### Build to json
```
    data,_ := doc.Build()
    fmt.Println(string(data))
    => this is the swagger schema in json format
```

## Object Rules
 swaggos will parse object tag to create swagger rules.follow options are supported:

```
    type User struct {
        // Model for field name
        // this example field name will be m1
        ModelName string `model:"m1" json:"m2"`
        // field name will be username
        Username string `json:"username"` 
        //  field name will be Password
        Password string 
        //  will be ignore
        Ignored `json:"-"`
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
        MaxItems []int `MaxItems:"3"`
        // array.length must >= 3
        MinItem []int `MinItems:"3"`
        // Enum values
        EnumValue int `enum:"a,b,c,d"`
        // ignore
        IgnoreField string `ignore:"true"`
    }
```
 
 
## Utils functions in path item
```
    path := doc.Get("/")
    // create a query field with description and required
    path.Query("name",DescRequired("description",true)).
        // create a field with description and required and default value
        Query("name2",DescRequiredDefault("desc",true,"default"))
```

other useful functions:

``` 
    // create a groups for user like '/users' prefix
    path.Tag("user group")
    
    // simple description for a api
    path.Summary("create a new user")

    // description for api
    path.Description("....")
       
    // set content-type
    path.ContentType("application/json","text/html")
   
    // form values 
    path.Form("key1",swaggos.Attribute{Required:true})

    // form files
    path.FormFile("file",swaggos.Attribute{Required:true})
    
    // form object reference
    path.FormObject(new(User))

    // query object
    path.QueryObject(new(User))
    
    // body
    path.Body(new(User))
```

## create response
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