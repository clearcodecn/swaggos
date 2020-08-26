package doc

type FieldType int

const (
	FieldTypeNull FieldType = iota
	FieldTypeArray
	FieldTypeBoolean
	FieldTypeInteger
	FieldTypeNumber
	FieldTypeObject
	FieldTypeString
)

var fieldTypeNames = map[string]FieldType{
	"null":    FieldTypeNull,
	"array":   FieldTypeArray,
	"boolean": FieldTypeBoolean,
	"integer": FieldTypeInteger,
	"number":  FieldTypeNumber,
	"object":  FieldTypeObject,
	"string":  FieldTypeString,
}

func FieldTypeByName(name string) FieldType {
	return fieldTypeNames[name]
}
