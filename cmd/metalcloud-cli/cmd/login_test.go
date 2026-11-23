package cmd

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/login"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func useTempSessions(t *testing.T) {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".metalcloud", "sessions.json")
	old := login.SessionsPath
	login.SessionsPath = func() (string, error) { return path, nil }
	t.Cleanup(func() { login.SessionsPath = old })
}

// runCLIWithKey is runCLI with a chosen api_key (runCLI always sets "test-key").
func runCLIWithKey(t *testing.T, srv *httptest.Server, apiKey string, args ...string) (string, error) {
	t.Helper()
	resetFlags(rootCmd)
	viper.Set(system.ConfigEndpoint, srv.URL)
	viper.Set(system.ConfigApiKey, apiKey)
	viper.Set(formatter.ConfigFormat, "json")
	system.AllowDevelop = true
	var err error
	out := captureStdout(t, func() {
		rootCmd.SetArgs(args)
		err = rootCmd.Execute()
	})
	return out, err
}

func saveTestSession(t *testing.T, endpoint, token string, expires time.Time) {
	t.Helper()
	err := login.SaveSession(endpoint, login.Session{
		Token:            token,
		ExpiresTimestamp: expires.UTC().Format(time.RFC3339),
		UserEmail:        "me@x.io",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestResolveCredentialPrefersApiKey(t *testing.T) {
	useTempSessions(t)
	now := time.Now()
	saveTestSession(t, "https://m.example.com", "sess", now.Add(time.Hour))

	token, fromSession, err := resolveCredential("msat_key", "https://m.example.com", now)
	if err != nil || token != "msat_key" || fromSession {
		t.Errorf("got %q, %v, %v", token, fromSession, err)
	}

	token, fromSession, err = resolveCredential("", "https://m.example.com/", now)
	if err != nil || token != "sess" || !fromSession {
		t.Errorf("got %q, %v, %v", token, fromSession, err)
	}
}

func TestResolveCredentialSkipsExpiredSession(t *testing.T) {
	useTempSessions(t)
	now := time.Now()
	saveTestSession(t, "https://m.example.com", "sess", now.Add(-time.Minute))

	_, _, err := resolveCredential("", "https://m.example.com", now)
	if err == nil || err.Error() != noCredentialText {
		t.Errorf("err = %v", err)
	}
}

func TestResolveCredentialNoCredentialText(t *testing.T) {
	useTempSessions(t)
	_, _, err := resolveCredential("", "https://m.example.com", time.Now())
	if err == nil || err.Error() != noCredentialText {
		t.Errorf("err = %v, want %q", err, noCredentialText)
	}
}

func TestStartupDeletesSessionOn401AndKeepsItOn403(t *testing.T) {
	for _, tc := range []struct {
		status      int
		wantErr     string
		wantSession bool
	}{
		{http.StatusUnauthorized, sessionEndedText, false},
		{http.StatusForbidden, "This route does not accept a CLI session", true},
	} {
		useTempSessions(t)
		var auth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(tc.status)
			_, _ = w.Write([]byte(`{"message":"This route does not accept a CLI session; use the UI or an API token."}`))
		}))
		saveTestSession(t, srv.URL, "sess-token", time.Now().Add(time.Hour))

		_, err := runCLIWithKey(t, srv, "", "site", "list")
		srv.Close()

		if auth != "Bearer sess-token" {
			t.Errorf("status %d: Authorization = %q", tc.status, auth)
		}
		if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
			t.Errorf("status %d: err = %v, want it to contain %q", tc.status, err, tc.wantErr)
		}
		s, _ := login.LoadSession(srv.URL)
		if (s != nil) != tc.wantSession {
			t.Errorf("status %d: session present = %v, want %v", tc.status, s != nil, tc.wantSession)
		}
	}
}

func TestStartup401WithApiKeyKeepsSession(t *testing.T) {
	useTempSessions(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	saveTestSession(t, srv.URL, "sess-token", time.Now().Add(time.Hour))

	_, err := runCLIWithKey(t, srv, "msat_key", "site", "list")
	if err == nil || strings.Contains(err.Error(), sessionEndedText) {
		t.Errorf("err = %v", err)
	}
	if s, _ := login.LoadSession(srv.URL); s == nil {
		t.Error("session deleted although the API key was used")
	}
}

func TestLogoutEndsSessionAndForgetsIt(t *testing.T) {
	useTempSessions(t)
	var logoutAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/logout" && r.Method == http.MethodPost {
			logoutAuth = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}
		t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()
	saveTestSession(t, srv.URL, "sess-token", time.Now().Add(time.Hour))

	// no API key: logout must not demand one
	if _, err := runCLIWithKey(t, srv, "", "logout"); err != nil {
		t.Fatal(err)
	}
	if logoutAuth != "Bearer sess-token" {
		t.Errorf("Authorization = %q", logoutAuth)
	}
	if s, _ := login.LoadSession(srv.URL); s != nil {
		t.Error("session not deleted")
	}
}

func TestLogoutDeletesEntryEvenWhenServerCallFails(t *testing.T) {
	useTempSessions(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	saveTestSession(t, srv.URL, "sess-token", time.Now().Add(time.Hour))

	if _, err := runCLIWithKey(t, srv, "", "logout"); err != nil {
		t.Fatal(err)
	}
	if s, _ := login.LoadSession(srv.URL); s != nil {
		t.Error("session not deleted")
	}
}
