package doc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Engine struct {
	engine *gin.Engine

	docs []*docItem
}

func (e *Engine) Get(relativePath string, handlers ...gin.HandlerFunc) {
	e.engine.GET(relativePath, handlers...)
	e.newDocItem(http.MethodGet, relativePath)
}

func (e *Engine) newDocItem(method string, path string) *docItem {
	item := newDocItem()
	e.docs = append(e.docs, item)
	return item.addRoute(method, path)
}

type Item struct {
	Method              string
	Path                string
	RequestContentType  string
	ResponseContentType string
	Body                []BasicParam
	Query               []BasicParam
	Form                []BasicParam
	Headers             []BasicParam
	FormData            []BasicParam

	Response interface{}
}

type BasicParam struct {
	Name        string
	Description string
	Type        string
	Required    bool
	Rule        string

	Items []BasicParam
}

type docItem struct {
	item *Item
}

func newDocItem() *docItem {
	return &docItem{item: new(Item)}
}

func (d *docItem) addRoute(method string, path string) *docItem {
	d.item.Method = method
	d.item.Path = path
	return d
}

func (d *docItem) Query(name string, typ string, desc string) *docItem {
	q := BasicParam{
		Name:        name,
		Type:        typ,
		Description: desc,
	}
	d.item.Query = append(d.item.Query, q)
	return d
}

func (d *docItem) RequiredQuery(name string, typ string, desc string) *docItem {
	q := BasicParam{
		Name:        name,
		Type:        typ,
		Description: desc,
		Required:    true,
	}
	d.item.Query = append(d.item.Query, q)
	return d
}

func (d *docItem) RequiredRuledQuery(name string, typ string, desc string, rule string) *docItem {
	q := BasicParam{
		Name:        name,
		Type:        typ,
		Description: desc,
		Required:    true,
		Rule:        rule,
	}
	d.item.Query = append(d.item.Query, q)
	return d
}

func (d *docItem) BindQuery(v interface{}) *docItem {
	return d
}
