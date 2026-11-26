package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/login"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/cobra"
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
	c := login.Client{Endpoint: "https://m.example.com"}

	token, kind, err := resolveCredential(context.Background(), "msat_key", "", c, now)
	if err != nil || token != "msat_key" || kind != credentialAPIKey {
		t.Errorf("got %q, %v, %v", token, kind, err)
	}

	c.Endpoint = "https://m.example.com/"
	token, kind, err = resolveCredential(context.Background(), "", "", c, now)
	if err != nil || token != "sess" || kind != credentialSession {
		t.Errorf("got %q, %v, %v", token, kind, err)
	}
}

func TestResolveCredentialSkipsExpiredSession(t *testing.T) {
	useTempSessions(t)
	now := time.Now()
	saveTestSession(t, "https://m.example.com", "sess", now.Add(-time.Minute))

	_, _, err := resolveCredential(context.Background(), "", "", login.Client{Endpoint: "https://m.example.com"}, now)
	if err == nil || err.Error() != noCredentialText {
		t.Errorf("err = %v", err)
	}
}

func TestResolveCredentialNoCredentialText(t *testing.T) {
	useTempSessions(t)
	_, _, err := resolveCredential(context.Background(), "", "", login.Client{Endpoint: "https://m.example.com"}, time.Now())
	if err == nil || err.Error() != noCredentialText {
		t.Errorf("err = %v, want %q", err, noCredentialText)
	}
}

// signInServer answers the service-account sign-in route with a fixed session token and counts the calls.
func signInServer(t *testing.T, sessionToken string, calls *int32) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls != nil {
			atomic.AddInt32(calls, 1)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"token":"` + sessionToken + `","expiresTimestamp":"` + time.Now().Add(time.Hour).UTC().Format(time.RFC3339) + `"}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func tokenFileWith(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestResolveCredentialApiKeyBeatsTokenFile(t *testing.T) {
	useTempSessions(t)
	var calls int32
	srv := signInServer(t, "robot", &calls)

	token, kind, err := resolveCredential(context.Background(), "msat_key", tokenFileWith(t, "provider"), login.Client{Endpoint: srv.URL}, time.Now())
	if err != nil || token != "msat_key" || kind != credentialAPIKey || calls != 0 {
		t.Errorf("got %q, %v, calls %d, %v", token, kind, calls, err)
	}
}

func TestResolveCredentialTokenFileBeatsPersonSession(t *testing.T) {
	useTempSessions(t)
	var calls int32
	srv := signInServer(t, "robot", &calls)
	saveTestSession(t, srv.URL, "person", time.Now().Add(time.Hour))

	token, kind, err := resolveCredential(context.Background(), "", tokenFileWith(t, "provider"), login.Client{Endpoint: srv.URL}, time.Now())
	if err != nil || token != "robot" || kind != credentialServiceAccount || calls != 1 {
		t.Errorf("got %q, %v, calls %d, %v", token, kind, calls, err)
	}
	if s, _ := login.LoadSession(srv.URL); s == nil || s.Token != "person" {
		t.Errorf("the person's session was touched: %+v", s)
	}
}

func TestSessionEndedOrDeletesOnlyTheServiceAccountSession(t *testing.T) {
	useTempSessions(t)
	const endpoint = "https://m.example.com"
	const tokenFile = "/run/secrets/robot-token"
	saveTestSession(t, endpoint, "person", time.Now().Add(time.Hour))
	if err := login.SaveServiceAccountSession(endpoint, tokenFile, login.Session{Token: "robot"}); err != nil {
		t.Fatal(err)
	}

	err := sessionEndedOr(errors.New("401 Unauthorized"), credentialServiceAccount, endpoint, tokenFile)
	want := "The service account's session ended; the next command signs in again from the token file."
	if err == nil || err.Error() != want {
		t.Errorf("err = %v, want %q", err, want)
	}
	if s, _ := login.LoadServiceAccountSession(endpoint, tokenFile); s != nil {
		t.Errorf("the service-account session is still there: %+v", s)
	}
	if s, _ := login.LoadSession(endpoint); s == nil || s.Token != "person" {
		t.Errorf("the person's session was touched: %+v", s)
	}

	// a 403 keeps the session and passes through; a person's session keeps its own text
	if err := login.SaveServiceAccountSession(endpoint, tokenFile, login.Session{Token: "robot"}); err != nil {
		t.Fatal(err)
	}
	forbidden := errors.New("403 Forbidden")
	if got := sessionEndedOr(forbidden, credentialServiceAccount, endpoint, tokenFile); got != forbidden {
		t.Errorf("403: got %v", got)
	}
	if s, _ := login.LoadServiceAccountSession(endpoint, tokenFile); s == nil {
		t.Error("a 403 removed the service-account session")
	}
	if got := sessionEndedOr(errors.New("401 Unauthorized"), credentialSession, endpoint, tokenFile); got == nil || got.Error() != sessionEndedText {
		t.Errorf("person's session: got %v", got)
	}
	if s, _ := login.LoadServiceAccountSession(endpoint, tokenFile); s == nil {
		t.Error("a person's 401 removed the service-account session")
	}
}

// useTokenFile points the token-file setting at path for one test, through the environment variable.
func useTokenFile(t *testing.T, path string) {
	t.Helper()
	t.Setenv("METALCLOUD_SERVICE_ACCOUNT_TOKEN_FILE", path)
}

func TestServiceAccountTokenFileRunsCommandsAndSignsInAgain(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var apiAuth []string
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/service-accounts/actions/sign-in", func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&signIns, 1)
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["token"] != "provider-token" {
				t.Errorf("posted token = %q", body["token"])
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"token":"robot-session","expiresTimestamp":"` + time.Now().Add(time.Hour).UTC().Format(time.RFC3339) + `"}`))
		})
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			apiAuth = append(apiAuth, r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`))
		})
	})
	api := httptest.NewServer(mux)
	defer api.Close()
	saveTestSession(t, api.URL, "person-session", time.Now().Add(time.Hour))
	useTokenFile(t, tokenFileWith(t, "provider-token\n"))

	for i := 0; i < 2; i++ {
		if _, err := runCLIWithKey(t, api, "", "site", "list"); err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
	}
	if signIns != 1 {
		t.Errorf("sign-ins = %d, want 1 (the second command reuses the saved session)", signIns)
	}
	if len(apiAuth) != 2 || apiAuth[0] != "Bearer robot-session" || apiAuth[1] != "Bearer robot-session" {
		t.Errorf("Authorization headers = %v", apiAuth)
	}
	if s, _ := login.LoadSession(api.URL); s == nil || s.Token != "person-session" {
		t.Errorf("the person's session was touched: %+v", s)
	}
}

// perAccountServer signs a service account in as "session-for-<provider token>" and records the
// Authorization header of every /api/v2/sites call.
func perAccountServer(t *testing.T, signIns *int32, sitesAuth *[]string) *httptest.Server {
	t.Helper()
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/service-accounts/actions/sign-in", func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(signIns, 1)
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"token":"session-for-` + body["token"] + `","expiresTimestamp":"` + time.Now().Add(time.Hour).UTC().Format(time.RFC3339) + `"}`))
		})
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			*sitesAuth = append(*sitesAuth, r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`))
		})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestTwoServiceAccountTokenFilesKeepTheirOwnSessions(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	api := perAccountServer(t, &signIns, &auth)
	saveTestSession(t, api.URL, "person-session", time.Now().Add(time.Hour))
	fileA := tokenFileWith(t, "tok-a")
	fileB := tokenFileWith(t, "tok-b")

	for _, f := range []string{fileA, fileB, fileA, fileB} {
		useTokenFile(t, f)
		if _, err := runCLIWithKey(t, api, "", "site", "list"); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{"Bearer session-for-tok-a", "Bearer session-for-tok-b", "Bearer session-for-tok-a", "Bearer session-for-tok-b"}
	if strings.Join(auth, ",") != strings.Join(want, ",") {
		t.Errorf("Authorization headers = %v, want %v", auth, want)
	}
	if signIns != 2 {
		t.Errorf("sign-ins = %d, want 2 (one per token file)", signIns)
	}
	if s, _ := login.LoadSession(api.URL); s == nil || s.Token != "person-session" {
		t.Errorf("the person's session was touched: %+v", s)
	}
}

func TestLoginWithAnotherTokenFileLeavesTheFirstAccountsSession(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	api := perAccountServer(t, &signIns, &auth)
	fileA := tokenFileWith(t, "tok-a")
	fileB := tokenFileWith(t, "tok-b")

	for _, f := range []string{fileA, fileB} {
		if _, err := runCLIWithKey(t, api, "", "login", "--service-account-token-file", f); err != nil {
			t.Fatal(err)
		}
	}
	useTokenFile(t, fileA)
	if _, err := runCLIWithKey(t, api, "", "site", "list"); err != nil {
		t.Fatal(err)
	}
	if len(auth) != 1 || auth[0] != "Bearer session-for-tok-a" || signIns != 2 {
		t.Errorf("Authorization headers = %v, sign-ins = %d", auth, signIns)
	}
}

// endedSessionServer answers the startup GET /api/v2/user with 401 for the tokens in ended, signs
// service accounts in as "fresh-N" and counts the sign-ins.
func endedSessionServer(t *testing.T, ended map[string]bool, signInStatus int, signIns *int32, sitesAuth *[]string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if ended[strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")] {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"Unauthorized"}`))
			return
		}
		_, _ = w.Write([]byte(fullUserJSON))
	})
	mux.HandleFunc("/api/v2/service-accounts/actions/sign-in", func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(signIns, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(signInStatus)
		if signInStatus == http.StatusCreated {
			_, _ = w.Write([]byte(`{"token":"fresh-` + string(rune('0'+n)) + `","expiresTimestamp":"` + time.Now().Add(time.Hour).UTC().Format(time.RFC3339) + `"}`))
		}
	})
	mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
		*sitesAuth = append(*sitesAuth, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func saveServiceAccountTestSession(t *testing.T, endpoint, tokenFile, token string) {
	t.Helper()
	err := login.SaveServiceAccountSession(endpoint, tokenFile, login.Session{
		Token:            token,
		ExpiresTimestamp: time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartup401OnAServiceAccountSessionSignsInOnceAndRuns(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	api := endedSessionServer(t, map[string]bool{"stale": true}, http.StatusCreated, &signIns, &auth)
	file := tokenFileWith(t, "provider-token")
	saveServiceAccountTestSession(t, api.URL, file, "stale")
	useTokenFile(t, file)

	if _, err := runCLIWithKey(t, api, "", "site", "list"); err != nil {
		t.Fatalf("err = %v", err)
	}
	if signIns != 1 {
		t.Errorf("sign-ins = %d, want 1", signIns)
	}
	if len(auth) != 1 || auth[0] != "Bearer fresh-1" {
		t.Errorf("the command ran with %v, want the new token", auth)
	}
	if s, _ := login.LoadServiceAccountSession(api.URL, file); s == nil || s.Token != "fresh-1" {
		t.Errorf("saved = %+v, want the new session", s)
	}
}

func TestStartup401AfterTheFreshSignInKeepsTheEndedText(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	api := endedSessionServer(t, map[string]bool{"stale": true, "fresh-1": true, "fresh-2": true}, http.StatusCreated, &signIns, &auth)
	file := tokenFileWith(t, "provider-token")
	saveServiceAccountTestSession(t, api.URL, file, "stale")
	useTokenFile(t, file)

	_, err := runCLIWithKey(t, api, "", "site", "list")
	if err == nil || err.Error() != serviceAccountSessionEndedText {
		t.Errorf("err = %v, want %q", err, serviceAccountSessionEndedText)
	}
	if signIns != 1 {
		t.Errorf("sign-ins = %d, want exactly 1 (no loop)", signIns)
	}
	if len(auth) != 0 {
		t.Errorf("the command ran: %v", auth)
	}
}

func TestStartup401OnAServiceAccountSessionFailedSignInReturnsTheSignInText(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	api := endedSessionServer(t, map[string]bool{"stale": true}, http.StatusUnauthorized, &signIns, &auth)
	file := tokenFileWith(t, "provider-token")
	saveServiceAccountTestSession(t, api.URL, file, "stale")
	useTokenFile(t, file)

	_, err := runCLIWithKey(t, api, "", "site", "list")
	want := "The provider token wasn't accepted; check the token file and the service account's subject."
	if err == nil || err.Error() != want {
		t.Errorf("err = %v, want %q", err, want)
	}
	if signIns != 1 {
		t.Errorf("sign-ins = %d, want 1", signIns)
	}
}

func TestStartup401OnAPersonsSessionDoesNotSignIn(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	api := endedSessionServer(t, map[string]bool{"person-session": true}, http.StatusCreated, &signIns, &auth)
	saveTestSession(t, api.URL, "person-session", time.Now().Add(time.Hour))

	_, err := runCLIWithKey(t, api, "", "site", "list")
	if err == nil || err.Error() != sessionEndedText {
		t.Errorf("err = %v, want %q", err, sessionEndedText)
	}
	if signIns != 0 {
		t.Errorf("sign-ins = %d, want 0", signIns)
	}
	if s, _ := login.LoadSession(api.URL); s != nil {
		t.Errorf("the person's ended session is still saved: %+v", s)
	}
}

// midCommandServer answers the startup calls with 200, signs service accounts in as "fresh-N" and
// answers /api/v2/sites with the status in sitesStatus (200 gives an empty list), recording the
// Authorization header of each call.
func midCommandServer(t *testing.T, sitesStatus *int32, signIns *int32, sitesAuth *[]string) *httptest.Server {
	t.Helper()
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/service-accounts/actions/sign-in", func(w http.ResponseWriter, r *http.Request) {
			n := atomic.AddInt32(signIns, 1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"token":"fresh-` + string(rune('0'+n)) + `","expiresTimestamp":"` + time.Now().Add(time.Hour).UTC().Format(time.RFC3339) + `"}`))
		})
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			*sitesAuth = append(*sitesAuth, r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			status := int(atomic.LoadInt32(sitesStatus))
			if status != http.StatusOK {
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"message":"refused"}`))
				return
			}
			_, _ = w.Write([]byte(`{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`))
		})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestMidCommand401OnAServiceAccountSessionEndsIt(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	sitesStatus := int32(http.StatusUnauthorized)
	api := midCommandServer(t, &sitesStatus, &signIns, &auth)
	file := tokenFileWith(t, "provider-token")
	saveServiceAccountTestSession(t, api.URL, file, "robot-session")
	saveTestSession(t, api.URL, "person-session", time.Now().Add(time.Hour))
	useTokenFile(t, file)

	_, err := runCLIWithKey(t, api, "", "site", "list")
	if err == nil || err.Error() != serviceAccountSessionEndedText {
		t.Fatalf("err = %v, want %q", err, serviceAccountSessionEndedText)
	}
	if len(auth) != 1 || auth[0] != "Bearer robot-session" {
		t.Errorf("the command called the API with %v, want one call with the saved session (no retry)", auth)
	}
	if signIns != 0 {
		t.Errorf("sign-ins = %d, want 0 (no retry inside the command)", signIns)
	}
	if s, _ := login.LoadServiceAccountSession(api.URL, file); s != nil {
		t.Errorf("the ended service account session is still saved: %+v", s)
	}
	if s, _ := login.LoadSession(api.URL); s == nil || s.Token != "person-session" {
		t.Errorf("the person's session was touched: %+v", s)
	}
}

func TestMidCommand401ThenTheNextCommandSignsInOnceFromTheFile(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	sitesStatus := int32(http.StatusUnauthorized)
	api := midCommandServer(t, &sitesStatus, &signIns, &auth)
	file := tokenFileWith(t, "provider-token")
	saveServiceAccountTestSession(t, api.URL, file, "robot-session")
	useTokenFile(t, file)

	_, _ = runCLIWithKey(t, api, "", "site", "list")
	atomic.StoreInt32(&sitesStatus, http.StatusOK)
	auth = nil
	if _, err := runCLIWithKey(t, api, "", "site", "list"); err != nil {
		t.Fatalf("next command: %v", err)
	}
	if signIns != 1 {
		t.Errorf("sign-ins = %d, want 1", signIns)
	}
	if len(auth) != 1 || auth[0] != "Bearer fresh-1" {
		t.Errorf("the next command called the API with %v, want the new session", auth)
	}
	if s, _ := login.LoadServiceAccountSession(api.URL, file); s == nil || s.Token != "fresh-1" {
		t.Errorf("saved = %+v, want the new session", s)
	}
}

func TestMidCommand401OnAPersonsSessionKeepsTheRawErrorAndTheSession(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	sitesStatus := int32(http.StatusUnauthorized)
	api := midCommandServer(t, &sitesStatus, &signIns, &auth)
	saveTestSession(t, api.URL, "person-session", time.Now().Add(time.Hour))

	_, err := runCLIWithKey(t, api, "", "site", "list")
	if err == nil || !strings.HasPrefix(err.Error(), "401 ") {
		t.Errorf("err = %v, want the raw 401 error", err)
	}
	if err != nil && (err.Error() == sessionEndedText || err.Error() == serviceAccountSessionEndedText) {
		t.Errorf("err = %q, want the unchanged raw error", err)
	}
	if signIns != 0 {
		t.Errorf("sign-ins = %d, want 0", signIns)
	}
	if s, _ := login.LoadSession(api.URL); s == nil || s.Token != "person-session" {
		t.Errorf("the person's session changed: %+v", s)
	}
}

func TestMidCommand403OnAServiceAccountSessionIsUnchanged(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	sitesStatus := int32(http.StatusForbidden)
	api := midCommandServer(t, &sitesStatus, &signIns, &auth)
	file := tokenFileWith(t, "provider-token")
	saveServiceAccountTestSession(t, api.URL, file, "robot-session")
	useTokenFile(t, file)

	_, err := runCLIWithKey(t, api, "", "site", "list")
	if err == nil || !strings.HasPrefix(err.Error(), "403 ") {
		t.Errorf("err = %v, want the raw 403 error", err)
	}
	if signIns != 0 {
		t.Errorf("sign-ins = %d, want 0", signIns)
	}
	if s, _ := login.LoadServiceAccountSession(api.URL, file); s == nil || s.Token != "robot-session" {
		t.Errorf("the service account session changed: %+v", s)
	}
}

// addFailingCommand registers a command under the root for one test; its RunE returns err.
func addFailingCommand(t *testing.T, err error) {
	t.Helper()
	c := &cobra.Command{Use: "mid-command-fail", RunE: func(*cobra.Command, []string) error { return err }}
	rootCmd.AddCommand(c)
	t.Cleanup(func() { rootCmd.RemoveCommand(c) })
}

func TestMidCommandWrapped401OnAServiceAccountSessionEndsIt(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"wrapped with %w", fmt.Errorf("failed to list sites: %w", errors.New("401 Unauthorized - {\"message\":\"x\"}"))},
		{"wrapped with %v", fmt.Errorf("failed to list sites: %v", errors.New("401 Unauthorized - {\"message\":\"x\"}"))},
		{"wrapped twice", fmt.Errorf("a: %w", fmt.Errorf("b: %v", errors.New("401 Unauthorized")))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			useTempSessions(t)
			var signIns int32
			var auth []string
			sitesStatus := int32(http.StatusOK)
			api := midCommandServer(t, &sitesStatus, &signIns, &auth)
			file := tokenFileWith(t, "provider-token")
			saveServiceAccountTestSession(t, api.URL, file, "robot-session")
			useTokenFile(t, file)
			addFailingCommand(t, tc.err)

			_, err := runCLIWithKey(t, api, "", "mid-command-fail")
			if err == nil || err.Error() != serviceAccountSessionEndedText {
				t.Errorf("err = %v, want %q", err, serviceAccountSessionEndedText)
			}
			if s, _ := login.LoadServiceAccountSession(api.URL, file); s != nil {
				t.Errorf("the ended session is still saved: %+v", s)
			}
			if signIns != 0 {
				t.Errorf("sign-ins = %d, want 0", signIns)
			}
		})
	}
}

func TestMidCommandErrorThatOnlyContains401ElsewhereIsUnchanged(t *testing.T) {
	for _, msg := range []string{
		"failed to get item 401 : 500 Internal Server Error",
		"failed to get item 401 Unauthorized-looking-name: 500 Internal Server Error",
		"500 Internal Server Error - request 4401 failed",
	} {
		t.Run(msg, func(t *testing.T) {
			useTempSessions(t)
			var signIns int32
			var auth []string
			sitesStatus := int32(http.StatusOK)
			api := midCommandServer(t, &sitesStatus, &signIns, &auth)
			file := tokenFileWith(t, "provider-token")
			saveServiceAccountTestSession(t, api.URL, file, "robot-session")
			useTokenFile(t, file)
			addFailingCommand(t, errors.New(msg))

			_, err := runCLIWithKey(t, api, "", "mid-command-fail")
			if err == nil || err.Error() != msg {
				t.Errorf("err = %v, want %q unchanged", err, msg)
			}
			if s, _ := login.LoadServiceAccountSession(api.URL, file); s == nil || s.Token != "robot-session" {
				t.Errorf("the session changed: %+v", s)
			}
		})
	}
}

func TestMidCommandWrapperDoesNotLeakWhenRunEIsNotReached(t *testing.T) {
	useTempSessions(t)
	var signIns int32
	var auth []string
	sitesStatus := int32(http.StatusOK)
	api := midCommandServer(t, &sitesStatus, &signIns, &auth)
	file := tokenFileWith(t, "provider-token")
	saveServiceAccountTestSession(t, api.URL, file, "robot-session")
	useTokenFile(t, file)

	ran := false
	c := &cobra.Command{
		Use:     "mid-command-prefail",
		PreRunE: func(*cobra.Command, []string) error { return errors.New("pre-run failed") },
		RunE:    func(*cobra.Command, []string) error { ran = true; return nil },
	}
	rootCmd.AddCommand(c)
	t.Cleanup(func() { rootCmd.RemoveCommand(c) })
	original := reflect.ValueOf(c.RunE).Pointer()

	if _, err := runCLIWithKey(t, api, "", "mid-command-prefail"); err == nil || err.Error() != "pre-run failed" {
		t.Fatalf("err = %v, want the pre-run error", err)
	}
	if ran {
		t.Fatal("RunE ran after a failing PreRunE")
	}

	// cobra has no hook after a failed run, so the next run puts the original RunE back before anything
	// else happens; the wrapper never survives into a later execution.
	if _, err := runCLIWithKey(t, api, "test-key", "site", "list"); err != nil {
		t.Fatalf("next run: %v", err)
	}
	if got := reflect.ValueOf(c.RunE).Pointer(); got != original {
		t.Errorf("the command's RunE is still the wrapper after the next run")
	}
}

func TestIsUnauthorizedMatchesTheStatusOnlyAtTheStartOrAfterAColon(t *testing.T) {
	for msg, want := range map[string]bool{
		"401 Unauthorized - {}":                      true,
		"failed to x: 401 Unauthorized - {}":         true,
		"a: b: 401 Unauthorized":                     true,
		"failed to get item 401 : 500 Internal":      false,
		"failed to get item 401 Unauthorized: 500 x": false,
		"403 Forbidden - 401 not here":               false,
		"4010 Unauthorized":                          false,
		"":                                           false,
	} {
		if got := isUnauthorized(errors.New(msg)); got != want {
			t.Errorf("isUnauthorized(%q) = %v, want %v", msg, got, want)
		}
	}
}

func TestLoginWithTokenFileSignsInAndSavesAMarkedSession(t *testing.T) {
	useTempSessions(t)
	var posted string
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/service-accounts/actions/sign-in" {
			t.Errorf("unexpected %s", r.URL.Path)
		}
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		posted = body["token"]
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"token":"robot-session","expiresTimestamp":"2030-01-01T00:00:00Z"}`))
	}))
	defer api.Close()
	file := tokenFileWith(t, "provider-token\n")

	// no API key and no saved session: login must not demand either
	out, err := runCLIWithKey(t, api, "", "login", "--service-account-token-file", file)
	if err != nil {
		t.Fatal(err)
	}
	if posted != "provider-token" {
		t.Errorf("posted token = %q", posted)
	}
	if s, _ := login.LoadServiceAccountSession(api.URL, file); s == nil || s.Token != "robot-session" || !s.ServiceAccount {
		t.Errorf("saved = %+v", s)
	}
	if s, _ := login.LoadSession(api.URL); s != nil {
		t.Errorf("a person's session was written: %+v", s)
	}
	if strings.Contains(out, "robot-session") || strings.Contains(out, "provider-token") {
		t.Errorf("the output carries a token: %q", out)
	}
	if !strings.Contains(out, "Signed in as a service account") {
		t.Errorf("output = %q", out)
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

func TestApiKeyWhitespaceIsTrimmedAndControlCharactersAreRefused(t *testing.T) {
	useTempSessions(t)
	var requests int32
	var auth []string
	mux := newMux(nil, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			auth = append(auth, r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`))
		})
	})
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()

	// a token pasted or read from a file often ends in a newline
	for _, key := range []string{"msat_key\n", "msat_key\r\n", " msat_key \n", "\tmsat_key"} {
		auth = nil
		if _, err := runCLIWithKey(t, api, key, "site", "list"); err != nil {
			t.Errorf("%q: %v", key, err)
		}
		if len(auth) != 1 || auth[0] != "Bearer msat_key" {
			t.Errorf("%q: Authorization headers = %q", key, auth)
		}
	}

	// whitespace alone is no token: the saved session is looked for, and there is none
	if _, err := runCLIWithKey(t, api, " \n", "site", "list"); err == nil || err.Error() != noCredentialText {
		t.Errorf("whitespace only: err = %v, want %q", err, noCredentialText)
	}

	// a control character inside is refused by name, before any request, without showing the token
	for _, tc := range []struct{ key, char string }{
		{"msat_\nkey", `'\n'`},
		{"msat_key\x00tail", `'\x00'`},
		{"msat\x7fkey", `'\x7f'`},
	} {
		atomic.StoreInt32(&requests, 0)
		_, err := runCLIWithKey(t, api, tc.key, "site", "list")
		if err == nil || !strings.Contains(err.Error(), "API token") || !strings.Contains(err.Error(), tc.char) {
			t.Errorf("%q: err = %v, want one naming the API token and %s", tc.key, err, tc.char)
		} else if strings.Contains(err.Error(), "msat") || strings.Contains(err.Error(), "key\x00") {
			t.Errorf("%q: the error shows the token: %v", tc.key, err)
		}
		if n := atomic.LoadInt32(&requests); n != 0 {
			t.Errorf("%q: %d requests sent", tc.key, n)
		}
	}
}
