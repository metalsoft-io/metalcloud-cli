package cmd

import (
	"testing"
)

func TestIsValidMAC(t *testing.T) {
	valid := []string{
		"aa:bb:cc:dd:ee:ff",
		"AA:BB:CC:DD:EE:FF",
		"00:11:22:33:44:55",
		"a1:b2:c3:d4:e5:f6",
		"aa-bb-cc-dd-ee-ff",
		"AA-BB-CC-DD-EE-FF",
	}
	for _, mac := range valid {
		if !isValidMAC(mac) {
			t.Errorf("expected %q to be valid", mac)
		}
	}

	invalid := []string{
		"",
		"aa:bb:cc:dd:ee",        // too short
		"aa:bb:cc:dd:ee:ff:00",  // too long
		"gg:bb:cc:dd:ee:ff",     // bad hex char
		"aa:bb:cc:dd:ee:f",      // odd nibble
		"aabbccddeeff",          // no separator
		"aa bb cc dd ee ff",     // space separator
		"aa:bb:cc:dd:ee:fg",     // g not hex
	}
	for _, mac := range invalid {
		if isValidMAC(mac) {
			t.Errorf("expected %q to be invalid", mac)
		}
	}
}

func TestCountTrue(t *testing.T) {
	tests := []struct {
		vals []bool
		want int
	}{
		{nil, 0},
		{[]bool{false}, 0},
		{[]bool{true}, 1},
		{[]bool{true, true}, 2},
		{[]bool{false, false, false}, 0},
		{[]bool{true, false, true}, 2},
		{[]bool{true, true, true, true}, 4},
	}
	for _, tc := range tests {
		got := countTrue(tc.vals...)
		if got != tc.want {
			t.Errorf("countTrue(%v) = %d, want %d", tc.vals, got, tc.want)
		}
	}
}

func TestParseMACToIPJSON(t *testing.T) {
	t.Run("valid map", func(t *testing.T) {
		data := []byte(`{"AA:BB:CC:DD:EE:FF":"10.0.0.1","11:22:33:44:55:66":"192.168.1.1"}`)
		m, err := parseMACToIPJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(m) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(m))
		}
		if m["AA:BB:CC:DD:EE:FF"] != "10.0.0.1" {
			t.Errorf("unexpected IP for AA:BB:CC:DD:EE:FF: %s", m["AA:BB:CC:DD:EE:FF"])
		}
	})

	t.Run("empty object", func(t *testing.T) {
		m, err := parseMACToIPJSON([]byte(`{}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(m) != 0 {
			t.Errorf("expected empty map, got %d entries", len(m))
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := parseMACToIPJSON([]byte(`not json`))
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		_, err := parseMACToIPJSON([]byte(``))
		if err == nil {
			t.Error("expected error for empty input")
		}
	})

	t.Run("invalid MAC key", func(t *testing.T) {
		_, err := parseMACToIPJSON([]byte(`{"not-a-mac":"10.0.0.1"}`))
		if err == nil {
			t.Error("expected error for invalid MAC key")
		}
	})

	t.Run("invalid IP value", func(t *testing.T) {
		_, err := parseMACToIPJSON([]byte(`{"AA:BB:CC:DD:EE:FF":"not-an-ip"}`))
		if err == nil {
			t.Error("expected error for invalid IP value")
		}
	})
}
