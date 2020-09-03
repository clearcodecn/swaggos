package ginwrapper

import (
	"fmt"
	"github.com/clearcodecn/swaggos"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Context gin.Context

type Group struct {
	basePath string
	router   *gin.RouterGroup
	doc      *swaggos.Swaggo
}

func NewGroup(basePath string, router *gin.RouterGroup, doc *swaggos.Swaggo) *Group {
	g := new(Group)
	g.basePath = basePath
	g.router = router
	g.doc = doc
	return g
}

func (g *Group) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return g.router.Use(middleware...)
}

func (g *Group) Group(relativePath string, handlers ...gin.HandlerFunc) *Group {
	router := g.router.Group(relativePath, handlers...)
	gg := NewGroup(g.trimPath(g.basePath+relativePath), router, g.doc)
	return gg
}

func (g *Group) BasePath() string {
	return g.router.BasePath()
}

func (g *Group) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return g.router.Handle(httpMethod, relativePath, handlers...)
}

func (g *Group) POST(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.POST(relativePath, handlers...)
	return g.doc.Post(g.trimPath(relativePath))
}

func (g *Group) GET(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.GET(relativePath, handlers...)
	return g.doc.Get(g.trimPath(relativePath))
}

func (g *Group) DELETE(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.DELETE(relativePath, handlers...)
	return g.doc.Delete(g.trimPath(relativePath))
}

func (g *Group) PATCH(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.PATCH(relativePath, handlers...)
	return g.doc.Patch(g.trimPath(relativePath))
}

func (g *Group) PUT(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.PUT(relativePath, handlers...)
	return g.doc.Put(g.trimPath(relativePath))
}

func (g *Group) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.OPTIONS(relativePath, handlers...)
	return g.doc.Options(g.trimPath(relativePath))
}

func (g *Group) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return g.router.HEAD(relativePath, handlers...)
}

func (g *Group) Any(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.Any(relativePath, handlers...)
	return g.doc.Post(g.trimPath(relativePath))
}

func (g *Group) StaticFile(relativePath, filepath string) gin.IRoutes {
	return g.router.StaticFile(relativePath, filepath)
}

func (g *Group) Static(relativePath, root string) gin.IRoutes {
	return g.router.Static(relativePath, root)
}

func (g *Group) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return g.router.StaticFS(relativePath, fs)
}

func (g *Group) trimPath(path string) string {
	if len(path) == 0 {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	path = strings.TrimPrefix(path, g.basePath)
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
