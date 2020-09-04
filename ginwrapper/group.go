package ginwrapper

import (
	"github.com/clearcodecn/swaggos"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Context gin.Context

type Group struct {
	router *gin.RouterGroup
	group  *swaggos.Group
}

func NewGroup(basePath string, router *gin.RouterGroup, doc *swaggos.Swaggos) *Group {
	g := new(Group)
	g.router = router
	g.group = swaggos.NewGroup(basePath, doc)
	return g
}

func (g *Group) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return g.router.Use(middleware...)
}

func (g *Group) Group(relativePath string, handlers ...gin.HandlerFunc) *Group {
	router := g.router.Group(relativePath, handlers...)
	gg := NewGroup(relativePath, router, g.group.Swaggos())
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
	return g.group.Post(relativePath)
}

func (g *Group) GET(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.GET(relativePath, handlers...)
	return g.group.Get(relativePath)
}

func (g *Group) DELETE(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.DELETE(relativePath, handlers...)
	return g.group.Delete(relativePath)
}

func (g *Group) PATCH(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.PATCH(relativePath, handlers...)
	return g.group.Patch(relativePath)
}

func (g *Group) PUT(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.PUT(relativePath, handlers...)
	return g.group.Put(relativePath)
}

func (g *Group) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.OPTIONS(relativePath, handlers...)
	return g.group.Options(relativePath)
}

func (g *Group) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return g.router.HEAD(relativePath, handlers...)
}

func (g *Group) Any(relativePath string, handlers ...gin.HandlerFunc) *swaggos.Path {
	g.router.Any(relativePath, handlers...)
	return g.group.Post(relativePath)
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
