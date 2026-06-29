package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// site_dhcp_oob_reservations_test.go covers:
//   site dhcp-oob-reservations list <site_id>
//   site dhcp-oob-reservations add <site_id> --mac ... --ip ...
//   site dhcp-oob-reservations remove <site_id> --mac ...

var siteItem = map[string]interface{}{
	"id":       1,
	"revision": 1.0,
	"slug":     "test-site",
	"name":     "test-site",
}

var siteConfigItem = map[string]interface{}{
	"serverPolicy": map[string]interface{}{
		"dhcpOption82ToIPMapping":             map[string]interface{}{"aa:bb:cc:dd:ee:ff": "10.0.0.1"},
		"dhcpBmcMacAddressWhitelistEnabled":   false,
		"dhcpBmcMacAddressWhitelist":          []interface{}{},
		"automaticallyAllocateServerTypes":    false,
		"automaticallySetServersAsAvailable":  false,
	},
}

func newSiteDhcpTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/sites/1/config", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodPut:
				_ = json.NewEncoder(w).Encode(siteConfigItem)
			default:
				_ = json.NewEncoder(w).Encode(siteConfigItem)
			}
		})
		mux.HandleFunc("/api/v2/sites/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(siteItem)
		})
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(siteItem))
		})
		mux.HandleFunc("/api/v2/subnets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList())
		})
	})
	return httptest.NewServer(mux)
}

func TestSiteDhcpOobReservationsList(t *testing.T) {
	srv := newSiteDhcpTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "site", "dhcp-oob-reservations", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "aa:bb:cc:dd:ee:ff") {
		t.Errorf("expected output to contain MAC address, got: %s", out)
	}
}

func TestSiteDhcpOobReservationsListRequiresArg(t *testing.T) {
	srv := newSiteDhcpTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "site", "dhcp-oob-reservations", "list")
	if err == nil {
		t.Fatal("expected error when no site arg provided, got nil")
	}
}
