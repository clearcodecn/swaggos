// package swaggos
package swaggos

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"text/template"
)

var (
	//go:embed swaggerui/*
	swaggerUIFs embed.FS

	//go:embed swaggerui/swagger-initializer.js
	swaggerInitializerJs string
)

func UI(baseURL string) http.Handler {
	tpl, err := template.New("").Parse(swaggerInitializerJs)
	if err != nil {
		panic(err)
	}
	var initBuf bytes.Buffer
	err = tpl.Execute(&initBuf, baseURL)
	if err != nil {
		panic(err)
	}
	swaggerInitializerJs = initBuf.String()

	var staticFS = fs.FS(swaggerUIFs)
	swagFs, err := fs.Sub(staticFS, "swaggerui")
	if err != nil {
		panic(err)
	}
	return &initFileHandlerWrapper{
		h: http.FileServer(http.FS(swagFs)),
	}
}

type initFileHandlerWrapper struct {
	h http.Handler
}

func (h *initFileHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "swagger-initializer.js") {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write([]byte(swaggerInitializerJs))
		return
	}
	h.h.ServeHTTP(w, r)
}
