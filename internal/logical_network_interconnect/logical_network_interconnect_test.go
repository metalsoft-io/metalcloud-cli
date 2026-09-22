package logical_network_interconnect

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

var interconnectItem = map[string]any{
	"id": "4", "revision": "6", "label": "dc1-dc2-ln", "name": "DC1 to DC2 network",
	"kind": "dci-evpn", "fabricInterconnectId": 12, "transportId": 899999999,
	"status":    "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var interconnectLinkItem = map[string]any{
	"id": "9", "logicalNetworkId": 44, "logicalNetworkInterconnectId": 4, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
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
		"/api/v2/logical-network-interconnects": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, interconnectItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{interconnectItem}, 1, 1))(w, r)
		},
		"/api/v2/logical-network-interconnects/4": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, interconnectItem)(w, r)
		},
		"/api/v2/logical-network-interconnects/4/links": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, interconnectLinkItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{interconnectLinkItem}, 1, 1))(w, r)
		},
		"/api/v2/logical-network-interconnects/4/links/9": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, interconnectLinkItem)(w, r)
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

func TestInterconnectListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error {
		return InterconnectList(ctx, LogicalNetworkInterconnectListFilters{Status: []string{"active"}, Kind: []string{"dci-evpn"}})
	})
	if !strings.Contains(out, "dc1-dc2-ln") {
		t.Errorf("list output missing interconnect: %s", out)
	}

	out = run(t, func() error { return InterconnectGet(ctx, "4") })
	if !strings.Contains(out, `"label":"dc1-dc2-ln"`) {
		t.Errorf("get output missing label: %s", out)
	}

	out = run(t, func() error { return InterconnectGet(ctx, "dc1-dc2-ln") })
	if !strings.Contains(out, `"id":"4"`) {
		t.Errorf("get by label output missing id: %s", out)
	}

	if _, err := GetInterconnectByIdOrLabel(ctx, "missing"); err == nil {
		t.Error("expected error for unknown label")
	}
}

func TestInterconnectCreateUpdateDelete(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	kind := sdk.LOGICALNETWORKINTERCONNECTKIND_DCI_EVPN
	run(t, func() error {
		return InterconnectCreate(ctx, sdk.CreateLogicalNetworkInterconnect{
			Label: "dc1-dc2-ln", Name: "DC1 to DC2 network", Kind: &kind, FabricInterconnectId: 12,
		})
	})
	created := rec.last()
	if created.Method != http.MethodPost || created.Path != "/api/v2/logical-network-interconnects" {
		t.Errorf("unexpected create request: %s %s", created.Method, created.Path)
	}
	if !strings.Contains(created.Body, `"fabricInterconnectId":12`) {
		t.Errorf("create body wrong: %s", created.Body)
	}

	run(t, func() error { return InterconnectUpdate(ctx, "4", []byte(`{"name":"renamed"}`)) })
	updated := rec.last()
	if got := updated.Header.Get("If-Match"); got != "6" {
		t.Errorf("update If-Match = %q, want %q", got, "6")
	}
	if !strings.Contains(updated.Body, `"name":"renamed"`) {
		t.Errorf("update body wrong: %s", updated.Body)
	}

	if err := InterconnectDelete(ctx, "4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deleted := rec.last()
	if deleted.Method != http.MethodDelete || deleted.Path != "/api/v2/logical-network-interconnects/4" {
		t.Errorf("unexpected delete request: %s %s", deleted.Method, deleted.Path)
	}
	if got := deleted.Header.Get("If-Match"); got != "6" {
		t.Errorf("delete If-Match = %q, want %q", got, "6")
	}
}

func TestInterconnectLinks(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return InterconnectLinkList(ctx, "4", []string{"44"}, []string{"active"}) })
	if !strings.Contains(out, "44") {
		t.Errorf("link list missing link: %s", out)
	}

	out = run(t, func() error { return InterconnectLinkGet(ctx, "4", "9") })
	if !strings.Contains(out, `"logicalNetworkId":44`) {
		t.Errorf("link get wrong: %s", out)
	}

	run(t, func() error { return InterconnectLinkAdd(ctx, "4", "44") })
	added := rec.last()
	if added.Method != http.MethodPost || added.Path != "/api/v2/logical-network-interconnects/4/links" {
		t.Errorf("unexpected link add request: %s %s", added.Method, added.Path)
	}
	if !strings.Contains(added.Body, `"logicalNetworkId":44`) {
		t.Errorf("link add body wrong: %s", added.Body)
	}

	if err := InterconnectLinkRemove(ctx, "4", "9"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	removed := rec.last()
	if removed.Method != http.MethodDelete || removed.Path != "/api/v2/logical-network-interconnects/4/links/9" {
		t.Errorf("unexpected link remove request: %s %s", removed.Method, removed.Path)
	}

	if err := InterconnectLinkGet(ctx, "4", "not-a-number"); err == nil {
		t.Error("expected error for non-numeric link ID")
	}
}

func TestInterconnectConfigExample(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := InterconnectConfigExample(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var example sdk.CreateLogicalNetworkInterconnect
	if err := json.Unmarshal([]byte(out), &example); err != nil {
		t.Fatalf("config example is not valid CreateLogicalNetworkInterconnect JSON: %v (%s)", err, out)
	}
	if example.Label == "" || example.FabricInterconnectId == 0 {
		t.Errorf("config example is incomplete: %s", out)
	}
}
