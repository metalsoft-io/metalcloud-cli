package main

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetTableHeader(t *testing.T) {

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      6,
			FieldPrecision: 2,
		},
	}
	expected := "| ID    | LABEL               | INST. |"

	actual := GetTableHeader(schema)

	if actual != expected {
		t.Errorf("Header is not correct, \nexpected:  %s\n     was: %s", expected, actual)
	}
}

func TestGetTableRow(t *testing.T) {
	RegisterTestingT(t)
	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      6,
			FieldPrecision: 2,
		},
		SchemaField{
			FieldName: "INTF",
			FieldSize: 5,
			FieldType: TypeInterface,
		},
	}

	row := []interface{}{10, "test", 33.3, map[string]string{"test": "test1", "test2": "test3"}}

	actual := GetTableRow(row, schema)

	Expect(actual).To(ContainSubstring("test1"))
	Expect(actual).To(ContainSubstring("test3"))
}

func TestGetTableDelimiter(t *testing.T) {
	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      6,
			FieldPrecision: 2,
		},
	}

	expected := "+-------+---------------------+-------+"

	actual := GetTableDelimiter(schema)

	if actual != expected {
		t.Errorf("Delimiter is not correct, \nexpected: %s\n     was: %s", expected, actual)
	}
}

func TestGetTableAsString(t *testing.T) {
	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      6,
			FieldPrecision: 2,
		},
	}

	expected :=
		`+-------+---------------------+-------+
| ID    | LABEL               | INST. |
+-------+---------------------+-------+
| 4     | str                 | 20.10 |
| 5     | st11r               | 22.10 |
| 6     | st11r444            | 2.10  |
+-------+---------------------+-------+
`
	data := [][]interface{}{
		{4, "str", 20.1},
		{5, "st11r", 22.1},
		{6, "st11r444", 2.1},
	}

	actual := GetTableAsString(data, schema)

	if actual != expected {
		t.Errorf("Delimiter is not correct, \nexpected:\n%s\nwas:\n%s", expected, actual)
	}
}

func TestGetTableAsJSONString(t *testing.T) {
	RegisterTestingT(t)
	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      6,
			FieldPrecision: 2,
		},
	}

	data := [][]interface{}{
		{4, "str", 20.1},
		{5, "st11r", 22.1},
		{6, "st11r444", 2.1},
	}

	ret, err := GetTableAsJSONString(data, schema)
	if err != nil {
		t.Errorf("%s", err)
	}

	var m []interface{}

	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())

	Expect(int(m[0].(map[string]interface{})["ID"].(float64))).To(Equal(data[0][0]))
	Expect(m[0].(map[string]interface{})["LABEL"]).To(Equal(data[0][1]))
	Expect(m[2].(map[string]interface{})["LABEL"]).To(Equal(data[2][1]))
	Expect(float64(m[2].(map[string]interface{})["INST."].(float64))).To(Equal(data[2][2]))

}

func TestGetTableAsCSVString(t *testing.T) {

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      6,
			FieldPrecision: 2,
		},
	}

	expected :=
		`ID,LABEL,INST.
4,str,20.100000
5,st11r,22.100000
6,st11r444,2.100000
`

	data := [][]interface{}{
		{4, "str", 20.1},
		{5, "st11r", 22.1},
		{6, "st11r444", 2.1},
	}

	actual, err := GetTableAsCSVString(data, schema)
	if err != nil {
		t.Errorf("%s", err)
	}
	if actual != expected {
		t.Errorf("Delimiter is not correct, \nexpected:\n%s\nwas:\n%s", expected, actual)
	}
}

func TestAdjustFieldSizes(t *testing.T) {
	RegisterTestingT(t)
	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 3,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5, //this is smaller than the largest
		},
		SchemaField{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      0,
			FieldPrecision: 4,
		},
	}

	data := [][]interface{}{
		{4, "12345", 20.1},
		{5, "12", 22.1},
		{6, "123456789", 1.2345},
	}

	AdjustFieldSizes(data, &schema)

	Expect(schema[0].FieldSize).To(Equal(3))
	Expect(schema[1].FieldSize).To(Equal(10))
	Expect(schema[2].FieldSize).To(Equal(8))

}
