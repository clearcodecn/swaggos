package swaggos

import (
	"net/http"
	"testing"
)

func TestUI(t *testing.T) {
	s := http.Server{
		Addr:    ":3333",
		Handler: UI("http://a.com/x.json"),
	}
	s.ListenAndServe()
}
