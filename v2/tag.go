package v2

import (
	"reflect"
	"strings"
)

type tags struct {
	tag map[string]string
}

func newTags(tag reflect.StructTag) *tags {
	jsonTag := tag.Get("json")
	formTag := tag.Get("form")
	docTag := tag.Get("doc")

	ts := new(tags)
	ts.tag = make(map[string]string)
	ts.tag["json"] = jsonTag
	ts.tag["form"] = formTag
	ts.tag["doc"] = docTag
	_, ok := tag.Lookup("required")
	if ok {
		ts.tag["required"] = ""
	}

	return ts
}

func (t *tags) name() string {
	tg := t.tag["json"]
	if tg == "-" || tg == "" {
		tg = t.tag["form"]
	}
	if tg == "" {
		tg = t.tag["query"]
	}
	return strings.Split(tg, ",")[0]
}

func (t *tags) ignore() bool {
	return t.tag["doc"] == "-"
}

func (t *tags) required() bool {
	arr := strings.Split(t.tag["doc"], ",")
	for _, a := range arr {
		if a == "required" {
			return true
		}
	}
	if _, ok := t.tag["required"]; ok {
		return true
	}
	return false
}
