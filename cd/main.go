package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	g := gin.Default()

	d := new(doc)

	d.engine = g

	d.Get("/", func(context *gin.Context) {}).Query(Name("hello")).Query(Name("world")).Body(Object(nil))

	gr := d.Group("/api/v1")
	{
		gr.Get("/", func(context *gin.Context) {}).Body(Object(nil)).Response(Object(nil))
	}
}

type doc struct {
	engine *gin.Engine
}

func (d *doc) Get(path string, handlers ...gin.HandlerFunc) *Request {
	d.engine.GET(path, handlers...)
	req := new(Request)
	req.Path = path
	req.Method = http.MethodGet
	return req
}

func (d *doc) Group(path string, handlers ...gin.HandlerFunc) *Group {
	g := d.engine.Group(path, handlers...)
	return &Group{group: g}
}


type Group struct {
	group *gin.RouterGroup
}

func (g *Group) Get(path string, handlers ...gin.HandlerFunc) *Request {
	g.group.GET(path, handlers...)
	req := new(Request)
	req.Path = path
	req.Method = http.MethodGet
	return req
}

type Request struct {
	Path   string
	Method string
	doc    *doc
}


type nameArg struct {
	*emptyArg
	name string
}

func (n *nameArg) Name() string {
	return n.name
}

func Name(name string) Arg {
	return &nameArg{name: name}
}

type Argument struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Rule        string `json:"rule"`
	In          string `json:"in"`
}

func Object(v interface{}) Arg {

}

func (r *Request) Query(arg Arg, args ...Arg) *Request {
	var argu = new(Argument)
	for _, a := range append([]Arg{arg}, args...) {
		switch a.(type) {
		case *nameArg:
			argu.Name = a.Name()
		}
	}
	//
}

func (r *Request) Body(arg Arg, args ...Arg) *Request {
	var argu = new(Argument)
	for _, a := range append([]Arg{arg}, args...) {
		switch a.(type) {
		case *nameArg:
			argu.Name = a.Name()
		}
	}
	//
}

func (r *Request) Response(v interface{}) {
}
