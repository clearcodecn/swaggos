package swaggos

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type AS string
type Arr []string
type MAp map[string]struct{}
type Fun func()
type ArrayMap []map[string]struct{}

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
	yd := new(Swaggos)
	yd.Define(new(A))

	def, err := json.MarshalIndent(yd.definitions, "", "  ")
	fmt.Println(string(def), err)
}

func TestDocs(t *testing.T) {
	d := NewSwaggo()
	d.JWT("Token").
		Oauth2Password("https://www.oauth2.com/token", []string{"openid"}).
		HostInfo("localhost:8899", "/api/v1")

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
	ObjectArray       []TObject
	ObjectArrayArray  [][]*TObject
	ObjectArrayArray2 [][][][]*TObject
}

type TObject struct {
	Attr string `json:"attr"`
}

func TestBuildSchema(t *testing.T) {
	y := NewSwaggo()
	v := new(ArrayTest)
	y.Define(v)
	data, _ := y.Build()
	fmt.Println(string(data))
}
