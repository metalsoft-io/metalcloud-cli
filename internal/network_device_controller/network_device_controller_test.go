package network_device_controller

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

var controllerItem = map[string]any{
	"id": "23", "revision": 3, "status": "active", "siteId": 1,
	"identifierString": "ndfc-controller-01", "description": nil,
	"datacenterName": "dc1", "managementAddress": "10.0.0.50", "managementPort": 443,
	"username": "admin", "driver": "cisco_ndfc", "tags": []any{},
}

var controllerCredentials = map[string]any{
	"username": "admin", "password": "secret", "host": "10.0.0.50", "port": 443,
	"datacenter": "dc1", "driver": "cisco_ndfc", "hostname": "ndfc-controller-01",
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

func newServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-device-controllers": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, controllerItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{controllerItem}, 1, 1))(w, r)
		},
		"/api/v2/network-device-controllers/23": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, controllerItem)(w, r)
		},
		"/api/v2/network-device-controllers/23/credentials":            testutils.JSONHandler(200, controllerCredentials),
		"/api/v2/network-device-controllers/23/actions/deploy-confirm": testutils.NoContentHandler(),
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

func TestNetworkDeviceControllerListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error {
		return NetworkDeviceControllerList(ctx, NetworkDeviceControllerListFilters{SiteId: []string{"1"}})
	})
	if !strings.Contains(out, "ndfc-controller-01") {
		t.Errorf("list output missing controller: %s", out)
	}

	out = run(t, func() error { return NetworkDeviceControllerGet(ctx, "23") })
	if !strings.Contains(out, `"driver":"cisco_ndfc"`) {
		t.Errorf("get output missing driver: %s", out)
	}

	out = run(t, func() error { return NetworkDeviceControllerGet(ctx, "ndfc-controller-01") })
	if !strings.Contains(out, `"id":"23"`) {
		t.Errorf("get by identifier output missing id: %s", out)
	}

	if _, err := GetNetworkDeviceControllerByIdOrLabel(ctx, "missing-controller"); err == nil {
		t.Error("expected error for unknown identifier")
	}
}

func TestNetworkDeviceControllerCreateUpdateDelete(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return NetworkDeviceControllerCreate(ctx, sdk.CreateNetworkDeviceController{
			DatacenterName:     "dc1",
			Driver:             sdk.SWITCHCONTROLLERDRIVER_CISCO_NDFC,
			ManagementAddress:  "10.0.0.50",
			ManagementPort:     443,
			Username:           "admin",
			ManagementPassword: "secret",
		})
	})
	created := rec.last()
	if created.Method != http.MethodPost || created.Path != "/api/v2/network-device-controllers" {
		t.Errorf("unexpected create request: %s %s", created.Method, created.Path)
	}
	if !strings.Contains(created.Body, `"driver":"cisco_ndfc"`) {
		t.Errorf("create body wrong: %s", created.Body)
	}

	run(t, func() error {
		return NetworkDeviceControllerUpdate(ctx, "23", []byte(`{"description":"updated"}`))
	})
	updated := rec.last()
	// The SDK revision is numeric; it must be rendered as a plain string tag.
	if got := updated.Header.Get("If-Match"); got != "3" {
		t.Errorf("update If-Match = %q, want %q", got, "3")
	}
	if !strings.Contains(updated.Body, `"description":"updated"`) {
		t.Errorf("update body wrong: %s", updated.Body)
	}

	if err := NetworkDeviceControllerDelete(ctx, "23"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deleted := rec.last()
	if deleted.Method != http.MethodDelete || deleted.Path != "/api/v2/network-device-controllers/23" {
		t.Errorf("unexpected delete request: %s %s", deleted.Method, deleted.Path)
	}
	if got := deleted.Header.Get("If-Match"); got != "3" {
		t.Errorf("delete If-Match = %q, want %q", got, "3")
	}
}

func TestNetworkDeviceControllerCredentialsAndDeployConfirm(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return NetworkDeviceControllerGetCredentials(ctx, "23") })
	if !strings.Contains(out, `"password":"secret"`) {
		t.Errorf("credentials output wrong: %s", out)
	}

	if err := NetworkDeviceControllerDeployConfirm(ctx, "23"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	confirmed := rec.last()
	if confirmed.Method != http.MethodPost || confirmed.Path != "/api/v2/network-device-controllers/23/actions/deploy-confirm" {
		t.Errorf("unexpected deploy-confirm request: %s %s", confirmed.Method, confirmed.Path)
	}
}

func TestNetworkDeviceControllerConfigExample(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := NetworkDeviceControllerConfigExample(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var example sdk.CreateNetworkDeviceController
	if err := json.Unmarshal([]byte(out), &example); err != nil {
		t.Fatalf("config example is not valid CreateNetworkDeviceController JSON: %v (%s)", err, out)
	}
	if example.DatacenterName == "" || example.ManagementAddress == "" {
		t.Errorf("config example is incomplete: %s", out)
	}
}
