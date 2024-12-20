# Swaggos

swaggos 是一个golang版本的swagger文档生成器，提供了native code包装器. 

## 安装
```
    go get -u github.com/clearcodecn/swaggos
```

## 使用

### 创建实例
> 创建一个新的实例，配置一些基本信息 `host` 和 `apiPrefix` 
``` 
    doc := swaggos.Default()
    doc.HostInfo("www.github.com","/api/v1")
    
    // default is application/json
    doc.Produces("application/json")
    
    // default is application/json
    doc.Consumes("application/json")
```

### Authorization
项目支持swagger的所有鉴权方式：`oauth2`, `basic auth`, `apiKey` 
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
#### 自定义Token
> 会在header中增加`access_token`参数
```
    doc.JWT("access_token")
```

### 公共参数
```
    // will create header param in each request
    doc.Header("name","description",true)
    doc.Query("name","description",true)
    doc.Form("name","description",true)
```

### 请求路劲
> path 是一个请求的实例，支持流式写法. 
```
    // 创建一个 path
    path := doc.Get("user_information")
    
    path.
        Tag("tagName"). // 创建 tag 
        Summary("summary"). // 总结
        Description("...."). // 描述
        ContentType("application/json","text/html"). // 请求/响应类型

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

### 路劲响应
```
    // 响应json，创建model
    path.JSON(new(Response))
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

### 分组
> 分组将对api进行分组，组下面的所有路劲会共享分组的公共参数
```
    	g := doc.Group("/api/v2")
    	g.Get("/user") // --> /api/v2/user
        // ... other methods
        g.Form ...
        g.Query ...
        g.Header ...
```

### 全局响应
```
    doc.Response(200, new(Success))
    doc.Response(400, new(Fail))
    doc.Response(500, new(ServerError))
```

### 结构体的tag支持
 

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
        Anonymous
    }
```
 

### 构建json和yaml
```
    data,_ := doc.Build()
    fmt.Println(string(data))
    => this is the swagger schema in json format

    data,_ := doc.Yaml()
    fmt.Println(string(data))
    => yaml format
```

### gin 集成 swagger ui 
```
    // 打开监听地址： http://localhost:端口/doc 访问. 
    ginutils.Serve("/doc",doc,ginEngine,nil)
```

> Tips: examples 目录下面有少量的示例


## Contact Me
QQ群: 642154119 
![wechat](./images/wechat.png) 