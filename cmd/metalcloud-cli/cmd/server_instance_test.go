package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// infraItem is a minimal Infrastructure JSON for ID=1, label="test-infra".
const infraItem = `{"id":1,"label":"test-infra","serviceStatus":"active","revision":1,"datacenterName":"dc1","siteId":1,"designIsLocked":0,"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","config":{"deployStatus":"not_started","label":"test-infra","deployType":"soft"}}`

// infraListBody is the paginated list response for GET /api/v2/infrastructures.
// GetInfrastructureByIdOrLabel calls GET /api/v2/infrastructures?search=<id_or_label>.
const infraListBody = `{"data":[` + infraItem + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`

// newInfraMux returns a ServeMux with /api/v2/user and /api/v2/infrastructures
// pre-registered, plus any extra routes from the extra callback.
// Used by infra-scoped commands that call GetInfrastructureByIdOrLabel.
func newInfraMux(extra func(mux *http.ServeMux)) *http.ServeMux {
	return newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(infraListBody))
		})
		if extra != nil {
			extra(mux)
		}
	})
}

var serverInstanceItem = map[string]interface{}{
	"id":                 1.0,
	"revision":           1.0,
	"label":              "si-1",
	"createdTimestamp":   "2024-01-01T00:00:00Z",
	"updatedTimestamp":   "2024-01-01T00:00:00Z",
	"infrastructureId":   1.0,
	"groupId":            1.0,
	"serviceStatus":      "active",
	"isVmInstance":       0.0,
	"isEndpointInstance": 0.0,
	"meta":               map[string]interface{}{"name": "si-1"},
}

var serverInstanceGroupItem = map[string]interface{}{
	"id":                       1.0,
	"revision":                 1.0,
	"label":                    "sig-1",
	"createdTimestamp":         "2024-01-01T00:00:00Z",
	"updatedTimestamp":         "2024-01-01T00:00:00Z",
	"infrastructureId":         1.0,
	"instanceCount":            1.0,
	"defaultServerTypeId":      0.0,
	"ipAllocateAuto":           1.0,
	"ipv4SubnetCreateAuto":     1.0,
	"processorCount":           1.0,
	"processorCoreCount":       1.0,
	"processorCoreMhz":         1000.0,
	"diskCount":                1.0,
	"diskSizeMbytes":           10240.0,
	"diskTypes":                []interface{}{},
	"virtualInterfacesEnabled": 0.0,
	"serviceStatus":            "active",
	"isVmGroup":                0.0,
	"isEndpointInstanceGroup":  0.0,
	"meta":                     map[string]interface{}{},
}

func newServerInstanceTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/server-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(serverInstanceItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/server-instance-groups", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(serverInstanceGroupItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestServerInstanceList(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-instance", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerInstanceGroupList(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-instance-group", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerInstanceList_Formats(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server-instance", "list", "1")
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

func TestServerInstanceGroupList_Formats(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server-instance-group", "list", "1")
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

// ---------------------------------------------------------------------------
// Fixtures and server for the server-instance / server-instance-group commands
// added on top of the original list/get set. Every map satisfies the required
// properties of its SDK model. Package level names are prefixed si/sig.
// ---------------------------------------------------------------------------

var siCmdInstanceConfig = map[string]interface{}{
	"revision": 3.0, "label": "si-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"groupId": 1.0, "serverTypeId": 5.0, "serverId": 100.0, "osTemplateId": 1.0,
	"hostname": "si-1.example.com", "deployType": "create", "deployStatus": "not_started",
}

var siCmdDrive = map[string]interface{}{
	"id": 8.0, "revision": 1.0, "label": "drive-1", "groupId": 4.0, "sizeMb": 40960.0,
	"storageType": "iscsi_ssd", "infrastructureId": 1.0,
	"infrastructure": map[string]interface{}{"id": 1.0, "label": "test-infra"},
	"serviceStatus":  "active", "storageUpdatedTimestamp": "2024-01-02T00:00:00Z",
	"provisioningProtocol": "iscsi", "meta": map[string]interface{}{"name": "drive-1"},
	"config": map[string]interface{}{
		"revision": 1.0, "label": "drive-1", "groupId": 4.0, "sizeMb": 40960.0,
		"storageType": "iscsi_ssd", "updatedTimestamp": "2024-01-02T00:00:00Z",
		"deployType": "create", "deployStatus": "not_started",
	},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var sigCmdDriveGroup = map[string]interface{}{
	"id": 4.0, "revision": 1.0, "label": "drive-group-1", "infrastructureId": 1.0,
	"driveSizeMbDefault": 40960.0, "expandWithServerInstanceGroup": 1.0,
	"storageType": "iscsi_ssd", "serviceStatus": "active", "allocationAffinity": "same_storage",
	"meta": map[string]interface{}{"name": "drive-group-1"},
	"config": map[string]interface{}{
		"revision": 1.0, "label": "drive-group-1", "infrastructureId": 1.0,
		"driveSizeMbDefault": 40960.0, "expandWithServerInstanceGroup": 1.0,
		"storageType": "iscsi_ssd", "updatedTimestamp": "2024-01-02T00:00:00Z",
		"deployType": "create", "deployStatus": "not_started",
	},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var siCmdInterface = map[string]interface{}{
	"id": 7.0, "revision": 2.0, "label": "iface-1", "infrastructureId": 1.0,
	"instanceId": 1.0, "index": 0.0, "capacityMbps": 10000.0, "dirtyBit": false,
	"serviceStatus": "active", "networkId": 9.0,
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var sigCmdInterface = map[string]interface{}{
	"id": 7.0, "revision": 2.0, "label": "sig-iface-1", "infrastructureId": 1.0,
	"groupId": 1.0, "index": 0.0, "serviceStatus": "active", "networkId": 9.0,
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var sigCmdNetworkEndpointGroup = map[string]interface{}{
	"id": "77", "name": "sig-1-neg", "siteId": 1.0, "revision": "2",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var sigCmdACL = map[string]interface{}{
	"id": "3", "ruleType": "ipv4", "direction": "in", "sequence": 10.0,
	"forwardingAction": "allow", "enforcementPoint": "svi",
	"endpointGroupId": "77", "logicalNetworkId": "7",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var siCmdInterfaceConfig = map[string]interface{}{
	"revision": 5.0, "label": "iface-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"instanceId": 1.0, "index": 0.0, "capacityMbps": 10000.0,
	"deployType": "create", "deployStatus": "not_started",
}

var sigCmdGroupConfig = map[string]interface{}{
	"revision": 4.0, "label": "sig-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"instanceCount": 1.0, "defaultServerTypeId": 5.0,
	"ipAllocateAuto": 1.0, "ipv4SubnetCreateAuto": 1.0,
	"processorCount": 2.0, "processorCoreCount": 8.0, "processorCoreMhz": 2400.0,
	"diskCount": 2.0, "diskSizeMbytes": 102400.0, "diskTypes": []interface{}{},
	"virtualInterfacesEnabled": 0.0,
	"deployType":               "create", "deployStatus": "not_started",
}

var siCmdStatistics = map[string]interface{}{
	"serverStatus": map[string]interface{}{"available": 3.0},
	"site":         map[string]interface{}{"dc1": 4.0},
}

// siHandle registers a handler that answers GET with get and every other verb
// with write.
func siHandle(mux *http.ServeMux, path string, get interface{}, write interface{}) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			jsonResponse(w, http.StatusOK, get)
			return
		}
		if write == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		jsonResponse(w, http.StatusOK, write)
	})
}

// newSIFullTestServer serves every route the extended server-instance and
// server-instance-group commands touch.
func newSIFullTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		siHandle(mux, "/api/v2/server-instances", paginatedList(serverInstanceItem), nil)
		siHandle(mux, "/api/v2/infrastructures/1/server-instances", paginatedList(serverInstanceItem), serverInstanceItem)
		siHandle(mux, "/api/v2/infrastructures/1/server-instance-groups", paginatedList(serverInstanceGroupItem), serverInstanceGroupItem)
		siHandle(mux, "/api/v2/server-instances/statistics", siCmdStatistics, nil)
		siHandle(mux, "/api/v2/server-instances/1", serverInstanceItem, nil)
		siHandle(mux, "/api/v2/server-instances/1/config", siCmdInstanceConfig, siCmdInstanceConfig)
		siHandle(mux, "/api/v2/server-instances/1/meta", nil, nil)
		siHandle(mux, "/api/v2/server-instances/1/actions/reset", nil, nil)
		siHandle(mux, "/api/v2/server-instances/1/drives", map[string]interface{}{"data": []interface{}{siCmdDrive}}, nil)
		siHandle(mux, "/api/v2/server-instances/1/interfaces", paginatedList(siCmdInterface), nil)
		siHandle(mux, "/api/v2/server-instances/1/interfaces/7", siCmdInterface, nil)
		siHandle(mux, "/api/v2/server-instances/1/interfaces/7/config", siCmdInterfaceConfig, siCmdInterfaceConfig)
		siHandle(mux, "/api/v2/infrastructures/1/actions/power-get", nil, map[string]interface{}{"1": "on"})
		siHandle(mux, "/api/v2/infrastructures/1/actions/power-set", nil, nil)
		siHandle(mux, "/api/v2/server-instance-groups/1", serverInstanceGroupItem, nil)
		siHandle(mux, "/api/v2/server-instance-groups/1/meta", nil, nil)
		siHandle(mux, "/api/v2/server-instance-groups/1/config", sigCmdGroupConfig, sigCmdGroupConfig)
		siHandle(mux, "/api/v2/server-instance-groups/1/drive-groups", map[string]interface{}{"data": []interface{}{sigCmdDriveGroup}}, nil)
		siHandle(mux, "/api/v2/server-instance-groups/1/interfaces", paginatedList(sigCmdInterface), nil)
		siHandle(mux, "/api/v2/server-instance-groups/1/interfaces/7", sigCmdInterface, nil)
		siHandle(mux, "/api/v2/server-instance-groups/1/config/networking", sigCmdNetworkEndpointGroup, sigCmdNetworkEndpointGroup)
		siHandle(mux, "/api/v2/server-instance-groups/1/config/networking/connections/5/security/rules", []interface{}{sigCmdACL}, sigCmdACL)
		siHandle(mux, "/api/v2/server-instance-groups/1/config/networking/connections/5/security/rules/3", sigCmdACL, sigCmdACL)
	})
	return httptest.NewServer(mux)
}

func TestServerInstanceListGlobal(t *testing.T) {
	srv := newSIFullTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-instance", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"label":"si-1"`) {
		t.Errorf("global list output missing instance: %s", out)
	}
}

func TestServerInstanceReadOnlySubcommands(t *testing.T) {
	srv := newSIFullTestServer()
	defer srv.Close()

	cases := [][]string{
		{"server-instance", "config-example"},
		{"server-instance", "statistics"},
		{"server-instance", "drives", "1"},
		{"server-instance", "interfaces", "1"},
		{"server-instance", "interface", "1", "7"},
		{"server-instance", "power-status-batch", "1", "1"},
		{"server-instance-group", "drive-groups", "1"},
		{"server-instance-group", "interfaces", "1"},
		{"server-instance-group", "interface", "1", "7"},
		{"server-instance-group", "network", "config", "1"},
		{"server-instance-group", "network", "acl", "list", "1", "5"},
		{"server-instance-group", "network", "acl", "get", "1", "5", "3"},
		{"server-instance-group", "network", "acl", "config-example"},
	}

	for _, args := range cases {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			out, err := runCLI(t, srv, args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out == "" {
				t.Error("expected output")
			}
		})
	}
}

func TestServerInstanceWriteSubcommands(t *testing.T) {
	srv := newSIFullTestServer()
	defer srv.Close()

	dir := t.TempDir()
	configFile := dir + "/config.json"
	if err := os.WriteFile(configFile, []byte(`{"label":"si-renamed"}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	metaFile := dir + "/meta.json"
	if err := os.WriteFile(metaFile, []byte(`{"tags":["prod"]}`), 0o600); err != nil {
		t.Fatalf("write meta: %v", err)
	}

	cases := [][]string{
		{"server-instance", "create", "1", "--label", "si-new"},
		{"server-instance", "create", "1", "--config-source", configFile},
		{"server-instance", "delete", "1"},
		{"server-instance", "update-config", "1", "--config-source", configFile},
		{"server-instance", "update-meta", "1", "--config-source", metaFile},
		{"server-instance", "reset", "1"},
		{"server-instance", "power-set-batch", "1", "off", "1"},
		{"server-instance", "update-interface-config", "1", "7", "--config-source", configFile},
		{"server-instance-group", "update-meta", "1", "--config-source", metaFile},
		{"server-instance-group", "network", "replace", "1"},
		{"server-instance-group", "network", "replace", "1", "--config-source", configFile},
		{"server-instance-group", "network", "acl", "add", "1", "5", "--rule-type", "ipv4", "--sequence", "10"},
		{"server-instance-group", "network", "acl", "update", "1", "5", "3", "--config-source", configFile},
		{"server-instance-group", "network", "acl", "remove", "1", "5", "3"},
	}

	for _, args := range cases {
		t.Run(strings.Join(args[:3], " "), func(t *testing.T) {
			if _, err := runCLI(t, srv, args...); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestServerInstanceCreateRequiresLabelOrConfigSource(t *testing.T) {
	srv := newSIFullTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-instance", "create", "1"); err == nil {
		t.Error("create without --label or --config-source should fail")
	}
	if _, err := runCLI(t, srv, "server-instance", "create", "1", "--label", "x", "--config-source", "pipe"); err == nil {
		t.Error("--label and --config-source should be mutually exclusive")
	}
}

func TestServerInstancePowerSetBatchRejectsUnknownAction(t *testing.T) {
	srv := newSIFullTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-instance", "power-set-batch", "1", "explode", "1"); err == nil {
		t.Error("an unknown power command should be rejected")
	}
}
