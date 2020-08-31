package v2

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/spec"
	"net/http"
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
	Fun
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
		HostInfo("http://localhost:8899/", "/api/v1", spec.InfoProps{})

	d.Get("/{id}").Query("orderBy", Attribute{
		Desc:     "排序",
		Required: false,
		Type:     "string",
		Format:   "string",
	}).
		JSON(new(A))

	data := d.Build()
	fmt.Println(string(data))
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Write(data)
	})
	http.ListenAndServe(":9991", nil)
}
