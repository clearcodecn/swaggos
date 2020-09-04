package main

import (
	"fmt"
	"github.com/clearcodecn/swaggos"
	"github.com/clearcodecn/swaggos/examples/model"
)

// group example
func main() {
	doc := swaggos.Default()

	doc.HostInfo("localhost:8080", "/api").
		Response(200, newSuccessExample()).
		Response(400, newErrorExample())

	group := doc.Group("/users")
	group.Get("/list").JSON(CommonResponseWithData([]model.User{}))
	group.Post("/create").Body(new(model.User)).JSON(CommonResponseWithData(1))
	group.Put("/update").Body(new(model.User)).JSON(CommonResponseWithData(1))
	// path item
	group.Get("/{id}").JSON(new(model.User))
	group.Delete("/{id}").JSON(CommonResponseWithData(1))

	data, _ := doc.Build()
	fmt.Println(string(data))

	data, _ = doc.Yaml()
	fmt.Println(string(data))
}

type CommonResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func newSuccessExample() interface{} {
	return &CommonResponse{
		Code:    0,
		Data:    []int{1, 2, 3, 4},
		Message: "success",
	}
}

func newErrorExample() interface{} {
	return &CommonResponse{
		Code:    1,
		Message: "error message",
	}
}

func CommonResponseWithData(v interface{}) *CommonResponse {
	return &CommonResponse{
		Data: v,
	}
}
