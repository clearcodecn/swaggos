package main

import (
	"github.com/clearcodecn/swaggos"
	"github.com/clearcodecn/swaggos/examples/model"
	"github.com/clearcodecn/swaggos/ginwrapper"
	"github.com/gin-gonic/gin"
)

func main() {
	g := ginwrapper.Default()
	doc := g.Doc()
	g.Gin().Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})
	doc.JWT("Authorization")
	doc.HostInfo("https://localhost:8080/", "/api/v1")
	group := g.Group("/api/v1")
	{
		group.GET("/users", listUsers).
			Query("order", swaggos.DescRequired("排序", false)).
			Query("q", swaggos.DescRequired("名称迷糊查询", false)).
			JSON([]model.User{})

		group.POST("/user/create", createUser).
			Body(new(model.User)).JSON(gin.H{"id": 1})

		group.DELETE("/user/*id", deleteUser).
			JSON(gin.H{"id": 1})

		group.PUT("/user/update", createUser).
			Body(new(model.User)).JSON(new(model.User))
	}
	g.ServeDoc()
	g.Gin().Run(":8888")
}

func listUsers(ctx *gin.Context)  {}
func createUser(ctx *gin.Context) {}
func deleteUser(ctx *gin.Context) {}
