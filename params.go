package ydoc

type ParamPosition int

const (
	ParamInPath ParamPosition = iota + 1
	ParamInQuery
	ParamInHeader
	ParamInBody
	ParamInFormData
)

var paramPositions = map[ParamPosition]string{
	ParamInPath:     "path",
	ParamInQuery:    "query",
	ParamInHeader:   "header",
	ParamInBody:     "body",
	ParamInFormData: "formData",
}

func getPosition(pos ParamPosition) string {
	return paramPositions[pos]
}

// 参数限制：
// 1. 如果存在body参数，那么只能有一个参数
// 2. 如果存在body参数，那么不能有form(file)参数
type ParameterObject struct {
	Name        string `json:"name"` // 如果在path上，那么这个名字必须和path里面的相同， required
	In          string `json:"in"`   // required,
	Description string `json:"description"`
	Required    bool   `json:"required"`

	// 如果in == body. schema 必须存在.
	Schema string `json:"schema"`
	// 如果不是 body, 那么下面的字段可选
	Type            string `json:"type"` // "string", "number", "integer", "boolean", "array" or "file". 请求头为  "multipart/form-data", " application/x-www-form-urlencoded"
	Format          string `json:"format"`
	AllowEmptyValue bool   `json:"allowEmptyValue"` // 是否允许为空,默认为true

	Items interface{} `json:"items"` // 如果字段的类型是 item , 那么这个字段用于描述 item子元素的内容.
	// 一些其他的验证字段.  https://swagger.io/specification/v2/#schemaObject
	Maximum int `json:"maximum"`
	// TODO..
}