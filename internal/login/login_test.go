package login

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func fastPolling(t *testing.T) {
	t.Helper()
	old := slowDownStep
	slowDownStep = 20 * time.Millisecond
	t.Cleanup(func() { slowDownStep = old })
}

func TestStartNotSupportedOn401And404(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusNotFound} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		}))
		_, err := Start(context.Background(), Client{Endpoint: srv.URL}, "host", "1.0")
		srv.Close()
		if !errors.Is(err, ErrNotSupported) {
			t.Errorf("status %d: err = %v, want ErrNotSupported", status, err)
		}
	}
	if ErrNotSupported.Error() != "This MetalCloud doesn't support signing in from the CLI; use an API token." {
		t.Errorf("ErrNotSupported text = %q", ErrNotSupported.Error())
	}
}

func TestStartSendsClientAndParsesResponse(t *testing.T) {
	var got map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2/cli-login" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"deviceCode":"dev","userCode":"ABCD-EFGH","verificationUri":"https://ui/cli-login","verificationUriComplete":"https://ui/cli-login?code=ABCD-EFGH","expiresIn":600,"interval":5}`))
	}))
	defer srv.Close()

	// a trailing slash on the endpoint must not double up
	s, err := Start(context.Background(), Client{Endpoint: srv.URL + "/"}, "host", "1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if got["clientName"] != "host" || got["clientVersion"] != "1.2.3" {
		t.Errorf("body = %v", got)
	}
	if s.DeviceCode != "dev" || s.UserCode != "ABCD-EFGH" || s.VerificationURIComplete != "https://ui/cli-login?code=ABCD-EFGH" ||
		s.Interval != 5*time.Second || s.ExpiresIn != 600*time.Second {
		t.Errorf("started = %+v", s)
	}
}

func startedFor(interval time.Duration) *Started {
	return &Started{DeviceCode: "dev", Interval: interval}
}

func TestWaitPendingThenApproved(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/cli-login/token" {
			t.Errorf("path = %s", r.URL.Path)
		}
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["deviceCode"] != "dev" {
			t.Errorf("deviceCode = %q", body["deviceCode"])
		}
		if atomic.AddInt32(&calls, 1) < 3 {
			_, _ = w.Write([]byte(`{"status":"pending"}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"approved","token":"tok","expiresTimestamp":"2026-12-01T10:00:00.000Z","user":{"id":7,"email":"me@x.io"}}`))
	}))
	defer srv.Close()

	a, err := Wait(context.Background(), Client{Endpoint: srv.URL}, startedFor(10*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if a.Token != "tok" || a.ExpiresTimestamp != "2026-12-01T10:00:00.000Z" || a.UserEmail != "me@x.io" || a.UserID != 7 {
		t.Errorf("approved = %+v", a)
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestWaitSlowsDownOn429(t *testing.T) {
	fastPolling(t)
	var times []time.Time
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		times = append(times, time.Now())
		if len(times) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte(`{"status":"approved","token":"t","expiresTimestamp":"2026-12-01T10:00:00Z","user":{"id":1,"email":"a@b"}}`))
	}))
	defer srv.Close()

	// interval 10 ms, 429 adds slowDownStep (20 ms): the second poll comes >= 30 ms after the first
	if _, err := Wait(context.Background(), Client{Endpoint: srv.URL}, startedFor(10*time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	if len(times) != 2 {
		t.Fatalf("polls = %d, want 2", len(times))
	}
	if gap := times[1].Sub(times[0]); gap < 30*time.Millisecond {
		t.Errorf("gap after 429 = %v, want >= 30ms", gap)
	}
	if slowDownStep != 20*time.Millisecond {
		t.Error("slowDownStep changed")
	}
}

func TestWaitDeniedAndExpired(t *testing.T) {
	for status, want := range map[int]error{http.StatusForbidden: ErrDenied, http.StatusGone: ErrExpired} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		}))
		_, err := Wait(context.Background(), Client{Endpoint: srv.URL}, startedFor(5*time.Millisecond))
		srv.Close()
		if !errors.Is(err, want) {
			t.Errorf("status %d: err = %v, want %v", status, err, want)
		}
	}
}

func TestWaitStopsOnContextCancel(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{"status":"pending"}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := Wait(ctx, Client{Endpoint: srv.URL}, startedFor(10*time.Millisecond))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context deadline", err)
	}
	if time.Since(start) > time.Second {
		t.Error("Wait did not stop promptly")
	}
	n := atomic.LoadInt32(&calls)
	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&calls) != n {
		t.Error("polling continued after cancel")
	}
}

func TestLogoutSendsBearerAndPost(t *testing.T) {
	var method, path, auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, auth = r.Method, r.URL.Path, r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := Logout(context.Background(), Client{Endpoint: srv.URL}, "sess-token"); err != nil {
		t.Fatal(err)
	}
	if method != http.MethodPost || path != "/api/v2/logout" || auth != "Bearer sess-token" {
		t.Errorf("got %s %s auth=%q", method, path, auth)
	}

	srv401 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv401.Close()
	if err := Logout(context.Background(), Client{Endpoint: srv401.URL}, "x"); err == nil {
		t.Error("expected an error on 401")
	}
}

func TestLoginWithNoBrowserDoesNotOpen(t *testing.T) {
	useTempSessionsPath(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/cli-login":
			_, _ = w.Write([]byte(`{"deviceCode":"dev","userCode":"ABCDEFGH","verificationUri":"https://ui/cli-login","verificationUriComplete":"https://ui/cli-login?code=ABCD-EFGH","expiresIn":600,"interval":0}`))
		case "/api/v2/cli-login/token":
			_, _ = w.Write([]byte(`{"status":"approved","token":"tok","expiresTimestamp":"2026-12-01T10:00:00.000Z","user":{"id":7,"email":"me@x.io"}}`))
		}
	}))
	defer srv.Close()

	var opened []string
	old := openBrowser
	openBrowser = func(url string) error { opened = append(opened, url); return nil }
	t.Cleanup(func() { openBrowser = old })

	var out bytes.Buffer
	if err := Login(context.Background(), Client{Endpoint: srv.URL}, "host", "1.0", true, &out); err != nil {
		t.Fatal(err)
	}
	if len(opened) != 0 {
		t.Errorf("browser opened with --no-browser: %v", opened)
	}
	for _, want := range []string{"https://ui/cli-login?code=ABCD-EFGH", "ABCD-EFGH", "Signed in as me@x.io", "2026-12-01 10:00 UTC"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "tok") {
		t.Errorf("output shows the token:\n%s", out.String())
	}
	s, err := LoadSession(srv.URL)
	if err != nil || s == nil || s.Token != "tok" || s.UserEmail != "me@x.io" || s.CreatedTimestamp == "" {
		t.Errorf("saved session = %+v, err = %v", s, err)
	}

	// and without --no-browser the complete URL is opened, a failure to open is not fatal
	openBrowser = func(url string) error { opened = append(opened, url); return errors.New("no browser") }
	out.Reset()
	if err := Login(context.Background(), Client{Endpoint: srv.URL}, "host", "1.0", false, &out); err != nil {
		t.Fatal(err)
	}
	if len(opened) != 1 || opened[0] != "https://ui/cli-login?code=ABCD-EFGH" {
		t.Errorf("opened = %v", opened)
	}
}

// serviceAccountServer answers the sign-in route and counts the calls. The provider tokens it saw
// are collected in tokens.
type serviceAccountServer struct {
	*httptest.Server
	calls  int32
	tokens []string
}

func newServiceAccountServer(t *testing.T, status int, sessionToken, expires string) *serviceAccountServer {
	t.Helper()
	s := &serviceAccountServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&s.calls, 1)
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2/service-accounts/actions/sign-in" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "" {
			t.Errorf("sign-in sent an Authorization header: %q", r.Header.Get("Authorization"))
		}
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		s.tokens = append(s.tokens, body["token"])
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status == http.StatusOK || status == http.StatusCreated {
			_ = json.NewEncoder(w).Encode(map[string]string{"token": sessionToken, "expiresTimestamp": expires})
		}
	}))
	t.Cleanup(s.Close)
	return s
}

func writeTokenFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSignInServiceAccountPostsTokenAndMarksSession(t *testing.T) {
	// the server answers 201 to a POST, as the backend does; 200 is accepted too
	for _, status := range []int{http.StatusOK, http.StatusCreated} {
		srv := newServiceAccountServer(t, status, "sess-1", "2030-01-01T00:00:00Z")

		before := time.Now()
		s, err := SignInServiceAccount(context.Background(), Client{Endpoint: srv.URL + "/"}, "provider-token")
		if err != nil {
			t.Fatalf("status %d: %v", status, err)
		}
		if len(srv.tokens) != 1 || srv.tokens[0] != "provider-token" {
			t.Errorf("status %d: posted tokens = %v", status, srv.tokens)
		}
		if s.Token != "sess-1" || s.ExpiresTimestamp != "2030-01-01T00:00:00Z" || !s.ServiceAccount {
			t.Errorf("status %d: session = %+v", status, s)
		}
		created, err := time.Parse(time.RFC3339, s.CreatedTimestamp)
		if err != nil || created.Before(before.Add(-2*time.Second)) {
			t.Errorf("status %d: CreatedTimestamp = %q (%v)", status, s.CreatedTimestamp, err)
		}
	}
}

func TestSignInServiceAccountFixedTexts(t *testing.T) {
	for _, tc := range []struct {
		status int
		want   string
	}{
		{http.StatusUnauthorized, "The provider token wasn't accepted; check the token file and the service account's subject."},
		{http.StatusForbidden, "This MetalCloud doesn't accept service accounts."},
		{http.StatusNotFound, "This MetalCloud doesn't support signing in from the CLI; use an API token."},
	} {
		srv := newServiceAccountServer(t, tc.status, "", "")
		_, err := SignInServiceAccount(context.Background(), Client{Endpoint: srv.URL}, "secret-provider-token")
		if err == nil || err.Error() != tc.want {
			t.Errorf("status %d: err = %v, want %q", tc.status, err, tc.want)
		}
		if err != nil && strings.Contains(err.Error(), "secret-provider-token") {
			t.Errorf("status %d: the error carries the token", tc.status)
		}
	}
}

func TestSignInServiceAccountOtherStatusAndEmptyAnswer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"message":"Unable to verify the token at this time."}`))
	}))
	_, err := SignInServiceAccount(context.Background(), Client{Endpoint: srv.URL}, "t")
	srv.Close()
	if err == nil || !strings.Contains(err.Error(), "503") || !strings.Contains(err.Error(), "Unable to verify the token at this time.") {
		t.Errorf("503: err = %v", err)
	}

	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	_, err = SignInServiceAccount(context.Background(), Client{Endpoint: srv.URL}, "t")
	srv.Close()
	if err == nil {
		t.Error("an answer without a session token was accepted")
	}
}

func TestServiceAccountCredentialSignsInWhenThereIsNoMarkedSession(t *testing.T) {
	path := useTempSessionsPath(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "sess-1", "2030-01-01T00:00:00Z")
	file := writeTokenFile(t, "  provider-token\n")
	now := time.Date(2029, 12, 31, 0, 0, 0, 0, time.UTC)

	token, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, file, srv.URL, now)
	if err != nil || token != "sess-1" {
		t.Fatalf("token = %q, err = %v", token, err)
	}
	if len(srv.tokens) != 1 || srv.tokens[0] != "provider-token" {
		t.Errorf("the file content was not trimmed: %q", srv.tokens)
	}
	saved, _ := LoadServiceAccountSession(srv.URL, file)
	if saved == nil || saved.Token != "sess-1" || !saved.ServiceAccount {
		t.Errorf("saved = %+v", saved)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Errorf("sessions file mode = %o, want 600", info.Mode().Perm())
	}
}

func TestServiceAccountCredentialReusesASessionValidForMoreThan60Seconds(t *testing.T) {
	useTempSessionsPath(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "new", "2030-01-01T00:00:00Z")
	file := writeTokenFile(t, "provider-token")
	now := time.Date(2029, 12, 31, 0, 0, 0, 0, time.UTC)
	err := SaveServiceAccountSession(srv.URL, file, Session{Token: "saved", ExpiresTimestamp: now.Add(61 * time.Second).Format(time.RFC3339)})
	if err != nil {
		t.Fatal(err)
	}

	token, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, file, srv.URL, now)
	if err != nil || token != "saved" || srv.calls != 0 {
		t.Errorf("token = %q, calls = %d, err = %v", token, srv.calls, err)
	}
}

func TestServiceAccountCredentialSignsInAgainAt59SecondsLeft(t *testing.T) {
	useTempSessionsPath(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "new", "2030-01-01T00:00:00Z")
	file := writeTokenFile(t, "provider-token")
	now := time.Date(2029, 12, 31, 0, 0, 0, 0, time.UTC)
	err := SaveServiceAccountSession(srv.URL, file, Session{Token: "saved", ExpiresTimestamp: now.Add(59 * time.Second).Format(time.RFC3339)})
	if err != nil {
		t.Fatal(err)
	}

	token, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, file, srv.URL, now)
	if err != nil || token != "new" || srv.calls != 1 {
		t.Errorf("token = %q, calls = %d, err = %v", token, srv.calls, err)
	}
	if saved, _ := LoadServiceAccountSession(srv.URL, file); saved == nil || saved.Token != "new" {
		t.Errorf("saved = %+v", saved)
	}
}

func TestServiceAccountCredentialReadsTheFileFreshAndKeepsThePersonsSession(t *testing.T) {
	useTempSessionsPath(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "robot", "2030-01-01T00:00:00Z")
	file := writeTokenFile(t, "first")
	now := time.Date(2029, 12, 31, 0, 0, 0, 0, time.UTC)
	if err := SaveSession(srv.URL, Session{Token: "person", ExpiresTimestamp: "2030-06-01T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}

	// the person's valid session is ignored: a service-account session is signed in
	token, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, file, srv.URL, now)
	if err != nil || token != "robot" {
		t.Fatalf("token = %q, err = %v", token, err)
	}
	if person, _ := LoadSession(srv.URL); person == nil || person.Token != "person" {
		t.Errorf("the person's session was touched: %+v", person)
	}

	// the file is read again at the next sign-in, not cached
	if err := os.WriteFile(file, []byte("second\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	later := now.Add(48 * time.Hour) // past the first session's end
	if _, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, file, srv.URL, later); err != nil {
		t.Fatal(err)
	}
	if len(srv.tokens) != 2 || srv.tokens[0] != "first" || srv.tokens[1] != "second" {
		t.Errorf("posted tokens = %q", srv.tokens)
	}
}

func TestServiceAccountCredentialFileProblemsDoNotLeakContent(t *testing.T) {
	useTempSessionsPath(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "robot", "2030-01-01T00:00:00Z")
	now := time.Now()

	_, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, filepath.Join(t.TempDir(), "missing"), srv.URL, now)
	if err == nil || !strings.Contains(err.Error(), "token file") {
		t.Errorf("missing file: err = %v", err)
	}

	_, err = ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, writeTokenFile(t, " \n"), srv.URL, now)
	if err == nil || !strings.Contains(err.Error(), "token file is empty") {
		t.Errorf("empty file: err = %v", err)
	}
	if srv.calls != 0 {
		t.Errorf("the server was called %d times", srv.calls)
	}
}

// readOnlySessionsPath points the sessions file into a folder that can't be written, as a read-only
// $HOME would. It skips the test when the folder is writable anyway (running as root).
func readOnlySessionsPath(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	if err := os.Chmod(home, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(home, 0o700) })
	if probe, err := os.Create(filepath.Join(home, "probe")); err == nil {
		probe.Close()
		t.Skip("the folder is writable even when read-only (running as root)")
	}
	path := filepath.Join(home, ".metalcloud", "sessions.json")
	old := SessionsPath
	SessionsPath = func() (string, error) { return path, nil }
	t.Cleanup(func() { SessionsPath = old })
	return path
}

func captureWarnings(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	old := warnOut
	warnOut = &buf
	t.Cleanup(func() { warnOut = old })
	return &buf
}

func TestServiceAccountCredentialUsesTheTokenWhenTheSessionCannotBeSaved(t *testing.T) {
	path := readOnlySessionsPath(t)
	warnings := captureWarnings(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "sess-1", "2030-01-01T00:00:00Z")
	file := writeTokenFile(t, "provider-token")
	now := time.Date(2029, 12, 31, 0, 0, 0, 0, time.UTC)

	token, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, file, srv.URL, now)
	if err != nil || token != "sess-1" {
		t.Fatalf("token = %q, err = %v", token, err)
	}
	if srv.calls != 1 {
		t.Errorf("sign-ins = %d, want 1", srv.calls)
	}
	out := warnings.String()
	if strings.Count(out, "\n") != 1 || !strings.Contains(out, path) || !strings.Contains(out, "sign in again on the next command") {
		t.Errorf("warning = %q", out)
	}
	if strings.Contains(out, "sess-1") || strings.Contains(out, "provider-token") {
		t.Errorf("the warning carries a token: %q", out)
	}
}

func TestServiceAccountCredentialSaveWorksWithoutAWarning(t *testing.T) {
	useTempSessionsPath(t)
	warnings := captureWarnings(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "sess-1", "2030-01-01T00:00:00Z")

	_, err := ServiceAccountCredential(context.Background(), Client{Endpoint: srv.URL}, writeTokenFile(t, "provider-token"), srv.URL, time.Date(2029, 12, 31, 0, 0, 0, 0, time.UTC))
	if err != nil || warnings.Len() != 0 {
		t.Errorf("err = %v, warnings = %q", err, warnings.String())
	}
}

func TestSignInServiceAccountFromFileStillFailsWhenTheSessionCannotBeSaved(t *testing.T) {
	readOnlySessionsPath(t)
	srv := newServiceAccountServer(t, http.StatusCreated, "sess-1", "2030-01-01T00:00:00Z")

	_, err := SignInServiceAccountFromFile(context.Background(), Client{Endpoint: srv.URL}, writeTokenFile(t, "provider-token"), srv.URL)
	if err == nil || !strings.Contains(err.Error(), "could not be saved") {
		t.Errorf("err = %v", err)
	}
}
