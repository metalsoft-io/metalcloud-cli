package server_instance

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestBuildServerFilterRules(t *testing.T) {
	rules, err := BuildServerFilterRules(
		[]string{"rack:equals:A12", "tags:matches_regexp:^gpu:.*$", "tags:anti_affinity"},
		[]string{"236", " 237 "},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	type want struct {
		property, operator string
		value              *string
	}
	str := func(s string) *string { return &s }

	expected := []want{
		{"rack", "equals", str("A12")},
		{"tags", "matches_regexp", str("^gpu:.*$")}, // the value keeps its colons
		{"tags", "anti_affinity", nil},              // no value for affinity operators
		{"server_id", "equals", str("236")},
		{"server_id", "equals", str("237")},
	}

	if len(rules) != len(expected) {
		t.Fatalf("expected %d rules, got %d", len(expected), len(rules))
	}
	for i, w := range expected {
		got := rules[i]
		if got.Property != w.property || got.Operator != w.operator {
			t.Errorf("rule %d: expected %s %s, got %s %s", i, w.property, w.operator, got.Property, got.Operator)
		}
		gotValue := got.Value.Get()
		switch {
		case w.value == nil && gotValue != nil:
			t.Errorf("rule %d: expected a null value, got %q", i, *gotValue)
		case w.value != nil && (gotValue == nil || *gotValue != *w.value):
			t.Errorf("rule %d: expected value %q, got %v", i, *w.value, gotValue)
		}
	}
}

func TestBuildServerFilterRules_Invalid(t *testing.T) {
	for _, rule := range []string{"rack", "", ":equals:A12", "rack::A12"} {
		if _, err := BuildServerFilterRules([]string{rule}, nil); err == nil {
			t.Errorf("rule %q: expected an error, got nil", rule)
		}
	}

	if _, err := BuildServerFilterRules(nil, []string{"not-a-number"}); err == nil {
		t.Error("expected an error for a non-numeric server ID")
	}
}

func TestBuildServerFilterRules_NoneGiven(t *testing.T) {
	rules, err := BuildServerFilterRules(nil, nil)
	if err != nil || len(rules) != 0 {
		t.Fatalf("expected no rules and no error, got %v and %v", rules, err)
	}
}

func TestServerInstanceGroupUpdate_SendsServerFilteringRules(t *testing.T) {
	const sigConfigJSON = `{"label":"sig-1","instanceCount":1,"defaultServerTypeId":5,"ipAllocateAuto":1,"ipv4SubnetCreateAuto":1,"processorCount":2,"processorCoreCount":8,"processorCoreMhz":2400,"diskCount":2,"diskSizeMbytes":102400,"diskTypes":[],"virtualInterfacesEnabled":0,"deployType":"deploy","deployStatus":"finished","revision":1,"updatedTimestamp":"2024-01-01T00:00:00Z"}`

	var mu sync.Mutex
	var patchBody map[string]any

	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instance-groups/10/config": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				raw, _ := io.ReadAll(r.Body)
				mu.Lock()
				_ = json.Unmarshal(raw, &patchBody)
				mu.Unlock()
			}
			testutils.RawHandler(http.StatusOK, sigConfigJSON)(w, r)
		},
	})
	defer ts.Close()

	rules, err := BuildServerFilterRules(nil, []string{"236"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := ServerInstanceGroupUpdate(setupTestContext(ts.URL), "10", "", 0, 0, rules); err != nil {
		t.Fatalf("ServerInstanceGroupUpdate: unexpected error: %v", err)
	}

	list, _ := patchBody["serverFilteringRules"].([]any)
	if len(list) != 1 {
		t.Fatalf("expected one serverFilteringRules entry, got %v", patchBody["serverFilteringRules"])
	}
	rule := list[0].(map[string]any)
	if rule["property"] != "server_id" || rule["operator"] != "equals" || rule["value"] != "236" {
		t.Errorf("unexpected rule: %v", rule)
	}
	if _, sent := patchBody["label"]; sent {
		t.Error("label must not be sent when it was not given")
	}
}
