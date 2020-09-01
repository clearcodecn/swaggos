package gin

import (
	"github.com/clearcodecn/yidoc"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	engine *gin.Engine
	doc    *yidoc.YiDoc
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
	doc := yidoc.NewYiDoc()
	e.doc = doc
	return e
}

func (e *Engine) Doc() *yidoc.YiDoc {
	return e.doc
}

func (e *Engine) Get(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.GET(path, handlers...)
	return e.doc.Get(path)
}

func (e *Engine) Post(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.POST(path, handlers...)
	return e.doc.Post(path)
}

func (e *Engine) Put(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.PUT(path, handlers...)
	return e.doc.Put(path)
}

func (e *Engine) Patch(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.PATCH(path, handlers...)
	return e.doc.Patch(path)
}

func (e *Engine) Delete(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.DELETE(path, handlers...)
	return e.doc.Delete(path)
}

func (e *Engine) Options(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.OPTIONS(path, handlers...)
	return e.doc.Options(path)
}

func (e *Engine) Head(path string, handlers ...gin.HandlerFunc) {
	e.engine.HEAD(path, handlers...)
}

func (e *Engine) Any(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.Any(path, handlers...)
	e.engine.POST(path, handlers...)
	return e.doc.Get(path)
}

func (e *Engine) Group(path string, handlers ...gin.HandlerFunc) *Group {
	group := e.engine.Group(path, handlers...)
	dgroup := NewGroup(path, group, e.doc)
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
