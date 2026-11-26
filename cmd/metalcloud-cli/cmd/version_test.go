package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

// newVersionAuthServer answers /api/v2/user only for "Bearer good", and with the gateway's 401 body
// for an expired API token otherwise.
func newVersionAuthServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(versionResponse))
	})
	mux.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer good" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"statusCode":401,"errorCode":"401","message":"API token has expired."}`))
			return
		}
		_, _ = w.Write([]byte(fullUserJSON))
	})
	return httptest.NewServer(mux)
}

func userLine(out string) string {
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "User") {
			return line
		}
	}
	return ""
}

func TestVersionSaysWhyTheUserCannotBeRead(t *testing.T) {
	useTempSessions(t)
	srv := newVersionAuthServer()
	defer srv.Close()

	// a refused token: the line says so and why, and version still succeeds
	out, err := runCLIWithKey(t, srv, "expired", "version")
	if err != nil {
		t.Fatalf("expired: %v", err)
	}
	if line := userLine(out); !strings.HasPrefix(line, "User: cannot be read: ") || !strings.Contains(line, "API token has expired.") {
		t.Errorf("expired: user line = %q", line)
	}

	out, err = runCLIWithKey(t, srv, "good", "version")
	if line := userLine(out); err != nil || !strings.HasPrefix(line, "User ID: 1, Email: test@test.com") {
		t.Errorf("good: user line = %q, err = %v", line, err)
	}

	// no credential at all
	out, err = runCLIWithKey(t, srv, "", "version")
	if line := userLine(out); err != nil || line != "User: cannot be read: "+noCredentialText {
		t.Errorf("none: user line = %q, err = %v", line, err)
	}
}

func TestVersionUsesTheSavedSession(t *testing.T) {
	useTempSessions(t)
	srv := newVersionAuthServer()
	defer srv.Close()
	saveTestSession(t, srv.URL, "good", time.Now().Add(time.Hour))

	out, err := runCLIWithKey(t, srv, "", "version")
	if line := userLine(out); err != nil || !strings.HasPrefix(line, "User ID: 1, Email: test@test.com") {
		t.Errorf("user line = %q, err = %v", line, err)
	}
}
