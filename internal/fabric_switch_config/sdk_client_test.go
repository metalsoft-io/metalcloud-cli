package fabric_switch_config

import "testing"

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
