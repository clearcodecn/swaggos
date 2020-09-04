package model

type User struct {
	Username     string `json:"username" required:"true"`
	Password     string `json:"password" required:"true" description:"密码" example:"123456" maxLength:"20" minLength:"6" pattern:"[a-zA-Z0-9]{6,20}"`
	Sex          int    `json:"sex" required:"false" default:"1" example:"1" format:"int64"`
	HeadImageURL string `json:"headImageUrl"`

	History string `json:"-"` // ignore
}

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
