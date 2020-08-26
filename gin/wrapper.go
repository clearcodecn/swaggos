package gin

import (
	"github.com/clearcodecn/yidoc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Engine struct {
	engine *gin.Engine
	doc    *yidoc.YiDoc
}

func NewEngine(engine *gin.Engine) *Engine {
	return &Engine{
		engine: engine,
		doc:    yidoc.New(),
	}
}

// methods wrapper
func (e *Engine) Get(path string, handlers ...gin.HandlerFunc) *yidoc.Path {
	e.engine.GET(path, handlers...)
	return e.doc.AddPath(http.MethodGet, path)
}
