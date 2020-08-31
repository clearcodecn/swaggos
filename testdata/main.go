package main

import (
	"encoding/json"
	"fmt"
)

type Ref struct {
	Str  string `json:"str"`
	Ref2
}

type Ref2 int

func main() {
	ref := new(Ref)
	d, _ := json.Marshal(ref)
	fmt.Println(string(d))
}
