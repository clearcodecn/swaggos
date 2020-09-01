package main

import (
	"github.com/clearcodecn/yidoc"
	"github.com/clearcodecn/yidoc/examples/dso/model"
	"github.com/clearcodecn/yidoc/gin"
	gin2 "github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
)

type BadReqErr struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func main() {

	g := gin.Default()
	doc := g.Doc()
	g.Gin().Use(func(context *gin2.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})
	info := spec.InfoProps{
		Version:        "0.1",
		Title:          "api",
		Description:    "api for group server,",
		Contact:        &spec.ContactInfo{ContactInfoProps: spec.ContactInfoProps{Name: "developer", Email: "vardump@foxmail.com"}},
	}
	doc.HostInfo("localhost:9999", "/api/v1", info).
		Oauth2("https://oauth.token.url", []string{"openid"}, []string{"openid"}).
		JWT("Authorization")

	gr := g.Group("/api/v1")
	{
		gr.POST("/group/create", helloHandler).
			Summary("group").
			Description("create a new group").
			Tag("group").
			Form("title", yidoc.Attribute{
				Description: "group title",
				Required:    true,
			}).
			Form("group file", yidoc.Attribute{
				Type:        yidoc.File,
				Required:    true,
				Description: "sequence file",
			}).
			JSON(new(model.Group)).
			BadRequest(new(BadReqErr))
	}

	g.ServeDoc()
	g.Gin().Run(":9992")
}

func helloHandler(ctx *gin2.Context) {}
