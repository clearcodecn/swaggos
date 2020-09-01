package yidoc

import (
	"github.com/go-openapi/spec"
	"reflect"
)

type documentTag struct {
	tag       reflect.StructTag
	attribute *Attribute
}

func newTags(tag reflect.StructTag) *documentTag {
	dt := new(documentTag)
	dt.tag = tag
	a := &Attribute{}
	a.parseTag(tag)
	dt.attribute = a
	return dt
}

func (t *documentTag) name() string {
	if t.attribute.Model != "" {
		return t.attribute.Model
	}
	if t.attribute.Json != "" {
		return t.attribute.Json
	}
	return ""
}

func (t *documentTag) ignore() bool {
	return t.attribute.Ignore
}

func (t *documentTag) required() bool {
	return t.attribute.Required
}

func (t *documentTag) Attribute() *Attribute {
	return t.attribute
}

func (t *documentTag) mergeSchema(schema spec.Schema) spec.Schema {
	schema.Description = t.attribute.Description
	schema.Example = t.attribute.Example
	schema.Nullable = t.attribute.Nullable
	schema.Format = t.attribute.Format
	schema.Title = t.attribute.Title
	schema.MaxLength = t.attribute.MaxLength
	schema.MinLength = t.attribute.MinLength
	schema.Pattern = t.attribute.Pattern
	schema.Maximum = t.attribute.Maximum
	schema.Minimum = t.attribute.Minimum
	schema.MaxItems = t.attribute.MaxItems
	schema.MinItems = t.attribute.MinItems
	schema.Enum = t.attribute.Enum
	return schema
}
