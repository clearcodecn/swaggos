package yidoc

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"testing"
)

type AS string
type Arr []string
type MAp map[string]struct{}
type Fun func()
type ArrayMap []map[string]struct{}

/**
{
	type: string
	name: "",
}
*/

type A struct {
	Int                  int `json:"int" doc:"required,number,整型数据"`
	Str                  string
	IntArray             [3]int
	IntString            [3]string
	IntArrayArray        [][][][][][]int
	ExampleArray         [][][]ExampleObject
	ExamplePtrArrayArray [][][][][][]*ExampleObject
	ObjectArray          [3]ObjectTest
	ObjectSlice          []*ExampleObject
	Object               *ExampleObject
	Bool                 bool
	Interface            interface{}
	Map                  map[string]interface{}
	ExampleObject
	MAp
	Arr
	ObjectOne ObjectTest
	AS
	Fun `json:"-"`
}

type ObjectTest struct {
	Int         int
	Str         string
	IntArray    [3]int
	IntString   [3]string
	Slice       []string
	ObjectSlice []*ExampleObject
	Object      *ExampleObject
	Bool        bool
	Interface   interface{}
	Map         map[string]interface{}
}

type ExampleObject struct {
	UserName string
	Password string
}

func TestYiDoc(t *testing.T) {
	yd := new(YiDoc)
	yd.Define("Object", new(A))

	def, err := json.MarshalIndent(yd.definitions, "", "  ")
	fmt.Println(string(def), err)
}

func TestDocs(t *testing.T) {
	d := NewYiDoc()
	d.JWT("Token").
		Oauth2("https://www.oauth2.com/token", []string{"openid"}, []string{"read", "write"}).
		HostInfo("localhost:8899", "/api/v1", spec.InfoProps{})

	d.Get("/{id}").Query("orderBy", Attribute{
		Description: "排序",
		Required:    false,
		Type:        "string",
		Format:      "string",
	}).
		Description("排序的用户").
		Tag("orders").
		Summary("排序").
		JSON(new(A))

	data, err := d.Build()
	fmt.Println(string(data))
	require.Nil(t, err)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Write(data)
	})
	http.ListenAndServe(":9991", nil)
}

type ArrayTest struct {
	//ObjectArray      []Object
	ObjectArrayArray [][][][]Object
}

type Object struct {
	Attr string `json:"attr"`
}

func TestBuildSchema(t *testing.T) {
	y := NewYiDoc()
	v := new(ArrayTest)

	refName := reflect.Indirect(reflect.ValueOf(v)).Type().Name()
	y.Define(refName, v)
	data, _ := y.Build()
	fmt.Println(string(data))
}
