# Swaggos

swaggos 是一个golang版本的swagger文档生成器，提供了native code包装器，并且支持主流的web框架包裹器

### 安装
```
    go get -u github.com/clearcodecn/swaggos
```

### 示例
> 目前只支持gin的包裹器
```go
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
示例将会生成该图例: [click here to see image](./images/ui.png)
您可以查看 examples 目录查看更多示例.

### 您也可以不使用包裹器
```go
func main() {
	doc := swaggos.Default()

	doc.HostInfo("localhost:8080", "/api").
		Response(200, newSuccessExample()).
		Response(400, newErrorExample())

	group := doc.Group("/users")
	group.Get("/list").JSON(CommonResponseWithData([]model.User{}))
	group.Post("/create").Body(new(model.User)).JSON(CommonResponseWithData(1))
	group.Put("/update").Body(new(model.User)).JSON(CommonResponseWithData(1))
	// path item
	group.Get("/{id}").JSON(new(model.User))
	group.Delete("/{id}").JSON(CommonResponseWithData(1))

	data, _ := doc.Build()
	fmt.Println(string(data))

	data, _ = doc.Yaml()
	fmt.Println(string(data))
}

```

### 使用

#### 增加请求头
```
     doc.Header("name","description",true)
    => generate a required header with key name
```

#### 增加jwt token
```
    doc.JWT("Authorization")
    => ui will create authorization in request headers.  
```

#### Oauth2 支持
```
    scopes:= []string{"openid"}
    doc.Oauth2("http://path/to/oauth/token/url",scopes,scopes)
    => ui will create a oauth2 password credentials client
```

### 增加 host 信息
```
    doc.HostInfo("yourhost.com","/your/api/prefix")
```

### 增加 响应 content-Type 类型
```
    doc.Produces("application/json")
```

### 增加 请求 content-Type 类型
```
    doc.Consumes("application/json")
```

### 生成json
```
    data,_ := doc.Build()
    fmt.Println(string(data))
    => this is the swagger schema in json format

    data,_ := doc.Yaml()
    fmt.Println(string(data))
    => yaml format
```

## struct的规则
   swaggos 会解析结构体的tag并将其赋值到 swagger 规则上面，下面是本项目支持的一些tag示例

```
    type User struct {
        // 字段名称  model > json
        // this example field name will be m1
        ModelName string `model:"m1" json:"m2"`
        // 字段名会是  username 
        Username string `json:"username"` 
        //  字段名会是 Password
        Password string 
        // 会被忽略
        Ignored `json:"-"`
        // 是否必须
        Required string `required:"true"`
        // 字段的描述
        Description string `description:"this is description"`
        // 字段的类型: string,integer,time,number,boolean,array...
        Type string `type:"time"`
        // 默认值 abc
        DefaultValue string `default:"abc"`
        // 最大值 100
        Max float64 `maximum:"100"`
        // 最小值 0
        Min float64 `min:"0"`
        // 最大长度 20
        MaxLength string `maxLength:"20"`
        // 最小长度 10
        MinLength string `minLength:"10"`
        // 正则表达式规则
        Pattern string `pattern:"\d{0,9}"`
        // 数组长度 小于3
        MaxItems []int `maxItems:"3"`
        // 数组长度大于3
        MinItem []int `minItems:"3"`
        // 枚举，用 , 分割
        EnumValue int `enum:"a,b,c,d"`
        // 忽略字段
        IgnoreField string `ignore:"true"`
        // 匿名字段规则：
        // 如果是一个基本类型，则直接添加, 
        // 如果是一个 数组，也将直接添加
        // 如果是一个结构体 但是带了json tag，将会作为一个字段
        // 如果是一个结构体 带没有json tag，将会将里面的子字段添加上该结构体上
        Anymouse
    }
```
 
 
## path上的工具方法
```
    path := doc.Get("/")
    // 创建一个 query 字段，包含了 描述和是否必须
    path.Query("name",DescRequired("description",true)).
    // 创建一个 query 字段，包含了 描述和是否必须 和默认值
    Query("name2",DescRequiredDefault("desc",true,"default"))
```

other useful functions:

``` 
    // 创建一个 swagger 的tag
    path.Tag("user group")
    
    // 请求的简单描述
    path.Summary("create a new user")

    // 请求的详细描述
    path.Description("....")
       
    // 设置请求-响应头
    path.ContentType("application/json","text/html")
   
    // form 字段
    path.Form("key1",swaggos.Attribute{Required:true})

    // 文件
    path.FormFile("file",swaggos.Attribute{Required:true})
    
    // form 用接头体解析
    path.FormObject(new(User))

    // query 用结构体解析
    path.QueryObject(new(User))
    
    // body 用结构体解析
    path.Body(new(User))

    // 响应json
    path.JSON(new(User))
```

## 响应
```
    // 响应带上具体的内容，将会创建具体的json示例
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


## 联系我
![wechat](./images/wechat.png) 