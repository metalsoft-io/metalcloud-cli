package dhcp_reservation

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// siteItem satisfies the required properties of sdk.Site. Sites are resolved
// by id or name, so both are asserted on.
var siteItem = map[string]any{
	"id": 1, "revision": 1, "slug": "dc-1", "name": "dc-1",
}

// reservationItem satisfies the required properties of sdk.DhcpReservation.
// Its revision (5) is the entity tag update/delete must echo in If-Match.
var reservationItem = map[string]any{
	"id": 12, "revision": 5, "siteId": 1, "ipVersion": "ipv4",
	"macAddress": "AA:BB:CC:DD:EE:FF", "priority": 100,
	"allocation": map[string]any{"kind": "manual", "ip": "192.168.1.10", "subnetId": 3},
}

type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

type recordedRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

func (r *recorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, recordedRequest{Method: req.Method, Path: req.URL.Path, Header: req.Header.Clone(), Body: string(body)})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *recorder) last() recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func newDhcpServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{siteItem}, 1, 1)),
		"/api/v2/sites/1/dhcp/ipv4/reservations": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(http.StatusCreated, reservationItem)(w, r)
				return
			}
			testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{reservationItem}, 1, 1))(w, r)
		},
		"/api/v2/sites/1/dhcp/ipv4/reservations/12": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(http.StatusOK, reservationItem)(w, r)
		},
	}
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func TestDhcpReservationListAndGet(t *testing.T) {
	rec := &recorder{}
	ts := newDhcpServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := DhcpReservationList(ctx, "1", "ipv4"); err != nil {
			t.Errorf("list: unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "AA:BB:CC:DD:EE:FF") {
		t.Errorf("list output missing the reservation: %s", out)
	}

	// The site is also resolvable by label.
	out = testutils.CaptureStdout(t, func() {
		if err := DhcpReservationGet(ctx, "dc-1", "ipv4", "12"); err != nil {
			t.Errorf("get: unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, `"id":12`) {
		t.Errorf("get output missing the reservation id: %s", out)
	}
	if rec.last().Path != "/api/v2/sites/1/dhcp/ipv4/reservations/12" {
		t.Errorf("unexpected path %s", rec.last().Path)
	}
}

func TestDhcpReservationCreate(t *testing.T) {
	rec := &recorder{}
	ts := newDhcpServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	testutils.CaptureStdout(t, func() {
		err := DhcpReservationCreate(ctx, "1", "ipv4", sdk.DhcpReservationCreate{
			MacAddress: sdk.PtrString("AA:BB:CC:DD:EE:FF"),
			Allocation: sdk.DhcpManualAllocationAsDhcpReservationAllocation(&sdk.DhcpManualAllocation{
				Kind: "manual", Ip: "192.168.1.10", SubnetId: 3,
			}),
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	req := rec.last()
	if req.Method != http.MethodPost || req.Path != "/api/v2/sites/1/dhcp/ipv4/reservations" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	for _, want := range []string{`"macAddress":"AA:BB:CC:DD:EE:FF"`, `"kind":"manual"`, `"subnetId":3`} {
		if !strings.Contains(req.Body, want) {
			t.Errorf("create body missing %s: %s", want, req.Body)
		}
	}
}

func TestDhcpReservationUpdate_SendsIfMatch(t *testing.T) {
	rec := &recorder{}
	ts := newDhcpServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	testutils.CaptureStdout(t, func() {
		config := []byte(`{"priority":50,"allocation":{"kind":"manual","ip":"192.168.1.11","subnetId":3}}`)
		if err := DhcpReservationUpdate(ctx, "1", "ipv4", "12", config); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	req := rec.last()
	if req.Method != http.MethodPut || req.Path != "/api/v2/sites/1/dhcp/ipv4/reservations/12" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	if got := req.Header.Get("If-Match"); got != "5" {
		t.Errorf("expected If-Match 5 (the reservation revision), got %q", got)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte(req.Body), &body)
	if body["priority"] != float64(50) {
		t.Errorf("update body should carry the new priority: %s", req.Body)
	}
}

func TestDhcpReservationDelete_SendsIfMatch(t *testing.T) {
	rec := &recorder{}
	ts := newDhcpServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DhcpReservationDelete(ctx, "1", "ipv4", "12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := rec.last()
	if req.Method != http.MethodDelete || req.Path != "/api/v2/sites/1/dhcp/ipv4/reservations/12" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	if got := req.Header.Get("If-Match"); got != "5" {
		t.Errorf("expected If-Match 5 (the reservation revision), got %q", got)
	}
}

func TestDhcpReservationInvalidIpVersion(t *testing.T) {
	rec := &recorder{}
	ts := newDhcpServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	err := DhcpReservationList(ctx, "1", "ipv5")
	if err == nil {
		t.Fatal("expected an error for an invalid IP version, got nil")
	}
	if !strings.Contains(err.Error(), "IpVersion") {
		t.Errorf("expected the SDK enum validation error, got: %v", err)
	}
}

func TestDhcpReservationUnknownSite(t *testing.T) {
	rec := &recorder{}
	ts := newDhcpServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DhcpReservationList(ctx, "does-not-exist", "ipv4"); err == nil {
		t.Fatal("expected an error for an unknown site, got nil")
	}
}

func TestDhcpReservationInvalidReservationId(t *testing.T) {
	rec := &recorder{}
	ts := newDhcpServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DhcpReservationGet(ctx, "1", "ipv4", "not-a-number"); err == nil {
		t.Fatal("expected an error for an invalid reservation ID, got nil")
	}
}

func TestDhcpReservationConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	out := testutils.CaptureStdout(t, func() {
		if err := DhcpReservationConfigExample(ctx); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	for _, want := range []string{`"kind":"manual"`, `"kind":"auto"`} {
		if !strings.Contains(out, want) {
			t.Errorf("config-example output missing %s: %s", want, out)
		}
	}
}
