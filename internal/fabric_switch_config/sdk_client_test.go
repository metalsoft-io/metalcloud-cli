package fabric_switch_config

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestManualStrategyBody guards the fix for the API rejecting a global scope with
// resourceId:0 ("scope.resourceId must not be less than 1"). The manual strategy
// body must carry scope {"kind":"global"} with NO resourceId key at all.
func TestManualStrategyBody(t *testing.T) {
	body := manualStrategyBody(7, "a_first")

	if body["kind"] != "manual" || body["subnetId"] != int64(7) || body["interfaceABinding"] != "a_first" {
		t.Fatalf("unexpected strategy body: %#v", body)
	}

	scope, ok := body["scope"].(map[string]any)
	if !ok {
		t.Fatalf("scope is not a map: %#v", body["scope"])
	}
	if scope["kind"] != "global" {
		t.Errorf("scope kind = %v, want global", scope["kind"])
	}
	if _, present := scope["resourceId"]; present {
		t.Errorf("global scope must omit resourceId, got %v", scope["resourceId"])
	}

	// Belt and suspenders: the marshaled JSON must not contain resourceId either.
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(raw), "resourceId") {
		t.Errorf("serialized strategy must not contain resourceId: %s", raw)
	}
}

// TestParseP2pLinksBody covers the lenient parse that replaces the SDK's typed
// decode (which fails with "data matches more than one schema in oneOf" once a
// link carries a manual/auto ipv4 strategy). The body below includes exactly
// that case.
func TestParseP2pLinksBody(t *testing.T) {
	body := []byte(`[
	  {
	    "id": 42, "revision": 3,
	    "interfaceA": {"type": "network_equipment_interface", "interfaceId": 1001},
	    "interfaceB": {"type": "network_equipment_interface", "interfaceId": 2002},
	    "config": {"ipv4": {"subnetAllocationStrategies": [
	      {"kind": "manual", "scope": {"kind": "global"}, "subnetId": 7, "interfaceABinding": "a_first"}
	    ]}}
	  },
	  {
	    "id": 43, "revision": 1,
	    "interfaceA": {"type": "network_equipment_interface", "interfaceId": 3003},
	    "interfaceB": {"type": "server_interface", "interfaceId": 9},
	    "config": {"ipv4": {"subnetAllocationStrategies": []}}
	  }
	]`)

	records, err := parseP2pLinksBody(body)
	if err != nil {
		t.Fatalf("parseP2pLinksBody: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("got %d records, want 2", len(records))
	}

	r0 := records[0]
	if r0.Id != 42 || r0.Revision != 3 {
		t.Errorf("link0 id/revision = %d/%d, want 42/3", r0.Id, r0.Revision)
	}
	if r0.InterfaceAId == nil || *r0.InterfaceAId != 1001 || r0.InterfaceBId == nil || *r0.InterfaceBId != 2002 {
		t.Errorf("link0 interface ids wrong: %v %v", r0.InterfaceAId, r0.InterfaceBId)
	}
	if !r0.HasIpv4Strategy {
		t.Errorf("link0 should report an existing ipv4 strategy")
	}

	r1 := records[1]
	// server_interface side is not a switch interface -> InterfaceBId stays nil.
	if r1.InterfaceAId == nil || *r1.InterfaceAId != 3003 {
		t.Errorf("link1 interfaceA id wrong: %v", r1.InterfaceAId)
	}
	if r1.InterfaceBId != nil {
		t.Errorf("link1 server_interface side should be nil, got %v", *r1.InterfaceBId)
	}
	if r1.HasIpv4Strategy {
		t.Errorf("link1 should have no ipv4 strategy")
	}
}
