package main

import (
	"encoding/json"
	"fmt"
)

type Ref struct {
	Str  string `json:"str"`
	Ref2
}

type Ref2 struct {
	Bar string `json:"bar"`
}

func main() {
	ref := new(Ref)
	d, _ := json.Marshal(ref)
	fmt.Println(string(d))
}
