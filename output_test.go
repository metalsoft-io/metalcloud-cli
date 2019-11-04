package main

import "testing"

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

	expected := "| 10    | test                | 33.30 | map[test :test1 test2:test3]|"
	row := []interface{}{10, "test", 33.3, map[string]string{"test": "test1", "test2": "test3"}}

	actual := GetTableRow(row, schema)

	if actual != expected {
		t.Errorf("Row is not correct, \nexpected: %s\n     was: %s", expected, actual)
	}
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
		`[
	{
		"ID": 6,
		"INST.": 2.1,
		"LABEL": "st11r444"
	},
	{
		"ID": 6,
		"INST.": 2.1,
		"LABEL": "st11r444"
	},
	{
		"ID": 6,
		"INST.": 2.1,
		"LABEL": "st11r444"
	}
]`
	data := [][]interface{}{
		{4, "str", 20.1},
		{5, "st11r", 22.1},
		{6, "st11r444", 2.1},
	}

	actual, err := GetTableAsJSONString(data, schema)
	if err != nil {
		t.Errorf("%s", err)
	}
	if actual != expected {
		t.Errorf("Delimiter is not correct, \nexpected:\n%s\nwas:\n%s", expected, actual)
	}
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
