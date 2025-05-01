package response_inspector

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

// mockReadCloser implements io.ReadCloser for testing
type mockReadCloser struct {
	io.Reader
	closeErr error
	readErr  error
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	if m.Reader != nil {
		return m.Reader.Read(p)
	}
	return 0, io.EOF
}

func (m *mockReadCloser) Close() error {
	return m.closeErr
}

func TestInspectResponse(t *testing.T) {
	// err is not nil, httpRes is nil
	err := errors.New("some error")
	if got := InspectResponse(nil, err); got == nil || got.Error() != "some error" {
		t.Errorf("expected error to be returned when err is not nil and httpRes is nil")
	}

	// err is not nil, httpRes.StatusCode >= 400
	httpRes := &http.Response{
		Status:     "400 Bad Request",
		StatusCode: 400,
		Body:       io.NopCloser(bytes.NewBufferString("error body")),
	}
	err = errors.New("some error")
	got := InspectResponse(httpRes, err)
	if got == nil || got.Error() == "some error" {
		t.Errorf("expected formatted error with status and body, got: %v", got)
	}

	// err is nil, httpRes.StatusCode >= 400
	httpRes = &http.Response{
		Status:     "404 Not Found",
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewBufferString("not found")),
	}
	got = InspectResponse(httpRes, nil)
	if got == nil || got.Error() == "some error" {
		t.Errorf("expected formatted error with status and body, got: %v", got)
	}

	// err is nil, httpRes.StatusCode < 400
	httpRes = &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("ok")),
	}
	got = InspectResponse(httpRes, nil)
	if got != nil {
		t.Errorf("expected nil error for 200 OK, got: %v", got)
	}
}

func TestParseResponseBody(t *testing.T) {
	// httpRes is nil
	result, err := ParseResponseBody(nil)
	if err == nil || result != nil {
		t.Errorf("expected error when httpRes is nil")
	}

	// valid JSON body
	body := `{"key":"value"}`
	httpRes := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(body)),
	}
	result, err = ParseResponseBody(httpRes)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key to be 'value', got: %v", result["key"])
	}

	// invalid JSON body
	httpRes = &http.Response{
		Body: io.NopCloser(bytes.NewBufferString("not json")),
	}
	_, err = ParseResponseBody(httpRes)
	if err == nil {
		t.Errorf("expected error for invalid JSON")
	}

	// error reading body
	badBody := &mockReadCloser{
		readErr: errors.New("read error"),
	}
	httpRes = &http.Response{
		Body: badBody,
	}
	_, err = ParseResponseBody(httpRes)
	if err == nil {
		t.Errorf("expected error when reading body fails")
	}
}
