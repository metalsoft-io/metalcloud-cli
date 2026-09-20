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

// ---------------------------------------------------------------------------
// update / sync / refresh / statistics / containers / cluster host commands
// ---------------------------------------------------------------------------

var vmpHostItem = map[string]interface{}{
	"id": 2.0, "name": "host-1", "poolId": 1.0,
	"address": "192.168.0.1", "port": 8443.0,
	"status": "online", "healthStatus": "normal",
	"allowVMsToBeCreated": true, "allowContainersToBeCreated": true,
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

var vmpHostInterfaceItem = map[string]interface{}{
	"id": 3.0, "hostId": 2.0, "status": "managed",
	"name": "phy0", "macAddress": "14:18:77:4b:5b:05",
	"networkDevices": []interface{}{},
}

var vmpNetworkDeviceAssignmentItem = map[string]interface{}{
	"id": 4.0, "hostInterfaceId": 3.0,
	"networkDeviceId": 5.0, "networkDeviceInterfaceName": "ethernet1/1/5",
}

var vmpStatisticsItem = map[string]interface{}{
	"totalRamGB": 67.0, "freeRamGB": 58.0, "usedRamGB": 9.0,
	"totalSpaceGB": 588.0, "usedSpaceGB": 17.0, "freeSpaceGB": 571.0,
	"gpuInfo": []interface{}{map[string]interface{}{"name": "G200eR2", "count": 1.0}},
}

// newVMPoolExtendedTestServer serves every route the new vm-pool sub-commands
// touch. lastBody, when non-nil, receives the JSON body of the last write.
func newVMPoolExtendedTestServer(lastBody *string) *httptest.Server {
	record := func(r *http.Request) {
		if lastBody == nil {
			return
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-pools/1", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			jsonResponse(w, http.StatusOK, vmPoolItem)
		})
		mux.HandleFunc("/api/v2/vm-pools/1/actions/sync", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"jobId": 7.0, "jobGroupId": 3.0})
		})
		mux.HandleFunc("/api/v2/vm-pools/1/actions/refresh-information", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, vmPoolItem)
		})
		mux.HandleFunc("/api/v2/vm-pools/1/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, vmpStatisticsItem)
		})
		mux.HandleFunc("/api/v2/vm-pools/1/containers", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList())
		})
		mux.HandleFunc("/api/v2/vm-pools/1/cluster-hosts/2", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			jsonResponse(w, http.StatusOK, vmpHostItem)
		})
		mux.HandleFunc("/api/v2/vm-pools/1/cluster-hosts/2/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, vmpStatisticsItem)
		})
		mux.HandleFunc("/api/v2/vm-pools/1/cluster-hosts/2/containers", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList())
		})
		mux.HandleFunc("/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			jsonResponse(w, http.StatusOK, vmpHostInterfaceItem)
		})
		mux.HandleFunc("/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				record(r)
				jsonResponse(w, http.StatusCreated, vmpNetworkDeviceAssignmentItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(vmpNetworkDeviceAssignmentItem))
		})
		mux.HandleFunc("/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices/4", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, vmpNetworkDeviceAssignmentItem)
		})
	})
	return httptest.NewServer(mux)
}

func vmpConfigFile(t *testing.T, contents string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vmpool-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(contents); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestVMPoolUpdateCmd(t *testing.T) {
	var body string
	srv := newVMPoolExtendedTestServer(&body)
	defer srv.Close()

	path := vmpConfigFile(t, `{"description":"updated"}`)
	if _, err := runCLI(t, srv, "vm-pool", "update", "1", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"description":"updated"`) {
		t.Errorf("update body missing description: %s", body)
	}
}

func TestVMPoolUpdateCmdRequiresConfigSource(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-pool", "update", "1"); err == nil {
		t.Fatal("expected error when --config-source is missing, got nil")
	}
}

func TestVMPoolUpdateConfigExampleCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "update-config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "managementHost") {
		t.Errorf("expected example to contain managementHost, got: %s", out)
	}
}

func TestVMPoolSyncCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "sync", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "jobGroupId") {
		t.Errorf("expected job info in output, got: %s", out)
	}
}

func TestVMPoolRefreshCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "refresh", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-vm-pool") {
		t.Errorf("expected refreshed pool in output, got: %s", out)
	}
}

func TestVMPoolStatisticsCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	for _, name := range []string{"statistics", "stats"} {
		out, err := runCLI(t, srv, "vm-pool", name, "1")
		if err != nil {
			t.Fatalf("vm-pool %s: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "totalRamGB") {
			t.Errorf("vm-pool %s: expected statistics in output, got: %s", name, out)
		}
	}
}

func TestVMPoolStatisticsCmdFlattensGpusInTable(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLIFormat(t, srv, "csv", "vm-pool", "statistics", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "G200eR2 x1") {
		t.Errorf("expected flattened GPU cell in csv output, got: %s", out)
	}
}

func TestVMPoolContainersCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-pool", "containers", "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMPoolClusterHostCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "cluster-host", "1", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "host-1") {
		t.Errorf("expected host in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "vm-pool", "cluster-host", "1"); err == nil {
		t.Error("expected error when host id is missing")
	}
}

func TestVMPoolUpdateClusterHostCmd(t *testing.T) {
	var body string
	srv := newVMPoolExtendedTestServer(&body)
	defer srv.Close()

	path := vmpConfigFile(t, `{"allowVMsToBeCreated":false}`)
	if _, err := runCLI(t, srv, "vm-pool", "update-cluster-host", "1", "2", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"allowVMsToBeCreated":false`) {
		t.Errorf("update body missing allowVMsToBeCreated: %s", body)
	}
}

func TestVMPoolClusterHostStatisticsCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "cluster-host-statistics", "1", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "totalRamGB") {
		t.Errorf("expected statistics in output, got: %s", out)
	}
}

func TestVMPoolClusterHostContainersCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-pool", "cluster-host-containers", "1", "2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMPoolClusterHostInterfaceCmd(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "cluster-host-interface", "1", "2", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "phy0") {
		t.Errorf("expected interface in output, got: %s", out)
	}
}

func TestVMPoolUpdateClusterHostInterfaceCmd(t *testing.T) {
	var body string
	srv := newVMPoolExtendedTestServer(&body)
	defer srv.Close()

	path := vmpConfigFile(t, `{"status":"managed"}`)
	if _, err := runCLI(t, srv, "vm-pool", "update-cluster-host-interface", "1", "2", "3", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"status":"managed"`) {
		t.Errorf("update body missing status: %s", body)
	}
}

func TestVMPoolNetworkDeviceAssignmentList(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	for _, group := range []string{"cluster-host-interface-network-device", "chi-network-device", "chind"} {
		out, err := runCLI(t, srv, "vm-pool", group, "list", "1", "2", "3")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", group, err)
		}
		if !strings.Contains(out, "ethernet1/1/5") {
			t.Errorf("%s list: expected assignment in output, got: %s", group, out)
		}
	}
}

func TestVMPoolNetworkDeviceAssignmentGet(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "chind", "get", "1", "2", "3", "4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ethernet1/1/5") {
		t.Errorf("expected assignment in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "vm-pool", "chind", "get", "1", "2", "3"); err == nil {
		t.Error("expected error when assignment id is missing")
	}
}

func TestVMPoolNetworkDeviceAssignmentAddFromConfig(t *testing.T) {
	var body string
	srv := newVMPoolExtendedTestServer(&body)
	defer srv.Close()

	path := vmpConfigFile(t, `{"networkDeviceId":5,"networkDeviceInterfaceName":"ethernet1/1/5"}`)
	if _, err := runCLI(t, srv, "vm-pool", "chind", "add", "1", "2", "3", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"networkDeviceId":5`, `"networkDeviceInterfaceName":"ethernet1/1/5"`} {
		if !strings.Contains(body, want) {
			t.Errorf("add body missing %s: %s", want, body)
		}
	}
}

func TestVMPoolNetworkDeviceAssignmentAddRequiresFlags(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-pool", "chind", "add", "1", "2", "3"); err == nil {
		t.Fatal("expected error when neither --config-source nor --network-device is given")
	}

	if _, err := runCLI(t, srv, "vm-pool", "chind", "add", "1", "2", "3", "--network-device", "5"); err == nil {
		t.Fatal("expected error when --network-device-interface is missing")
	}
}

func TestVMPoolNetworkDeviceAssignmentRemove(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	for _, name := range []string{"remove", "rm", "delete"} {
		if _, err := runCLI(t, srv, "vm-pool", "chind", name, "1", "2", "3", "4"); err != nil {
			t.Fatalf("chind %s: unexpected error: %v", name, err)
		}
	}
}

func TestVMPoolNetworkDeviceAssignmentConfigExample(t *testing.T) {
	srv := newVMPoolExtendedTestServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "chind", "config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "networkDeviceInterfaceName") {
		t.Errorf("expected example in output, got: %s", out)
	}
}
