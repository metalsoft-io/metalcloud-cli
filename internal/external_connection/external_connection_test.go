package external_connection

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

var externalConnectionItem = map[string]any{
	"id": "12", "label": "dc1-ext", "name": "DC1 external", "fabricId": 7,
	"revision":  "4",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var externalConnectionInterfaceItem = map[string]any{
	"id": 5, "networkDeviceInterfaceId": 101, "networkDeviceInterfaceName": "Ethernet1/1",
	"networkDeviceId": 45, "revision": "2",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var externalConnectionLogicalNetworkItem = map[string]any{
	"id": 3, "externalConnectionId": 12, "logicalNetworkId": 44, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var externalConnectionNetworkDevice = testutils.NetworkDeviceFixture("45", 1, "border-leaf-01")

// recorder captures the requests received by the mock server so tests can
// assert on method, path, headers and body.
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

func newServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/external-connections": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, externalConnectionItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{externalConnectionItem}, 1, 1))(w, r)
		},
		"/api/v2/external-connections/12": testutils.JSONHandler(200, externalConnectionItem),
		"/api/v2/external-connections/12/interfaces": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, externalConnectionInterfaceItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{externalConnectionInterfaceItem}, 1, 1))(w, r)
		},
		"/api/v2/external-connections/12/interfaces/5": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, externalConnectionInterfaceItem)(w, r)
		},
		"/api/v2/external-connections/12/logical-networks": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, externalConnectionLogicalNetworkItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{externalConnectionLogicalNetworkItem}, 1, 1))(w, r)
		},
		"/api/v2/external-connections/12/logical-networks/3": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, externalConnectionLogicalNetworkItem)(w, r)
		},
		"/api/v2/external-connections/network-devices/45/interfaces": testutils.JSONHandler(200, map[string]any{
			"data": []any{map[string]any{
				"networkDeviceId": 45, "networkDeviceInterfaceId": 101,
				"networkDeviceInterfaceName": "Ethernet1/1", "externalConnectionId": 12,
				"externalConnectionInterfaceId": 5,
			}},
		}),
		"/api/v2/network-devices/45": testutils.JSONHandler(200, externalConnectionNetworkDevice),
		"/api/v2/network-devices":    testutils.JSONHandler(200, testutils.PaginatedResponse([]any{externalConnectionNetworkDevice}, 1, 1)),
	}
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func run(t *testing.T, fn func() error) string {
	t.Helper()
	var out string
	var err error
	out = testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

func TestExternalConnectionListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error {
		return ExternalConnectionList(ctx, ExternalConnectionListFilters{FabricId: []string{"7"}})
	})
	if !strings.Contains(out, "dc1-ext") {
		t.Errorf("list output missing external connection: %s", out)
	}

	out = run(t, func() error { return ExternalConnectionGet(ctx, "12") })
	if !strings.Contains(out, `"label":"dc1-ext"`) {
		t.Errorf("get output missing label: %s", out)
	}
}

func TestExternalConnectionGetByLabel(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return ExternalConnectionGet(ctx, "dc1-ext") })
	if !strings.Contains(out, `"id":"12"`) {
		t.Errorf("get by label output missing id: %s", out)
	}

	if _, err := GetExternalConnectionByIdOrLabel(ctx, "missing"); err == nil {
		t.Error("expected error for unknown label")
	}
}

func TestExternalConnectionCreateAndUpdate(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return ExternalConnectionCreate(ctx, sdk.CreateExternalConnection{
			Label: "dc1-ext", Name: "DC1 external", FabricId: 7,
			ExternalConnectionInterfaces: []sdk.CreateExternalConnectionInterface{{NetworkDeviceInterfaceId: 101}},
		})
	})
	created := rec.last()
	if created.Method != http.MethodPost || created.Path != "/api/v2/external-connections" {
		t.Errorf("unexpected create request: %s %s", created.Method, created.Path)
	}
	if !strings.Contains(created.Body, `"networkDeviceInterfaceId":101`) {
		t.Errorf("create body missing interfaces: %s", created.Body)
	}

	run(t, func() error {
		return ExternalConnectionUpdate(ctx, "12", []byte(`{"name":"renamed"}`))
	})
	updated := rec.last()
	if updated.Method != http.MethodPatch && updated.Method != http.MethodPut {
		t.Errorf("unexpected update method: %s", updated.Method)
	}
	if got := updated.Header.Get("If-Match"); got != "4" {
		t.Errorf("update If-Match = %q, want %q", got, "4")
	}
	if !strings.Contains(updated.Body, `"name":"renamed"`) {
		t.Errorf("update body missing name: %s", updated.Body)
	}
}

func TestExternalConnectionDelete(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	// The delete route is served by the "/12" handler which returns the object;
	// the SDK ignores the body for deletes.
	if err := ExternalConnectionDelete(ctx, "12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deleted := rec.last()
	if deleted.Method != http.MethodDelete || deleted.Path != "/api/v2/external-connections/12" {
		t.Errorf("unexpected delete request: %s %s", deleted.Method, deleted.Path)
	}
}

func TestExternalConnectionInterfaces(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return ExternalConnectionInterfaceList(ctx, "12") })
	if !strings.Contains(out, "Ethernet1/1") {
		t.Errorf("interface list missing interface: %s", out)
	}

	out = run(t, func() error { return ExternalConnectionInterfaceGet(ctx, "12", "5") })
	if !strings.Contains(out, `"networkDeviceInterfaceId":101`) {
		t.Errorf("interface get missing id: %s", out)
	}

	run(t, func() error { return ExternalConnectionInterfaceAdd(ctx, "12", "101") })
	added := rec.last()
	if added.Method != http.MethodPost || added.Path != "/api/v2/external-connections/12/interfaces" {
		t.Errorf("unexpected add request: %s %s", added.Method, added.Path)
	}
	if !strings.Contains(added.Body, `"networkDeviceInterfaceId":101`) {
		t.Errorf("add body wrong: %s", added.Body)
	}

	run(t, func() error { return ExternalConnectionInterfaceUpdate(ctx, "12", "5", "102") })
	updated := rec.last()
	if got := updated.Header.Get("If-Match"); got != "2" {
		t.Errorf("interface update If-Match = %q, want the interface revision %q", got, "2")
	}
	if !strings.Contains(updated.Body, `"networkDeviceInterfaceId":102`) {
		t.Errorf("interface update body wrong: %s", updated.Body)
	}

	if err := ExternalConnectionInterfaceRemove(ctx, "12", "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	removed := rec.last()
	if removed.Method != http.MethodDelete || removed.Path != "/api/v2/external-connections/12/interfaces/5" {
		t.Errorf("unexpected remove request: %s %s", removed.Method, removed.Path)
	}
	if got := removed.Header.Get("If-Match"); got != "2" {
		t.Errorf("interface remove If-Match = %q, want %q", got, "2")
	}

	if err := ExternalConnectionInterfaceGet(ctx, "12", "not-a-number"); err == nil {
		t.Error("expected error for non-numeric interface ID")
	}
}

func TestExternalConnectionLogicalNetworks(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return ExternalConnectionLogicalNetworkList(ctx, "12") })
	if !strings.Contains(out, "44") {
		t.Errorf("logical network list missing entry: %s", out)
	}

	out = run(t, func() error { return ExternalConnectionLogicalNetworkGet(ctx, "12", "3") })
	if !strings.Contains(out, `"logicalNetworkId":44`) {
		t.Errorf("logical network get wrong: %s", out)
	}

	run(t, func() error { return ExternalConnectionLogicalNetworkAdd(ctx, "12", "44") })
	added := rec.last()
	if added.Method != http.MethodPost || added.Path != "/api/v2/external-connections/12/logical-networks" {
		t.Errorf("unexpected add request: %s %s", added.Method, added.Path)
	}
	if !strings.Contains(added.Body, `"logicalNetworkId":44`) {
		t.Errorf("add body wrong: %s", added.Body)
	}

	if err := ExternalConnectionLogicalNetworkRemove(ctx, "12", "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	removed := rec.last()
	if removed.Method != http.MethodDelete || removed.Path != "/api/v2/external-connections/12/logical-networks/3" {
		t.Errorf("unexpected remove request: %s %s", removed.Method, removed.Path)
	}
}

func TestNetworkDeviceInterfacesGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return NetworkDeviceInterfacesGet(ctx, "45") })
	if !strings.Contains(out, `"externalConnectionInterfaceId":5`) {
		t.Errorf("network device interfaces output wrong: %s", out)
	}
	if got := rec.last().Path; got != "/api/v2/external-connections/network-devices/45/interfaces" {
		t.Errorf("unexpected path: %s", got)
	}
}

func TestExternalConnectionConfigExample(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := ExternalConnectionConfigExample(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var example sdk.CreateExternalConnection
	if err := json.Unmarshal([]byte(out), &example); err != nil {
		t.Fatalf("config example is not valid CreateExternalConnection JSON: %v (%s)", err, out)
	}
	if example.Label == "" || example.FabricId == 0 {
		t.Errorf("config example is incomplete: %s", out)
	}
}
