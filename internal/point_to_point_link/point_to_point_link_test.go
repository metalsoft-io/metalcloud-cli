package point_to_point_link

import (
	"encoding/json"
	"testing"
)

// TestP2pLinkDisplayParse confirms the lenient display parse decodes a link that
// carries a manual ipv4 strategy - the exact shape that makes the SDK's typed
// oneOf decode fail with "data matches more than one schema".
func TestP2pLinkDisplayParse(t *testing.T) {
	body := []byte(`[
	  {
	    "id": 42, "label": "fab5-sw4-swp33s0-to-sw2-swp1s0",
	    "description": "leaf<->spine", "routingActivation": "default",
	    "serviceStatus": "active", "revision": 3,
	    "config": {"ipv4": {"subnetAllocationStrategies": [
	      {"kind": "manual", "scope": {"kind": "global"}, "subnetId": 7, "interfaceABinding": "a_first"}
	    ]}}
	  },
	  {"id": 43, "label": "x", "description": null, "routingActivation": "default", "serviceStatus": "active", "revision": 1}
	]`)

	var links []p2pLinkDisplay
	if err := json.Unmarshal(body, &links); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("got %d links, want 2", len(links))
	}
	if links[0].Id != 42 || links[0].Label != "fab5-sw4-swp33s0-to-sw2-swp1s0" ||
		links[0].RoutingActivation != "default" || links[0].ServiceStatus != "active" || links[0].Revision != 3 {
		t.Errorf("link0 fields wrong: %+v", links[0])
	}
	// null description decodes to "".
	if links[1].Description != "" {
		t.Errorf("link1 null description should be empty, got %q", links[1].Description)
	}
}
