package network_endpoint_group

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

var groupItem = map[string]any{
	"id": "3", "name": "dc1-endpoint-group", "siteId": 1, "revision": "5",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var groupLogicalNetworkItem = map[string]any{
	"logicalNetworkId": "44", "tagged": true, "accessMode": "l2",
	"networkEndpointGroupId": "3", "mtu": 9000,
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
		"/api/v2/network-endpoint-groups": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, groupItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{groupItem}, 1, 1))(w, r)
		},
		"/api/v2/network-endpoint-groups/3": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, groupItem)(w, r)
		},
		"/api/v2/network-endpoint-groups/3/logical-networks": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.NoContentHandler()(w, r)
				return
			}
			// This endpoint is not paginated: a bare {"data": [...]} envelope.
			testutils.JSONHandler(200, map[string]any{"data": []any{groupLogicalNetworkItem}})(w, r)
		},
		"/api/v2/network-endpoint-groups/3/logical-networks/44": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, groupLogicalNetworkItem)(w, r)
		},
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

func TestNetworkEndpointGroupListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error {
		return NetworkEndpointGroupList(ctx, NetworkEndpointGroupListFilters{SiteId: []string{"1"}})
	})
	if !strings.Contains(out, "dc1-endpoint-group") {
		t.Errorf("list output missing group: %s", out)
	}

	out = run(t, func() error { return NetworkEndpointGroupGet(ctx, "3") })
	if !strings.Contains(out, `"name":"dc1-endpoint-group"`) {
		t.Errorf("get output missing name: %s", out)
	}

	out = run(t, func() error { return NetworkEndpointGroupGet(ctx, "dc1-endpoint-group") })
	if !strings.Contains(out, `"id":"3"`) {
		t.Errorf("get by name output missing id: %s", out)
	}

	if _, err := GetNetworkEndpointGroupByIdOrLabel(ctx, "missing-group"); err == nil {
		t.Error("expected error for unknown name")
	}
}

func TestNetworkEndpointGroupCreateUpdateDelete(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return NetworkEndpointGroupCreate(ctx, sdk.CreateNetworkEndpointGroup{
			Name: "dc1-endpoint-group", SiteId: sdk.PtrInt64(1),
		})
	})
	created := rec.last()
	if created.Method != http.MethodPost || created.Path != "/api/v2/network-endpoint-groups" {
		t.Errorf("unexpected create request: %s %s", created.Method, created.Path)
	}
	if !strings.Contains(created.Body, `"name":"dc1-endpoint-group"`) {
		t.Errorf("create body wrong: %s", created.Body)
	}

	run(t, func() error { return NetworkEndpointGroupUpdate(ctx, "3", []byte(`{"name":"renamed"}`)) })
	updated := rec.last()
	if got := updated.Header.Get("If-Match"); got != "5" {
		t.Errorf("update If-Match = %q, want %q", got, "5")
	}
	if !strings.Contains(updated.Body, `"name":"renamed"`) {
		t.Errorf("update body wrong: %s", updated.Body)
	}

	if err := NetworkEndpointGroupDelete(ctx, "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deleted := rec.last()
	if deleted.Method != http.MethodDelete || deleted.Path != "/api/v2/network-endpoint-groups/3" {
		t.Errorf("unexpected delete request: %s %s", deleted.Method, deleted.Path)
	}
}

func TestNetworkEndpointGroupLogicalNetworks(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return NetworkEndpointGroupLogicalNetworkList(ctx, "3") })
	if !strings.Contains(out, `"logicalNetworkId":"44"`) {
		t.Errorf("logical network list wrong: %s", out)
	}

	out = run(t, func() error { return NetworkEndpointGroupLogicalNetworkGet(ctx, "3", "44") })
	if !strings.Contains(out, `"accessMode":"l2"`) {
		t.Errorf("logical network get wrong: %s", out)
	}

	err := NetworkEndpointGroupLogicalNetworkAdd(ctx, "3", sdk.CreateNetworkEndpointGroupLogicalNetwork{
		LogicalNetworkId: "44", Tagged: true, AccessMode: sdk.NETWORKENDPOINTGROUPALLOWEDACCESSMODE_L2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	added := rec.last()
	if added.Method != http.MethodPost || added.Path != "/api/v2/network-endpoint-groups/3/logical-networks" {
		t.Errorf("unexpected add request: %s %s", added.Method, added.Path)
	}
	if !strings.Contains(added.Body, `"logicalNetworkId":"44"`) || !strings.Contains(added.Body, `"accessMode":"l2"`) {
		t.Errorf("add body wrong: %s", added.Body)
	}

	run(t, func() error {
		return NetworkEndpointGroupLogicalNetworkUpdate(ctx, "3", "44", sdk.UpdateNetworkEndpointGroupLogicalNetwork{
			Mtu: sdk.PtrInt32(9000),
		})
	})
	updated := rec.last()
	// The attachment has no revision of its own: the group's revision guards it.
	if got := updated.Header.Get("If-Match"); got != "5" {
		t.Errorf("logical network update If-Match = %q, want the group revision %q", got, "5")
	}
	if !strings.Contains(updated.Body, `"mtu":9000`) {
		t.Errorf("logical network update body wrong: %s", updated.Body)
	}

	if err := NetworkEndpointGroupLogicalNetworkRemove(ctx, "3", "44"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	removed := rec.last()
	if removed.Method != http.MethodDelete || removed.Path != "/api/v2/network-endpoint-groups/3/logical-networks/44" {
		t.Errorf("unexpected remove request: %s %s", removed.Method, removed.Path)
	}

	if err := NetworkEndpointGroupLogicalNetworkGet(ctx, "3", "not-a-number"); err == nil {
		t.Error("expected error for non-numeric logical network ID")
	}
}

func TestNetworkEndpointGroupConfigExample(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := NetworkEndpointGroupConfigExample(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var example sdk.CreateNetworkEndpointGroup
	if err := json.Unmarshal([]byte(out), &example); err != nil {
		t.Fatalf("config example is not valid CreateNetworkEndpointGroup JSON: %v (%s)", err, out)
	}
	if example.Name == "" {
		t.Errorf("config example is incomplete: %s", out)
	}
}
