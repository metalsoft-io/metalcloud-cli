package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// vm_pool_test.go covers get, create, delete — list is already in vm_test.go.

func newVMPoolWriteTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-pools/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				_ = json.NewEncoder(w).Encode(vmPoolItem)
			}
		})
		mux.HandleFunc("/api/v2/vm-pools", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(vmPoolItem)
			default:
				_ = json.NewEncoder(w).Encode(paginatedList(vmPoolItem))
			}
		})
	})
	return httptest.NewServer(mux)
}

func TestVMPoolGet(t *testing.T) {
	srv := newVMPoolWriteTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-vm-pool") {
		t.Errorf("expected output to contain 'test-vm-pool', got: %s", out)
	}
}

func TestVMPoolGetRequiresArg(t *testing.T) {
	srv := newVMPoolWriteTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-pool", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestVMPoolCreate(t *testing.T) {
	srv := newVMPoolWriteTestServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "vmpool-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"siteId":1,"managementHost":"vcenter.example.com","managementPort":443,"name":"new-pool","type":"vmware","networkFabricId":1,"username":"admin","password":"secret"}`)
	f.Close()

	_, err = runCLI(t, srv, "vm-pool", "create", "--config-source", f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMPoolDelete(t *testing.T) {
	srv := newVMPoolWriteTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-pool", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMPoolDeleteRequiresArg(t *testing.T) {
	srv := newVMPoolWriteTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-pool", "delete")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}
