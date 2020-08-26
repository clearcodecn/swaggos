package main

import (
	"encoding/json"
	"fmt"
)

type A string
type B []string
type Object struct {
	ObjA string `json:"obj_a"`
	Objb string `json:"objb"`
	ObjC string `json:"obj_c"`
}

type Ts struct {
	A `json:"xxxxx"`
	B
	X string
}

func main() {
	x := Ts{
		A("hello"),
		[]string{"a", "b", "c"},
		"hello",
	}
	xx, _ := json.Marshal(x)
	fmt.Println(string(xx))
}
