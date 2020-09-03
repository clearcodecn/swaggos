package main

import (
	"github.com/clearcodecn/swaggos"
	"github.com/clearcodecn/swaggos/ginwrapper"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username     string `json:"username" required:"true"`
	Password     string `json:"password" required:"true" description:"密码" example:"123456" maxLength:"20" minLength:"6" pattern:"[a-zA-Z0-9]{6,20}"`
	Sex          int    `json:"sex" required:"false" default:"1" example:"1" format:"int64"`
	HeadImageURL string `json:"headImageUrl"`

	History string `json:"-"` // ignore
}

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
			JSON([]User{})

		group.POST("/user/create", createUser).
			Body(new(User)).JSON(gin.H{"id": 1})

		group.DELETE("/user/*id", deleteUser).
			JSON(gin.H{"id": 1})

		group.PUT("/user/update", createUser).
			Body(new(User)).JSON(new(User))
	}
	g.ServeDoc()
	g.Gin().Run(":8888")
}

func listUsers(ctx *gin.Context)  {}
func createUser(ctx *gin.Context) {}
func deleteUser(ctx *gin.Context) {}
