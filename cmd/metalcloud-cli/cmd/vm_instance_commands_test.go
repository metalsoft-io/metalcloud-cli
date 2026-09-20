package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// vm_instance_commands_test.go covers the VM instance and VM instance group
// commands added on top of the original list/get/config/power set.

var vmiGroupItem = map[string]interface{}{
	"id":                1.0,
	"revision":          1.0,
	"label":             "vmig-1",
	"instanceGroupType": "vm",
	"infrastructureId":  1.0,
	"infrastructure":    map[string]interface{}{"id": 1.0},
	"serviceStatus":     "active",
	"diskSizeGB":        20.0,
	"instanceCount":     2.0,
	"createdTimestamp":  "2024-01-01T00:00:00Z",
	"updatedTimestamp":  "2024-01-01T00:00:00Z",
	"config": map[string]interface{}{
		"revision":         1.0,
		"label":            "vmig-1",
		"deployType":       "soft",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	},
	"meta": map[string]interface{}{},
}

var vmiInstanceConfigItem = map[string]interface{}{
	"revision":         1.0,
	"label":            "test-vm-instance",
	"typeId":           1.0,
	"deployType":       "soft",
	"deployStatus":     "not_started",
	"diskSizeGB":       20.0,
	"ramGB":            8.0,
	"cpuCores":         4.0,
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

var vmigConfigItem = map[string]interface{}{
	"revision":         1.0,
	"label":            "vmig-1",
	"instanceCount":    2.0,
	"deployType":       "soft",
	"deployStatus":     "not_started",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

// vmigInterfaceItem mirrors the live payload: the label only exists inside
// 'config', so the interface endpoints are read raw by the CLI.
var vmigInterfaceItem = map[string]interface{}{
	"id":               2.0,
	"revision":         4.0,
	"serviceStatus":    "active",
	"groupId":          1.0,
	"infrastructureId": 1.0,
	"index":            1.0,
	"networkId":        nil,
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"meta":             map[string]interface{}{},
	"config": map[string]interface{}{
		"revision":         3.0,
		"label":            "if1",
		"index":            1.0,
		"deployType":       "create",
		"deployStatus":     "finished",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	},
}

var vmigNetworkEndpointGroupItem = map[string]interface{}{
	"id":               "1019",
	"name":             "vmig-1",
	"revision":         "1",
	"siteId":           1.0,
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

var vmigConnectionItem = map[string]interface{}{
	"id":                      "535",
	"tagged":                  true,
	"accessMode":              "l2",
	"mtu":                     1500.0,
	"providesDefaultRoute":    false,
	"disableAutoIpAllocation": false,
}

// vmiContextBody is the raw variables/os-installation-data payload: 'server' is
// null and 'network' is an object, both of which the strict SDK models reject.
const vmiContextBody = `{"site":{"id":1},"siteConfig":{},"server":null,` +
	`"serverInstance":{"id":1},"serverInstanceGroup":{"id":1},` +
	`"infrastructure":{"id":1},"driveGroups":[],"drives":[],"fileShares":[],` +
	`"buckets":[],"sharedDrives":[],"network":{"interfaces":[]},` +
	`"variables":{},"secrets":{},"userSSHKeys":[],"managementSSHKey":null}`

func vmiNewCommandsServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		// VM instance routes.
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(vmInstanceItem)
				return
			}
			_ = json.NewEncoder(w).Encode(paginatedList(vmInstanceItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmInstanceItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/config", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmiInstanceConfigItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/meta", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmInstanceItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/actions/apply-type/5", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmInstanceItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/credentials", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"username":        "root",
				"initialPassword": "secret",
			})
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/variables", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vmiContextBody))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/os-installation-data", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vmiContextBody))
		})

		// VM instance group routes.
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmiGroupItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/config", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmigConfigItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/meta", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmiGroupItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/actions/apply-type/5", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmiGroupItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/interfaces", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmigInterfaceItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/interfaces/2", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmigInterfaceItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/vm-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmInstanceItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/config/networking", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmigNetworkEndpointGroupItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/config/networking/connections", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(vmigConnectionItem)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{vmigConnectionItem}})
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups/1/config/networking/connections/535", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmigConnectionItem)
		})
	})
	return httptest.NewServer(mux)
}

// VM instance commands

func TestVMInstanceConfigExampleCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-instance", "config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "typeId") {
		t.Errorf("expected the example to list typeId, got: %s", out)
	}
}

func TestVMInstanceCreateCmd_FromFlags(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "create", "1", "--vm-type-id", "1", "--group-id", "1", "--disk-size-gb", "20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceCreateCmd_RequiresFlags(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-instance", "create", "1"); err == nil {
		t.Fatal("expected error when neither --config-source nor --vm-type-id is given")
	}
}

func TestVMInstanceDeleteCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "delete", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceUpdateConfigCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	config := vmiWriteTempJSON(t, `{"label":"renamed"}`)

	_, err := runCLI(t, srv, "vm-instance", "update-config", "1", "1", "--config-source", config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceUpdateMetaCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	config := vmiWriteTempJSON(t, `{"tags":["prod"]}`)

	_, err := runCLI(t, srv, "vm-instance", "update-meta", "1", "1", "--config-source", config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceApplyTypeCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "apply-type", "1", "1", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceVariablesCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-instance", "variables", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "serverInstance") {
		t.Errorf("expected the variables document, got: %s", out)
	}
}

func TestVMInstanceOsInstallationDataCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "os-installation-data", "1", "1", "--remove-empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceCredentialsCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "credentials", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// VM instance group commands

func TestVMInstanceGroupConfigCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "config", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupUpdateMetaCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	config := vmiWriteTempJSON(t, `{"tags":["prod"]}`)

	_, err := runCLI(t, srv, "vm-instance-group", "update-meta", "1", "1", "--config-source", config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupApplyTypeCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "apply-type", "1", "1", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupInterfacesCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-instance-group", "interfaces", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "if1") {
		t.Errorf("expected the interface label from 'config', got: %s", out)
	}
}

func TestVMInstanceGroupInterfaceCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "interface", "1", "1", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupInstancesCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "instances", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// VM instance group network commands

func TestVMInstanceGroupNetworkListCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "network", "list", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupNetworkConnectionsCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "network", "connections", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupNetworkGetCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "network", "get", "1", "1", "535")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupNetworkConfigExampleCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-instance-group", "network", "config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "logicalNetworkId") {
		t.Errorf("expected the example to list logicalNetworkId, got: %s", out)
	}
}

func TestVMInstanceGroupNetworkConnectCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "network", "connect", "1", "1",
		"--logical-network-id", "5", "--access-mode", "l2", "--tagged", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupNetworkConnectCmd_RequiresFlags(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-instance-group", "network", "connect", "1", "1"); err == nil {
		t.Fatal("expected error when neither --config-source nor --logical-network-id is given")
	}
}

func TestVMInstanceGroupNetworkUpdateCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "network", "update", "1", "1", "535", "--tagged", "false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGroupNetworkDisconnectCmd(t *testing.T) {
	srv := vmiNewCommandsServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "network", "disconnect", "1", "1", "535")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// vmiWriteTempJSON writes content to a temporary file and returns its path, so
// --config-source can be exercised without reading from the pipe.
func vmiWriteTempJSON(t *testing.T, content string) string {
	t.Helper()
	path := t.TempDir() + "/config.json"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}
