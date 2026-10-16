package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// groupBodyRecorder serves the routes server-instance-group create and update touch and
// records the JSON body of every write request.
type groupBodyRecorder struct {
	mu     sync.Mutex
	bodies []map[string]interface{}
}

func (g *groupBodyRecorder) record(r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	var body map[string]interface{}
	_ = json.Unmarshal(raw, &body)

	g.mu.Lock()
	g.bodies = append(g.bodies, body)
	g.mu.Unlock()
}

func (g *groupBodyRecorder) server(t *testing.T) *httptest.Server {
	t.Helper()

	// Slice flags keep their values between runCLI calls; do not leak ours into other tests.
	t.Cleanup(func() {
		serverInstanceGroupFlags.serverFilterRules = nil
		serverInstanceGroupFlags.serverIds = nil
	})

	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/server-types/100", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, serverTypeItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/server-instance-groups", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				g.record(r)
			}
			jsonResponse(w, http.StatusOK, serverInstanceGroupItem)
		})
		mux.HandleFunc("/api/v2/server-instance-groups/1/config", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				g.record(r)
			}
			jsonResponse(w, http.StatusOK, sigCmdGroupConfig)
		})
	})
	return httptest.NewServer(mux)
}

// rules returns the serverFilteringRules of the most recent recorded body as "property:operator:value".
func (g *groupBodyRecorder) rules(t *testing.T) []string {
	t.Helper()

	if len(g.bodies) == 0 {
		t.Fatal("no request body was recorded")
	}

	var result []string
	list, _ := g.bodies[len(g.bodies)-1]["serverFilteringRules"].([]interface{})
	for _, item := range list {
		rule := item.(map[string]interface{})
		value := ""
		if v, ok := rule["value"].(string); ok {
			value = v
		}
		result = append(result, rule["property"].(string)+":"+rule["operator"].(string)+":"+value)
	}
	return result
}

func containsRule(rules []string, want string) bool {
	for _, rule := range rules {
		if rule == want {
			return true
		}
	}
	return false
}

func TestServerInstanceGroupCreate_PinsToServer(t *testing.T) {
	rec := &groupBodyRecorder{}
	srv := rec.server(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-instance-group", "create", "1", "web", "100", "1", "--server-id", "236"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rules := rec.rules(t); !containsRule(rules, "server_id:equals:236") {
		t.Errorf("expected a server_id equals 236 rule, got %v", rules)
	}
}

func TestServerInstanceGroupCreate_SendsFilterRule(t *testing.T) {
	rec := &groupBodyRecorder{}
	srv := rec.server(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-instance-group", "create", "1", "web", "100", "2",
		"--server-filter-rule", "rack:equals:A12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rules := rec.rules(t); !containsRule(rules, "rack:equals:A12") {
		t.Errorf("expected a rack equals A12 rule, got %v", rules)
	}
}

func TestServerInstanceGroupUpdate_PinsToServer(t *testing.T) {
	rec := &groupBodyRecorder{}
	srv := rec.server(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-instance-group", "update", "1", "--server-id", "236"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rules := rec.rules(t); !containsRule(rules, "server_id:equals:236") {
		t.Errorf("expected a server_id equals 236 rule, got %v", rules)
	}
}
