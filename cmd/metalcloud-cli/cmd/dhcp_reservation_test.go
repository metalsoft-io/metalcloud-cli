package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// dhcpResSiteItem satisfies sdk.Site's required fields; dhcpResItem satisfies
// sdk.DhcpReservation's. The reservation revision (5) is the entity tag the
// update/delete commands must echo back in If-Match.
var dhcpResSiteItem = map[string]interface{}{
	"id": 1, "revision": 1, "slug": "dc-1", "name": "dc-1",
}

var dhcpResItem = map[string]interface{}{
	"id": 12, "revision": 5, "siteId": 1, "ipVersion": "ipv4",
	"macAddress": "AA:BB:CC:DD:EE:FF", "priority": 100,
	"allocation": map[string]interface{}{"kind": "manual", "ip": "192.168.1.10", "subnetId": 3},
}

func newDhcpReservationTestServer(lastIfMatch *string, lastBody *string) *httptest.Server {
	record := func(r *http.Request) {
		*lastIfMatch = r.Header.Get("If-Match")
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(dhcpResSiteItem))
		})
		mux.HandleFunc("/api/v2/sites/1/dhcp/ipv4/reservations", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				record(r)
				jsonResponse(w, http.StatusCreated, dhcpResItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(dhcpResItem))
		})
		mux.HandleFunc("/api/v2/sites/1/dhcp/ipv4/reservations/12", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				*lastIfMatch = r.Header.Get("If-Match")
				w.WriteHeader(http.StatusNoContent)
			case http.MethodPut:
				record(r)
				jsonResponse(w, http.StatusOK, dhcpResItem)
			default:
				jsonResponse(w, http.StatusOK, dhcpResItem)
			}
		})
	}))
}

func TestDhcpReservationListAndAliases(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	for _, name := range []string{"dhcp-reservation", "dhcp-reservations"} {
		out, err := runCLI(t, srv, name, "list", "1", "ipv4")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "AA:BB:CC:DD:EE:FF") {
			t.Errorf("%s list: output missing the reservation: %s", name, out)
		}
	}

	// The site is also resolvable by label.
	if _, err := runCLI(t, srv, "dhcp-reservation", "ls", "dc-1", "ipv4"); err != nil {
		t.Fatalf("list by site label: unexpected error: %v", err)
	}
}

func TestDhcpReservationList_ArgsAndIpVersion(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "dhcp-reservation", "list", "1"); err == nil {
		t.Error("expected an error when the ip version is missing")
	}
	if _, err := runCLI(t, srv, "dhcp-reservation", "list", "1", "ipv5"); err == nil {
		t.Error("expected an error for an invalid ip version")
	}
}

func TestDhcpReservationGet(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "dhcp-reservation", "get", "1", "ipv4", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"id":12`) {
		t.Errorf("expected the reservation id in output, got: %s", out)
	}
}

func TestDhcpReservationConfigExample(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "dhcp-reservation", "config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"kind":"manual"`, `"kind":"auto"`} {
		if !strings.Contains(out, want) {
			t.Errorf("config-example output missing %s: %s", want, out)
		}
	}
}

func TestDhcpReservationCreate_ManualFlags(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "dhcp-reservation", "create", "1", "ipv4",
		"--mac-address", "AA:BB:CC:DD:EE:FF", "--ip", "192.168.1.10", "--subnet-id", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"kind":"manual"`, `"ip":"192.168.1.10"`, `"subnetId":3`, `"macAddress":"AA:BB:CC:DD:EE:FF"`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body missing %s: %s", want, body)
		}
	}
}

func TestDhcpReservationCreate_AutoFlags(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "dhcp-reservation", "create", "1", "ipv4",
		"--circuit-id", "leaf-01:swp1", "--subnet-pool-id", "3,4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"kind":"auto"`, `"subnetPoolIds":[3,4]`, `"circuitId":"leaf-01:swp1"`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body missing %s: %s", want, body)
		}
	}
}

func TestDhcpReservationCreate_RequiresAllocation(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "dhcp-reservation", "create", "1", "ipv4"); err == nil {
		t.Fatal("expected an error when no allocation flags and no --config-source are given")
	}
	// A manual allocation needs both --ip and --subnet-id.
	if _, err := runCLI(t, srv, "dhcp-reservation", "create", "1", "ipv4", "--ip", "192.168.1.10"); err == nil {
		t.Error("expected an error when --ip is given without --subnet-id")
	}
}

func TestDhcpReservationUpdate_SendsIfMatch(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "dhcp-res-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"priority":50,"allocation":{"kind":"manual","ip":"192.168.1.11","subnetId":3}}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "dhcp-reservation", "update", "1", "ipv4", "12", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
	if ifMatch != "5" {
		t.Errorf("expected If-Match 5 (the reservation revision), got %q", ifMatch)
	}
	if !strings.Contains(body, `"priority":50`) {
		t.Errorf("update body missing the new priority: %s", body)
	}
}

func TestDhcpReservationUpdate_RequiresConfigSource(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "dhcp-reservation", "update", "1", "ipv4", "12"); err == nil {
		t.Fatal("expected an error when --config-source is missing")
	}
}

func TestDhcpReservationDelete_SendsIfMatch(t *testing.T) {
	var ifMatch, body string
	srv := newDhcpReservationTestServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "dhcp-reservation", "rm", "1", "ipv4", "12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ifMatch != "5" {
		t.Errorf("expected If-Match 5 (the reservation revision), got %q", ifMatch)
	}
}
