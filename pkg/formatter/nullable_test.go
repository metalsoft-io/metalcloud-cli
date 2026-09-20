package formatter

import (
	"reflect"
	"testing"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestExtractValueUnwrapsAllNullableTypes(t *testing.T) {
	type record struct {
		Str    sdk.NullableString
		Int64  sdk.NullableInt64
		Int32  sdk.NullableInt32
		Bool   sdk.NullableBool
		Unset  sdk.NullableString
		Null   sdk.NullableInt64
		Plain  string
		Nested *sdk.NullableString
	}
	nested := *sdk.NewNullableString(sdk.PtrString("nested"))
	r := record{
		Str:    *sdk.NewNullableString(sdk.PtrString("hello")),
		Int64:  *sdk.NewNullableInt64(sdk.PtrInt64(42)),
		Int32:  *sdk.NewNullableInt32(sdk.PtrInt32(7)),
		Bool:   *sdk.NewNullableBool(sdk.PtrBool(true)),
		Null:   *sdk.NewNullableInt64(nil),
		Plain:  "plain",
		Nested: &nested,
	}
	v := reflect.ValueOf(r)

	cases := map[string]interface{}{
		"Str":    "hello",
		"Int64":  int64(42),
		"Int32":  int64(7),
		"Bool":   true,
		"Unset":  "",
		"Null":   "",
		"Plain":  "plain",
		"Nested": "nested",
	}
	for field, want := range cases {
		got := extractValue(v.FieldByName(field))
		if got != want {
			t.Errorf("%s: got %#v, want %#v", field, got, want)
		}
	}
}
