package main

import (
	"reflect"
	"testing"
	"time"
)

func TestTableSortWithSchema(t *testing.T) {

	data := [][]interface{}{
		{4, "str", 20.1},
		{6, "st11r", 22.1},
		{5, "wt11r444", 2.3},
		{5, "wt11r444", 2.1},
		{5, "at11r43", 2.2},
		{4, "xxxx", 2.2},
	}

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

	TableSorter(schema).OrderBy("LABEL", "ID", "INST.").Sort(data)

	expected := [][]interface{}{
		{5, "at11r43", 2.2},
		{6, "st11r", 22.1},
		{4, "str", 20.1},
		{5, "wt11r444", 2.1},
		{5, "wt11r444", 2.3},
		{4, "xxxx", 2.2},
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("the sorted array was not correct \nwas:\n%+v\n expected\n %+v\n", data, expected)
	}
}

func TestTableSortWithSchemaWithDateTime(t *testing.T) {

	data := [][]interface{}{
		{4, "str", "2013-11-29T13:00:01Z"},
		{6, "st11r", "2013-11-29T13:00:02Z"},
		{6, "st11r", "2014-11-29T13:00:03Z"},
		{6, "st11r", "2012-11-29T13:00:03Z"},
	}

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
			FieldName:   "DATE",
			FieldType:   TypeDateTime,
			FieldSize:   6,
			FieldFormat: defaultTimeFormat,
		},
	}

	TableSorter(schema).OrderBy("DATE").Sort(data)

	expected := [][]interface{}{
		{6, "st11r", "2012-11-29T13:00:03Z"},
		{4, "str", "2013-11-29T13:00:01Z"},
		{6, "st11r", "2013-11-29T13:00:02Z"},
		{6, "st11r", "2014-11-29T13:00:03Z"},
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("the sorted array was not correct \nwas:\n%+v\n expected\n %+v\n", data, expected)
	}
}

func TestDefaultTimeFormat(t *testing.T) {

	layout := defaultTimeFormat

	s := "2012-11-29T13:00:03Z"

	tm, err := time.Parse(layout, s)

	if err != nil {
		t.Errorf("error converting time string %s", err)
	}

	if tm.Year() != 2012 || tm.Second() != 3 {
		t.Error("Date was not parsed correctly")
	}

}
