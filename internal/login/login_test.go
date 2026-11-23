package login

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
