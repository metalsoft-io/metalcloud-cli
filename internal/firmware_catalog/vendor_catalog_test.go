package firmware_catalog

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestDownloadBinary(t *testing.T) {
	// Create a test server to simulate binary downloads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test binary content"))
	}))
	defer server.Close()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "binaries-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("MissingVendorURL", func(t *testing.T) {
		vc := &VendorCatalog{}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: "",
		}

		_, err := vc.downloadBinary(binary)

		if err == nil {
			t.Error("Expected error for missing vendor URL, got nil")
		}
		if !strings.Contains(err.Error(), "no vendor download URL") {
			t.Errorf("Expected error about missing URL, got: %v", err)
		}
	})

	t.Run("HTTPError", func(t *testing.T) {
		vc := &VendorCatalog{}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: server.URL + "/error",
		}

		_, err := vc.downloadBinary(binary)

		if err == nil {
			t.Error("Expected HTTP error, got nil")
		}
		if !strings.Contains(err.Error(), "non-OK response") {
			t.Errorf("Expected error about non-OK response, got: %v", err)
		}
	})

	t.Run("InvalidURL", func(t *testing.T) {
		vc := &VendorCatalog{}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: "http://invalid-url-that-should-not-resolve.example",
		}

		_, err := vc.downloadBinary(binary)

		if err == nil {
			t.Error("Expected error for invalid URL, got nil")
		}
	})

	t.Run("DownloadToTempFile", func(t *testing.T) {
		vc := &VendorCatalog{}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: server.URL,
		}

		localPath, err := vc.downloadBinary(binary)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Check result - when using temporary files
		if !strings.HasSuffix(localPath, ".bin") {
			t.Errorf("Expected local path ending on '.bin', got: %s", localPath)
		}
	})

	t.Run("DownloadToLocalPath", func(t *testing.T) {
		vc := &VendorCatalog{
			VendorLocalBinariesPath: tempDir,
		}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: server.URL,
		}

		localPath, err := vc.downloadBinary(binary)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Check if file exists in the specified local path
		expectedPath := filepath.Join(tempDir, "test-binary")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected file to exist at %s, but it doesn't", expectedPath)
		}

		// Verify file content
		content, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}
		if string(content) != "test binary content" {
			t.Errorf("Expected content 'test binary content', got: %s", string(content))
		}

		// Check local path
		if !strings.HasSuffix(localPath, "test-binary") {
			t.Errorf("Expected local path to end with 'test-binary', got: %s", localPath)
		}
	})

	t.Run("WithRepoHttpUrl", func(t *testing.T) {
		vc := &VendorCatalog{
			VendorLocalBinariesPath: tempDir,
			RepoBaseUrl:             "http://firmware-repo.example.com",
		}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: server.URL,
		}

		localPath, err := vc.downloadBinary(binary)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Check if local path is correct
		if !strings.HasSuffix(localPath, "test-binary") {
			t.Errorf("Expected local path ending on 'test-binary', got: %s", localPath)
		}
	})

	t.Run("FileCannotBeCreated", func(t *testing.T) {
		// Create and immediately remove a directory to make it unavailable
		nonExistentDir := filepath.Join(tempDir, "non-existent-dir")
		os.Mkdir(nonExistentDir, 0755)
		os.Remove(nonExistentDir)

		vc := &VendorCatalog{
			VendorLocalBinariesPath: nonExistentDir,
		}
		binary := &sdk.FirmwareBinary{
			ExternalId:        sdk.PtrString("test-binary"),
			VendorDownloadUrl: server.URL,
		}

		_, err := vc.downloadBinary(binary)

		if err == nil {
			t.Error("Expected file creation error, got nil")
		}
	})
}

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
