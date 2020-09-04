package ginwrapper

import (
	"fmt"
	"github.com/clearcodecn/swaggos"
	"github.com/gin-gonic/gin"
	"strings"
)

type Engine struct {
	engine *gin.Engine
	doc    *swaggos.Swaggos
}

func Default() *Engine {
	g := gin.Default()
	return New(g)
}

func (e *Engine) Gin() *gin.Engine {
	return e.engine
}

func New(g *gin.Engine) *Engine {
	e := new(Engine)
	e.engine = g
	doc := swaggos.Default()
	e.doc = doc
	return e
}

func (e *Engine) Doc() *swaggos.Swaggos {
	return e.doc
}

func (e *Engine) Get(path string, handlers ...gin.HandlerFunc) *swaggos.Path {
	e.engine.GET(path, handlers...)
	return e.doc.Get(path)
}

func (e *Engine) Post(path string, handlers ...gin.HandlerFunc) *swaggos.Path {
	e.engine.POST(path, handlers...)
	return e.doc.Post(trimPath(path))
}

func (e *Engine) Put(path string, handlers ...gin.HandlerFunc) *swaggos.Path {
	e.engine.PUT(path, handlers...)
	return e.doc.Put(trimPath(path))
}

func (e *Engine) Patch(path string, handlers ...gin.HandlerFunc) *swaggos.Path {
	e.engine.PATCH(path, handlers...)
	return e.doc.Patch(trimPath(path))
}

func (e *Engine) Delete(path string, handlers ...gin.HandlerFunc) *swaggos.Path {
	e.engine.DELETE(path, handlers...)
	return e.doc.Delete(trimPath(path))
}

func (e *Engine) Options(path string, handlers ...gin.HandlerFunc) *swaggos.Path {
	e.engine.OPTIONS(path, handlers...)
	return e.doc.Options(trimPath(path))
}

func (e *Engine) Head(path string, handlers ...gin.HandlerFunc) {
	e.engine.HEAD(path, handlers...)
}

func (e *Engine) Any(path string, handlers ...gin.HandlerFunc) *swaggos.Path {
	e.engine.Any(path, handlers...)
	e.engine.POST(path, handlers...)
	return e.doc.Get(trimPath(path))
}

func (e *Engine) Group(path string, handlers ...gin.HandlerFunc) *Group {
	group := e.engine.Group(path, handlers...)
	dgroup := NewGroup(trimPath(path), group, e.doc)
	return dgroup
}

func (e *Engine) ServeDoc() {
	e.engine.GET("/_/swagger.json", func(ctx *gin.Context) {
		data, err := e.doc.Build()
		if err != nil {
			ctx.AbortWithError(400, err)
			return
		}
		ctx.Writer.Write(data)
	})
}

func trimPath(path string) string {
	path = "/" + strings.Trim(path, "/")
	arr := strings.Split(path, "/")
	var pathArr []string
	for _, a := range arr {
		if strings.Contains(a, "*") {
			a = strings.Replace(a, "*", "", -1)
			pathArr = append(pathArr, fmt.Sprintf("{%s}", a))
		} else {
			pathArr = append(pathArr, a)
		}
	}
	return strings.Join(pathArr, "/")
}
