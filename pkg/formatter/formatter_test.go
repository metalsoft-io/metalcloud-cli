package formatter

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type testStruct struct {
	ID    int
	Name  string
	Extra string
}

type paginatedStruct struct {
	Data []testStruct
}

func TestGenerateTable_Struct(t *testing.T) {
	obj := testStruct{ID: 7, Name: "struct"}
	tbl := generateTable(obj, nil)
	var buf bytes.Buffer
	tbl.SetOutputMirror(&buf)
	tbl.Render()
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "NAME") {
		t.Errorf("expected headers in output, got: %s", out)
	}
}

func TestGenerateTable_Slice(t *testing.T) {
	obj := []testStruct{{ID: 8, Name: "a"}, {ID: 9, Name: "b"}}
	tbl := generateTable(obj, nil)
	var buf bytes.Buffer
	tbl.SetOutputMirror(&buf)
	tbl.Render()
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Errorf("expected slice rows in output, got: %s", out)
	}
}

func TestGenerateTable_PaginatedStruct(t *testing.T) {
	obj := paginatedStruct{Data: []testStruct{{ID: 10, Name: "pag"}}}
	tbl := generateTable(obj, nil)
	var buf bytes.Buffer
	tbl.SetOutputMirror(&buf)
	tbl.Render()
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "pag") {
		t.Errorf("expected paginated data in output, got: %s", out)
	}
}

func TestGenerateTable_String(t *testing.T) {
	obj := "hello"
	tbl := generateTable(obj, nil)
	var buf bytes.Buffer
	tbl.SetOutputMirror(&buf)
	tbl.Render()
	out := buf.String()
	if !strings.Contains(out, "RESULT") || !strings.Contains(out, "hello") {
		t.Errorf("expected string result in output, got: %s", out)
	}
}

func TestGenerateTable_Map(t *testing.T) {
	obj := map[string]interface{}{"foo": 1, "bar": 2}
	tbl := generateTable(obj, nil)
	var buf bytes.Buffer
	tbl.SetOutputMirror(&buf)
	tbl.Render()
	out := buf.String()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Errorf("expected map keys in output, got: %s", out)
	}
}

func TestGetPaginatedData(t *testing.T) {
	slice := []testStruct{{ID: 1}}
	val, ok := getPaginatedData(slice)
	if !ok || val.Len() != 1 {
		t.Errorf("expected slice to be paginated data")
	}
	obj := paginatedStruct{Data: []testStruct{{ID: 2}}}
	val, ok = getPaginatedData(obj)
	if !ok || val.Len() != 1 {
		t.Errorf("expected struct with Data to be paginated data")
	}
	obj2 := struct{ Foo int }{Foo: 1}
	_, ok = getPaginatedData(obj2)
	if ok {
		t.Errorf("expected struct without Data to not be paginated data")
	}
	_, ok = getPaginatedData(123)
	if ok {
		t.Errorf("expected non-struct to not be paginated data")
	}
}

func TestGetFieldNamesAndValues(t *testing.T) {
	obj := testStruct{ID: 1, Name: "foo"}
	names, values, _ := getFieldNamesAndValues(obj, nil)
	if len(names) != 3 || len(values) != 3 {
		t.Errorf("expected 3 fields, got %d", len(names))
	}
	ptr := &obj
	names2, values2, _ := getFieldNamesAndValues(ptr, nil)
	if len(names2) != 3 || len(values2) != 3 {
		t.Errorf("expected 3 fields for pointer, got %d", len(names2))
	}
	_, _, _ = getFieldNamesAndValues(123, nil) // should not panic
}

func TestGetColumnsCount(t *testing.T) {
	cfg := map[string]RecordFieldConfig{
		"A": {Order: 1},
		"B": {Order: 2, Hidden: true},
		"C": {Order: 3, InnerFields: map[string]RecordFieldConfig{
			"D": {Order: 1},
		}},
	}
	count, maxOrder := getColumnsCount(&cfg)
	if count == 0 || maxOrder == 0 {
		t.Errorf("expected nonzero count and maxOrder")
	}
}

func TestPopulate(t *testing.T) {
	cfg := map[string]RecordFieldConfig{
		"ID":   {Order: 1},
		"Name": {Order: 2},
	}
	obj := testStruct{ID: 1, Name: "foo"}
	names := make(table.Row, 2)
	values := make(table.Row, 2)
	var configs []table.ColumnConfig
	populate(obj, &cfg, &names, &values, &configs)
	if names[0] != "ID" || values[1] != "foo" {
		t.Errorf("populate did not fill names/values as expected")
	}
}

func TestAddField(t *testing.T) {
	names := make(table.Row, 1)
	values := make(table.Row, 1)
	var configs []table.ColumnConfig
	cfg := RecordFieldConfig{Order: 1, Title: "MyID", Transformer: func(i interface{}) string { return "x" }, MaxWidth: 10}
	addField(cfg, "ID", 42, &names, &values, &configs)
	if names[0] != "MyID" || values[0] != 42 {
		t.Errorf("addField did not set names/values as expected")
	}
	if len(configs) == 0 {
		t.Errorf("expected column config")
	}
	cfg2 := RecordFieldConfig{Order: 1, Hidden: true}
	addField(cfg2, "ID", 42, &names, &values, &configs) // should not set
}

func TestExtractValue(t *testing.T) {
	s := "abc"
	b := true
	i := 42
	u := uint(42)
	f := 3.14
	arr := [2]int{1, 2}
	slice := []string{"a", "b"}
	m := map[string]interface{}{"x": 1}
	type dummy struct{}
	ptr := &s
	if extractValue(reflect.ValueOf(s)) != "abc" {
		t.Errorf("extractValue string failed")
	}
	if extractValue(reflect.ValueOf(b)) != true {
		t.Errorf("extractValue bool failed")
	}
	if extractValue(reflect.ValueOf(i)) != int64(42) {
		t.Errorf("extractValue int failed")
	}
	if extractValue(reflect.ValueOf(u)) != uint64(42) {
		t.Errorf("extractValue uint failed")
	}
	if extractValue(reflect.ValueOf(f)) != 3.14 {
		t.Errorf("extractValue float failed")
	}
	if v, ok := extractValue(reflect.ValueOf(arr)).([]interface{}); !ok || len(v) != 2 {
		t.Errorf("extractValue array failed")
	}
	if v, ok := extractValue(reflect.ValueOf(slice)).([]interface{}); !ok || len(v) != 2 {
		t.Errorf("extractValue slice failed")
	}
	if v, ok := extractValue(reflect.ValueOf(m)).([]string); !ok || len(v) != 1 {
		t.Errorf("extractValue map failed")
	}
	if extractValue(reflect.ValueOf(dummy{})) != nil {
		t.Errorf("extractValue struct failed")
	}
	if extractValue(reflect.ValueOf(ptr)) != "abc" {
		t.Errorf("extractValue pointer failed")
	}
	var invalid reflect.Value
	if extractValue(invalid) != nil {
		t.Errorf("extractValue invalid failed")
	}
}

func TestFormatStatusValue(t *testing.T) {
	statuses := []string{
		"available", "ready", "used", "unavailable", "registering", "cleaning", "cleaning_required",
		"updating_firmware", "pending_registration", "used_registering", "used_diagnostics",
		"decommissioned", "removed_from_rack", "defective", "active", "ordered", "draft", "unknown",
	}
	for _, s := range statuses {
		out := FormatStatusValue(s)
		if !strings.HasPrefix(out, text.EscapeStart) {
			t.Errorf("FormatStatusValue(%q) missing color escape in output '%s'", s, out)
		}
		if !strings.Contains(out, s) {
			t.Errorf("FormatStatusValue(%q) missing value in output '%s'", s, out)
		}
	}
	if FormatStatusValue("random string") != text.FgYellow.Sprintf("%s", "random string") {
		t.Errorf("FormatStatusValue random string failed")
	}
	if FormatStatusValue(123) != "123" {
		t.Errorf("FormatStatusValue non-string failed")
	}
}

func TestFormatDateTimeValue(t *testing.T) {
	now := time.Now().UTC()
	str := now.Format("2006-01-02T15:04:05Z")
	out := FormatDateTimeValue(str)
	if out == "" {
		t.Errorf("FormatDateTimeValue valid string failed")
	}
	str2 := now.Format("2006-01-02T15:04:05.000Z")
	out2 := FormatDateTimeValue(str2)
	if out2 == "" {
		t.Errorf("FormatDateTimeValue valid string with ms failed")
	}
	str3 := "0000-00-00T00:00:00Z"
	if FormatDateTimeValue(str3) != "" {
		t.Errorf("FormatDateTimeValue zero date failed")
	}
	if FormatDateTimeValue(now) == "" {
		t.Errorf("FormatDateTimeValue time.Time failed")
	}
	if FormatDateTimeValue(nil) != "" {
		t.Errorf("FormatDateTimeValue nil failed")
	}
	if FormatDateTimeValue(123) != "123" {
		t.Errorf("FormatDateTimeValue fallback failed")
	}
}
