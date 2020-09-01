package yidoc

import (
	"reflect"
	"strconv"
	"strings"
)

type Attribute struct {
	Model       string        `json:"model"`
	Description string        `json:"description"`
	Required    bool          `json:"required"`
	Type        AttributeType `json:"type"`
	Example     interface{}   `json:"example"`

	Nullable  bool          `json:"nullable,omitempty"`
	Format    string        `json:"format,omitempty"`
	Title     string        `json:"title,omitempty"`
	Default   interface{}   `json:"default,omitempty"`
	Maximum   *float64      `json:"maximum,omitempty"`
	Minimum   *float64      `json:"minimum,omitempty"`
	MaxLength *int64        `json:"maxLength,omitempty"`
	MinLength *int64        `json:"minLength,omitempty"`
	Pattern   string        `json:"pattern,omitempty"`
	MaxItems  *int64        `json:"maxItems,omitempty"`
	MinItems  *int64        `json:"minItems,omitempty"`
	Enum      []interface{} `json:"enum,omitempty"`
	Ignore    bool          `json:"ignore"`
	Json      string        `json:"json"`
}

func (a *Attribute) parseTag(t reflect.StructTag) {
	a.Description = t.Get("description")
	// required
	a.Required = t.Get("nullable") == "true"
	example := t.Get("example")
	if example != "" {
		a.Example = example
	}
	a.Nullable = t.Get("nullable") == "true"
	a.Format = t.Get("format")
	a.Title = t.Get("title")
	a.Default = t.Get("default")
	a.Maximum = str2f64Ptr(t.Get("maximum"))
	a.Minimum = str2f64Ptr(t.Get("minimum"))
	a.MaxLength = str2i64Ptr(t.Get("maxLength"))
	a.MinLength = str2i64Ptr(t.Get("minLength"))
	a.Pattern = t.Get("pattern")
	a.MaxItems = str2i64Ptr(t.Get("maxItems"))
	a.MinItems = str2i64Ptr(t.Get("minItems"))
	var enum []interface{}
	arr := strings.Split(t.Get("enum"), ",")
	for _, a := range arr {
		if a != "" {
			enum = append(enum, a)
		}
	}
	a.Enum = enum
	a.Model = t.Get("model")
	j := t.Get("json")
	if j == "-" {
		a.Json = "-"
	} else {
		a.Json = strings.Split(j, ",")[0]
	}
	if a.Json == "-" || a.Model == "-" {
		a.Ignore = true
	}
}

func str2f64Ptr(s string) *float64 {
	if s == "" {
		return nil
	}
	f, _ := strconv.ParseFloat(s, 64)
	return &f
}

func str2i64Ptr(s string) *int64 {
	if s == "" {
		return nil
	}
	f, _ := strconv.ParseInt(s, 10, 64)
	return &f
}
