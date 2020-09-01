package gin

import (
	"github.com/clearcodecn/yidoc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Group struct {
	basePath string
	router   *gin.RouterGroup
	doc      *yidoc.YiDoc
}

func NewGroup(basePath string, router *gin.RouterGroup, doc *yidoc.YiDoc) *Group {
	g := new(Group)
	g.basePath = basePath
	g.router = router
	g.doc = doc
	return g
}

func (group *Group) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return group.router.Use(middleware...)
}

func (group *Group) Group(relativePath string, handlers ...gin.HandlerFunc) *Group {
	router := group.router.Group(relativePath, handlers...)
	g := NewGroup(group.basePath+relativePath, router, group.doc)
	return g
}

func (group *Group) BasePath() string {
	return group.router.BasePath()
}

func (group *Group) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return group.router.Handle(httpMethod, relativePath, handlers...)
}

func (group *Group) POST(relativePath string, handlers ...gin.HandlerFunc) *yidoc.Path {
	group.router.POST(relativePath, handlers...)
	return group.doc.Post(group.basePath + relativePath)
}

func (group *Group) GET(relativePath string, handlers ...gin.HandlerFunc) *yidoc.Path {
	group.router.GET(relativePath, handlers...)
	return group.doc.Get(group.basePath + relativePath)
}

func (group *Group) DELETE(relativePath string, handlers ...gin.HandlerFunc) *yidoc.Path {
	group.router.DELETE(relativePath, handlers...)
	return group.doc.Delete(group.basePath + relativePath)
}

func (group *Group) PATCH(relativePath string, handlers ...gin.HandlerFunc) *yidoc.Path {
	group.router.PATCH(relativePath, handlers...)
	return group.doc.Patch(group.basePath + relativePath)
}

func (group *Group) PUT(relativePath string, handlers ...gin.HandlerFunc) *yidoc.Path {
	group.router.PUT(relativePath, handlers...)
	return group.doc.Put(group.basePath + relativePath)
}

func (group *Group) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) *yidoc.Path {
	group.router.OPTIONS(relativePath, handlers...)
	return group.doc.Options(group.basePath + relativePath)
}

func (group *Group) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return group.router.HEAD(relativePath, handlers...)
}

func (group *Group) Any(relativePath string, handlers ...gin.HandlerFunc) *yidoc.Path {
	group.router.Any(relativePath, handlers...)
	return group.doc.Post(group.basePath + relativePath)
}

func (group *Group) StaticFile(relativePath, filepath string) gin.IRoutes {
	return group.router.StaticFile(relativePath, filepath)
}

func (group *Group) Static(relativePath, root string) gin.IRoutes {
	return group.router.Static(relativePath, root)
}

func (group *Group) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return group.router.StaticFS(relativePath, fs)
}
