package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var subnetItem = map[string]interface{}{
	"id": 1.0, "label": "test-subnet", "name": "test-subnet",
	"annotations": map[string]interface{}{},
	"createdAt":   "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
	"revision": 1.0, "tags": map[string]interface{}{},
	"parentSubnetId": 0.0, "ipVersion": "ipv4",
	"networkAddress": "10.0.0.0", "prefixLength": 24.0,
	"netmask": "255.255.255.0", "defaultGatewayAddress": "10.0.0.1",
	"isPool": false, "allocationDenylist": []interface{}{}, "childOverlapAllowRules": []interface{}{},
}

func newSubnetTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/subnets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(subnetItem))
		})
		mux.HandleFunc("/api/v2/subnets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(subnetItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestSubnetList(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "subnet", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "10.0.0.0") {
		t.Errorf("expected output to contain '10.0.0.0', got: %s", out)
	}
	if !strings.Contains(out, "test-subnet") {
		t.Errorf("expected output to contain 'test-subnet', got: %s", out)
	}
}

func TestSubnetListAlias(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "subnets", "ls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "10.0.0.0") {
		t.Errorf("expected output to contain '10.0.0.0', got: %s", out)
	}
}

func TestSubnetGet(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "subnet", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-subnet") {
		t.Errorf("expected output to contain 'test-subnet', got: %s", out)
	}
}

func TestSubnetGetRequiresArg(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "subnet", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestSubnetHelp(t *testing.T) {
	out, err := runCLI(t, nil, "subnet", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "subnet") {
		t.Errorf("expected help output to contain 'subnet', got: %s", out)
	}
}

func TestSubnetCreate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/subnets", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, subnetItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(subnetItem))
		})
		mux.HandleFunc("/api/v2/subnets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(subnetItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "subnet-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"x","networkAddress":"10.0.0.0","prefixLength":24,"ipVersion":"ipv4","isPool":false}`)
	f.Close()

	_, execErr := runCLI(t, srv, "subnet", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestSubnetUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/subnets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(subnetItem))
		})
		mux.HandleFunc("/api/v2/subnets/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, subnetItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(subnetItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "subnet-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"updated"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "subnet", "update", "1", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestSubnetDelete(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/subnets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(subnetItem))
		})
		mux.HandleFunc("/api/v2/subnets/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(subnetItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, execErr := runCLI(t, srv, "subnet", "delete", "1")
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- subnet capacity, remove-ip and remove-ip-range ---

var subnetCapacityItem = map[string]interface{}{
	"subnetId": 1, "prefix": "10.0.0.0/24", "isPool": false,
	"totalIpCount": "254", "freeIpCount": "200", "usedIpCount": "54",
	"freePrefixes": nil, "freePrefixCounts": nil,
	"requestedPrefixLength": nil, "allocatablePrefixCount": nil,
}

// newSubnetSubResourceServer serves the subnet plus its capacity and the two
// deletable sub-resources, recording the If-Match of the last delete.
func newSubnetSubResourceServer(lastIfMatch *string) *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/subnets/1", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, subnetItem)
		})
		mux.HandleFunc("/api/v2/subnets/1/capacity", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, subnetCapacityItem)
		})
		mux.HandleFunc("/api/v2/subnets/1/ips/7", func(w http.ResponseWriter, r *http.Request) {
			*lastIfMatch = r.Header.Get("If-Match")
			w.WriteHeader(http.StatusNoContent)
		})
		mux.HandleFunc("/api/v2/subnets/1/ip-ranges/9", func(w http.ResponseWriter, r *http.Request) {
			*lastIfMatch = r.Header.Get("If-Match")
			w.WriteHeader(http.StatusNoContent)
		})
	}))
}

func TestSubnetCapacity(t *testing.T) {
	var ifMatch string
	srv := newSubnetSubResourceServer(&ifMatch)
	defer srv.Close()

	out, err := runCLI(t, srv, "subnet", "capacity", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "10.0.0.0/24") {
		t.Errorf("expected the prefix in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "subnet", "capacity"); err == nil {
		t.Error("expected an error when no subnet id is given")
	}
}

func TestSubnetRemoveIp(t *testing.T) {
	var ifMatch string
	srv := newSubnetSubResourceServer(&ifMatch)
	defer srv.Close()

	if _, err := runCLI(t, srv, "subnet", "remove-ip", "1", "7"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ifMatch != "1" {
		t.Errorf("expected If-Match 1 (the subnet revision), got %q", ifMatch)
	}

	if _, err := runCLI(t, srv, "subnet", "remove-ip", "1"); err == nil {
		t.Error("expected an error when the ip id is missing")
	}
}

func TestSubnetRemoveIpRange(t *testing.T) {
	var ifMatch string
	srv := newSubnetSubResourceServer(&ifMatch)
	defer srv.Close()

	if _, err := runCLI(t, srv, "subnet", "remove-ip-range", "1", "9"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ifMatch != "1" {
		t.Errorf("expected If-Match 1 (the subnet revision), got %q", ifMatch)
	}

	if _, err := runCLI(t, srv, "subnet", "remove-ip-range", "1"); err == nil {
		t.Error("expected an error when the range id is missing")
	}
}

func TestSubnetList_Formats(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "subnet", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}
