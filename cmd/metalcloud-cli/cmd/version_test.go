package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// version_test.go covers:
//   version — exits 0, output contains version info

var versionResponse = `{"version":"6.0.0","minCliVersion":"1.0.0","maxCliVersion":"99.99.99"}`

func newVersionTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/version", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(versionResponse))
		})
	})
	return httptest.NewServer(mux)
}

func TestVersionCommand(t *testing.T) {
	srv := newVersionTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "CLI Version") {
		t.Errorf("expected output to contain 'CLI Version', got: %s", out)
	}
}

func TestVersionCommandAlias(t *testing.T) {
	srv := newVersionTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "ver")
	if err != nil {
		t.Fatalf("unexpected error running 'ver' alias: %v", err)
	}
	if !strings.Contains(out, "CLI Version") {
		t.Errorf("expected output to contain 'CLI Version', got: %s", out)
	}
}
