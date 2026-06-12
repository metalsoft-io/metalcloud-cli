package site

import (
	"net"
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func siteListResponse(id int, name string) map[string]interface{} {
	return map[string]interface{}{
		"data": []interface{}{siteItem(id, name)},
		"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
}

func siteConfigWithMapping(mapping map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"serverPolicy": map[string]interface{}{
			"dhcpOption82ToIPMapping":            mapping,
			"dhcpBmcMacAddressWhitelistEnabled":  false,
			"dhcpBmcMacAddressWhitelist":         []interface{}{},
			"automaticallyAllocateServerTypes":   false,
			"automaticallySetServersAsAvailable": false,
		},
	}
}

// --- Pure function tests ---

func TestIsValidMAC(t *testing.T) {
	tests := []struct {
		mac   string
		valid bool
	}{
		{"aa:bb:cc:dd:ee:ff", true},
		{"AA:BB:CC:DD:EE:FF", true},
		{"aa-bb-cc-dd-ee-ff", true},
		{"aabb.ccdd.eeff", true},
		{"not-a-mac", false},
		{"", false},
		{"gg:hh:ii:jj:kk:ll", false},
	}
	for _, tc := range tests {
		_, err := net.ParseMAC(tc.mac)
		got := err == nil
		if got != tc.valid {
			t.Errorf("net.ParseMAC(%q) valid=%v, want %v", tc.mac, got, tc.valid)
		}
	}
}

func TestIpBelongsToAnySubnet(t *testing.T) {
	_, net1, _ := net.ParseCIDR("192.168.1.0/24")
	_, net2, _ := net.ParseCIDR("10.0.0.0/8")
	subnets := []*net.IPNet{net1, net2}

	tests := []struct {
		ip      string
		belongs bool
	}{
		{"192.168.1.100", true},
		{"10.5.5.5", true},
		{"172.16.0.1", false},
		{"8.8.8.8", false},
	}
	for _, tc := range tests {
		ip := net.ParseIP(tc.ip)
		got := ipBelongsToAnySubnet(ip, subnets)
		if got != tc.belongs {
			t.Errorf("ipBelongsToAnySubnet(%s) = %v, want %v", tc.ip, got, tc.belongs)
		}
	}
}

func TestIpBelongsToAnySubnet_EmptySubnets(t *testing.T) {
	ip := net.ParseIP("192.168.1.1")
	if ipBelongsToAnySubnet(ip, nil) {
		t.Error("ipBelongsToAnySubnet with nil subnets: expected false, got true")
	}
}

func TestCheckDuplicateIPs_NoDuplicates(t *testing.T) {
	entries := map[string]string{
		"aa:bb:cc:dd:ee:01": "192.168.1.1",
		"aa:bb:cc:dd:ee:02": "192.168.1.2",
	}
	if err := checkDuplicateIPs(entries); err != nil {
		t.Fatalf("checkDuplicateIPs() unexpected error: %v", err)
	}
}

func TestCheckDuplicateIPs_WithDuplicate(t *testing.T) {
	entries := map[string]string{
		"aa:bb:cc:dd:ee:01": "192.168.1.1",
		"aa:bb:cc:dd:ee:02": "192.168.1.1",
	}
	if err := checkDuplicateIPs(entries); err == nil {
		t.Fatal("checkDuplicateIPs() expected error for duplicate IPs, got nil")
	}
}

func TestCheckDuplicateIPs_Empty(t *testing.T) {
	if err := checkDuplicateIPs(map[string]string{}); err != nil {
		t.Fatalf("checkDuplicateIPs() unexpected error on empty map: %v", err)
	}
}

// --- DhcpOobReservationsList ---

func TestDhcpOobReservationsList_HappyPath(t *testing.T) {
	mapping := map[string]interface{}{
		"aa:bb:cc:dd:ee:01": "192.168.1.1",
		"aa:bb:cc:dd:ee:02": "192.168.1.2",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":          testutils.JSONHandler(200, siteListResponse(1, "test-site")),
		"/api/v2/sites/1/config": testutils.JSONHandler(200, siteConfigWithMapping(mapping)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DhcpOobReservationsList(ctx, "test-site", "mac"); err != nil {
		t.Errorf("DhcpOobReservationsList: expected nil error, got: %v", err)
	}
}

func TestDhcpOobReservationsList_SortByIP(t *testing.T) {
	mapping := map[string]interface{}{
		"aa:bb:cc:dd:ee:01": "192.168.1.2",
		"aa:bb:cc:dd:ee:02": "192.168.1.1",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":          testutils.JSONHandler(200, siteListResponse(1, "test-site")),
		"/api/v2/sites/1/config": testutils.JSONHandler(200, siteConfigWithMapping(mapping)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DhcpOobReservationsList(ctx, "test-site", "ip"); err != nil {
		t.Errorf("DhcpOobReservationsList (sort by ip): expected nil error, got: %v", err)
	}
}

func TestDhcpOobReservationsList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":          testutils.JSONHandler(200, siteListResponse(1, "test-site")),
		"/api/v2/sites/1/config": testutils.JSONHandler(200, siteConfigWithMapping(map[string]interface{}{})),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DhcpOobReservationsList(ctx, "test-site", "mac"); err != nil {
		t.Errorf("DhcpOobReservationsList empty: expected nil error, got: %v", err)
	}
}

func TestDhcpOobReservationsList_SiteNotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(200, map[string]interface{}{
			"data": []interface{}{},
			"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DhcpOobReservationsList(ctx, "nonexistent", "mac"); err == nil {
		t.Error("DhcpOobReservationsList with nonexistent site: expected error, got nil")
	}
}

// --- DhcpOobReservationsRemove ---

func TestDhcpOobReservationsRemove_MacNotFound(t *testing.T) {
	mapping := map[string]interface{}{
		"aa:bb:cc:dd:ee:01": "192.168.1.1",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":          testutils.JSONHandler(200, siteListResponse(1, "test-site")),
		"/api/v2/sites/1/config": testutils.JSONHandler(200, siteConfigWithMapping(mapping)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	err := DhcpOobReservationsRemove(ctx, "test-site", []string{"ff:ff:ff:ff:ff:ff"})
	if err == nil {
		t.Error("DhcpOobReservationsRemove with non-existent MAC: expected error, got nil")
	}
}
