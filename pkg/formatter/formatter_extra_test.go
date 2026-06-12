package formatter

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestIsNativeFormat(t *testing.T) {
	tests := []struct {
		format string
		want   bool
	}{
		{"json", true},
		{"JSON", true},
		{"yaml", true},
		{"YAML", true},
		{"text", false},
		{"csv", false},
		{"md", false},
		{"", false},
	}
	for _, tc := range tests {
		viper.Set(ConfigFormat, tc.format)
		got := IsNativeFormat()
		if got != tc.want {
			t.Errorf("IsNativeFormat() with format=%q: got %v, want %v", tc.format, got, tc.want)
		}
	}
}

func TestIsTextFormat(t *testing.T) {
	tests := []struct {
		format string
		want   bool
	}{
		{"text", true},
		{"TEXT", true},
		{"", true},
		{"json", false},
		{"yaml", false},
		{"csv", false},
		{"md", false},
	}
	for _, tc := range tests {
		viper.Set(ConfigFormat, tc.format)
		got := IsTextFormat()
		if got != tc.want {
			t.Errorf("IsTextFormat() with format=%q: got %v, want %v", tc.format, got, tc.want)
		}
	}
}

func TestFormatIntegerValue(t *testing.T) {
	tests := []struct {
		in   interface{}
		want string
	}{
		{nil, ""},
		{int(1234), "1,234"},
		{int64(1000000), "1,000,000"},
		{uint32(42), "42"},
		{"500", "500"},
		{"not-a-number", "not-a-number"},
	}
	for _, tc := range tests {
		got := FormatIntegerValue(tc.in)
		if got != tc.want {
			t.Errorf("FormatIntegerValue(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatIdValue(t *testing.T) {
	tests := []struct {
		in   interface{}
		want string
	}{
		{nil, ""},
		{int(1234), "1234"},
		{int64(1000000), "1000000"},
		{uint32(42), "42"},
		{"99", "99"},
	}
	for _, tc := range tests {
		got := FormatIdValue(tc.in)
		if got != tc.want {
			t.Errorf("FormatIdValue(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatBooleanValue(t *testing.T) {
	tests := []struct {
		in   interface{}
		want string
	}{
		{nil, ""},
		{true, "true"},
		{false, "false"},
		{float64(1), "true"},
		{float64(0), "false"},
		{int(1), "true"},
		{int(0), "false"},
		{"true", "true"},
		{"false", "false"},
		{"1", "true"},
		{"0", "false"},
	}
	for _, tc := range tests {
		got := FormatBooleanValue(tc.in)
		if got != tc.want {
			t.Errorf("FormatBooleanValue(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatStringListValue(t *testing.T) {
	tests := []struct {
		in   interface{}
		want string
	}{
		{nil, ""},
		{[]string{"a", "b", "c"}, "a, b, c"},
		{[]string{}, ""},
		{[]string{"single"}, "single"},
		{[]interface{}{"x", "y"}, "x, y"},
		{"plain string", "plain string"},
	}
	for _, tc := range tests {
		got := FormatStringListValue(tc.in)
		if got != tc.want {
			t.Errorf("FormatStringListValue(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestPrintResult_JSON(t *testing.T) {
	viper.Set(ConfigFormat, "json")
	type row struct{ Name string }
	err := PrintResult([]row{{"foo"}}, nil)
	if err != nil {
		t.Errorf("PrintResult json: unexpected error: %v", err)
	}
}

func TestPrintResult_YAML(t *testing.T) {
	viper.Set(ConfigFormat, "yaml")
	type row struct{ Name string }
	err := PrintResult([]row{{"bar"}}, nil)
	if err != nil {
		t.Errorf("PrintResult yaml: unexpected error: %v", err)
	}
}

func TestPrintResult_CSV(t *testing.T) {
	viper.Set(ConfigFormat, "csv")
	type row struct{ Name string }
	err := PrintResult([]row{{"baz"}}, nil)
	if err != nil {
		t.Errorf("PrintResult csv: unexpected error: %v", err)
	}
}

func TestPrintResult_MD(t *testing.T) {
	viper.Set(ConfigFormat, "md")
	type row struct{ Name string }
	err := PrintResult([]row{{"qux"}}, nil)
	if err != nil {
		t.Errorf("PrintResult md: unexpected error: %v", err)
	}
}

func TestPrintResult_Text(t *testing.T) {
	viper.Set(ConfigFormat, "text")
	type row struct{ Name string }
	err := PrintResult([]row{{"text-val"}}, nil)
	if err != nil {
		t.Errorf("PrintResult text: unexpected error: %v", err)
	}
}

func TestPrintResult_UnsupportedFormat(t *testing.T) {
	viper.Set(ConfigFormat, "xml")
	err := PrintResult("anything", nil)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

// TestPrintResult_InterfaceFields_NoPanic — regression test: raw structs use
// interface{} fields; table formats must not panic on FieldByName over them.
func TestPrintResult_InterfaceFields_NoPanic(t *testing.T) {
	type rawRow struct {
		Id     interface{}
		Label  *string
		Nested interface{}
	}
	label := "x"
	rows := []rawRow{
		{Id: float64(1), Label: &label, Nested: map[string]interface{}{"inner": map[string]interface{}{"deep": "v"}}},
		{Id: "str-id", Label: nil, Nested: nil},
	}
	cfg := &PrintConfig{FieldsConfig: map[string]RecordFieldConfig{
		"Id":    {Order: 1},
		"Label": {Order: 2},
		// Dotted path through interface{} holding a JSON map — panicked before fix.
		"Nested.Inner.Deep": {Order: 3},
	}}
	for _, format := range []string{"text", "csv", "md", "json", "yaml"} {
		viper.Set(ConfigFormat, format)
		if err := PrintResult(rows, cfg); err != nil {
			t.Errorf("format %s: unexpected error: %v", format, err)
		}
	}
}

// TestLocateField_MapTraversal — dotted paths resolve through raw JSON maps
// using camelCase keys.
func TestLocateField_MapTraversal(t *testing.T) {
	val := map[string]interface{}{
		"ethernetFabric": map[string]interface{}{"fabricType": "ethernet"},
	}
	type holder struct{ Cfg interface{} }
	h := holder{Cfg: val}
	field := locateField("Cfg.EthernetFabric.FabricType", reflect.ValueOf(h))
	if !field.IsValid() {
		t.Fatal("expected to resolve Cfg.EthernetFabric.FabricType through map")
	}
	if field.Interface() != "ethernet" {
		t.Errorf("expected 'ethernet', got %v", field.Interface())
	}
}

// TestLocateField_MapMissingKey — missing map key with camelCase fallback must
// return invalid, not panic (MapIndex on zero Value regression).
func TestLocateField_MapMissingKey(t *testing.T) {
	type holder struct{ Cfg interface{} }
	h := holder{Cfg: map[string]interface{}{"otherKey": "x"}}
	field := locateField("Cfg.EthernetFabric.FabricType", reflect.ValueOf(h))
	if field.IsValid() {
		t.Errorf("expected invalid field for missing key, got %v", field)
	}
}
