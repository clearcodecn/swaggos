package main

import (
	"fmt"
	"github.com/clearcodecn/swaggos"
	"github.com/clearcodecn/swaggos/examples/model"
)

func main() {
	doc := swaggos.Default()

	doc.HostInfo("localhost:8080", "/api")
	doc.Get("/users").JSON([]model.RuleUser{})

	data, _ := doc.Yaml()
	fmt.Println(string(data))
}
