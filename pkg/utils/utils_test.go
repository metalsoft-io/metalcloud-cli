//go:build !integration

package utils

import (
	"bytes"
	"os"
	"testing"
)

func TestCreateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World!", "hello-world"},
		{"  Leading and trailing  ", "leading-and-trailing"},
		{"Special@#%Characters", "special-characters"},
		{"", ""},
		{"!!!", ""},
		{"Multiple   Spaces", "multiple-spaces"},
		{"Already-slug", "already-slug"},
		{"MiXeD CaSe", "mixed-case"},
	}

	for _, tt := range tests {
		got := CreateSlug(tt.input)
		if got != tt.expected {
			t.Errorf("CreateSlug(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetFloat32FromString(t *testing.T) {
	tests := []struct {
		input       string
		expected    float32
		expectError bool
	}{
		{"123.45", 123.45, false},
		{"0", 0, false},
		{"-42.1", -42.1, false},
		{"notanumber", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := GetFloat32FromString(tt.input)
		if tt.expectError && err == nil {
			t.Errorf("GetFloat32FromString(%q) expected error, got nil", tt.input)
		}
		if !tt.expectError && (err != nil || got != tt.expected) {
			t.Errorf("GetFloat32FromString(%q) = %v, %v; want %v, nil", tt.input, got, err, tt.expected)
		}
	}
}

func TestReadConfigFromFile(t *testing.T) {
	content := []byte("test config data")
	tmpfile, err := os.CreateTemp("", "testconfig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	read, err := ReadConfigFromFile(tmpfile.Name())
	if err != nil {
		t.Errorf("ReadConfigFromFile error: %v", err)
	}
	if !bytes.Equal(read, content) {
		t.Errorf("ReadConfigFromFile = %q; want %q", read, content)
	}

	_, err = ReadConfigFromFile("nonexistentfile")
	if err == nil {
		t.Error("ReadConfigFromFile with nonexistent file expected error, got nil")
	}
}

func TestReadConfigFromPipe(t *testing.T) {
	// Save original os.Stdin
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()

	// Simulate piped input
	r, w, _ := os.Pipe()
	testData := []byte("piped config data")
	go func() {
		w.Write(testData)
		w.Close()
	}()
	os.Stdin = r

	read, err := ReadConfigFromPipe()
	if err != nil {
		t.Errorf("ReadConfigFromPipe error: %v", err)
	}
	if !bytes.Equal(read, testData) {
		t.Errorf("ReadConfigFromPipe = %q; want %q", read, testData)
	}

	// Simulate no data piped (os.Stdin is a char device)
	// This is tricky to test in unit tests, so we skip this part.
}

func TestReadConfigFromPipeOrFile(t *testing.T) {
	// Test with empty string
	read, err := ReadConfigFromPipeOrFile("")
	if err != nil {
		t.Errorf("ReadConfigFromPipeOrFile(\"\") error: %v", err)
	}
	if read != nil {
		t.Errorf("ReadConfigFromPipeOrFile(\"\") = %v; want nil", read)
	}

	// Test with file
	content := []byte("file data")
	tmpfile, err := os.CreateTemp("", "testconfig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	read, err = ReadConfigFromPipeOrFile(tmpfile.Name())
	if err != nil {
		t.Errorf("ReadConfigFromPipeOrFile(file) error: %v", err)
	}
	if !bytes.Equal(read, content) {
		t.Errorf("ReadConfigFromPipeOrFile(file) = %q; want %q", read, content)
	}

	// Test with "pipe"
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()
	r, w, _ := os.Pipe()
	testData := []byte("pipe or file data")
	go func() {
		w.Write(testData)
		w.Close()
	}()
	os.Stdin = r

	read, err = ReadConfigFromPipeOrFile("pipe")
	if err != nil {
		t.Errorf("ReadConfigFromPipeOrFile(\"pipe\") error: %v", err)
	}
	if !bytes.Equal(read, testData) {
		t.Errorf("ReadConfigFromPipeOrFile(\"pipe\") = %q; want %q", read, testData)
	}
}
