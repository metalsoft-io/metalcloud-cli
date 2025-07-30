package firmware_catalog

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// TestDownloadBinaryWithMockResponses tests downloadBinary with custom HTTP client mocks
func TestDownloadBinaryWithMockResponses(t *testing.T) {
	// Store the original http.DefaultClient to restore it later
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "binaries-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("MockHTTPFailure", func(t *testing.T) {
		// Mock client that always fails
		http.DefaultClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return nil, &mockNetworkError{message: "simulated network failure"}
				},
			},
		}

		vc := &VendorCatalog{
			VendorLocalBinariesPath: tempDir,
		}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: "http://example.com/binary",
		}

		_, err := vc.downloadBinary(binary)

		if err == nil {
			t.Error("Expected network error, got nil")
		}
		if !strings.Contains(err.Error(), "simulated network failure") {
			t.Errorf("Expected network error, got: %v", err)
		}
	})

	t.Run("MockHTTPSuccess", func(t *testing.T) {
		// Mock client that returns success with specific content
		http.DefaultClient = &http.Client{
			Transport: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("mocked binary data")),
					}, nil
				},
			},
		}

		vc := &VendorCatalog{
			VendorLocalBinariesPath: tempDir,
		}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: "http://example.com/binary",
		}

		_, err := vc.downloadBinary(binary)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify file content
		expectedPath := filepath.Join(tempDir, "test-binary")
		content, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}
		if string(content) != "mocked binary data" {
			t.Errorf("Expected content 'mocked binary data', got: %s", string(content))
		}
	})
}

// Helper types for mocking
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTripFunc(req)
}

type mockNetworkError struct {
	message string
}

func (e *mockNetworkError) Error() string {
	return e.message
}

func (e *mockNetworkError) Timeout() bool {
	return false
}

func (e *mockNetworkError) Temporary() bool {
	return true
}
