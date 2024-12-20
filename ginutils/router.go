package ginutils

import (
	"fmt"
	"github.com/clearcodecn/swaggos"
	"github.com/gin-gonic/gin"
	"strings"
)

func Serve(path string, swag *swaggos.Swaggos, e *gin.Engine, basicAuthAccounts gin.Accounts) {
	path = "/" + strings.Trim(path, "/")
	if len(basicAuthAccounts) > 0 {
		e.Use(gin.BasicAuth(basicAuthAccounts))
	}
	// swagger json 服务
	e.GET(fmt.Sprintf("%s/_doc", path), gin.WrapH(swag))

	// swagger ui 服务
	e.Any(fmt.Sprintf("%s/*action", path), gin.WrapH(swaggos.UI(path, path+"/_doc")))

	fmt.Println("swagger server at: " + path)
}
