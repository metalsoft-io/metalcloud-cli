package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const ctTs = "2024-01-01T00:00:00Z"

var ctItem = map[string]interface{}{
	"id": 42, "name": "small", "label": "small-ct", "displayName": "Small",
	"cpuCores": 4, "ramGB": 16,
}

var ctContainerItem = map[string]interface{}{
	"id": 100, "name": "container-100", "siteId": 1, "infrastructureId": 1234,
	"userId": 7, "instanceId": 5678, "containerInstanceId": 5678,
	"host": "host-1", "hosts": "host-1",
	"cpuCores": 4, "ramGB": 16, "diskSizeGB": 40,
	"typeId": 42, "poolId": 3, "administrationState": "active",
	"powerState": "on", "powerStateLastUpdatedTimestamp": ctTs,
	"createdTimestamp": ctTs, "allocationTimestamp": ctTs,
	"disks": []interface{}{},
}

var ctInfraItem = map[string]interface{}{
	"id": 1234, "label": "test-infra", "revision": 1, "serviceStatus": "active",
	"datacenterName": "dc1", "siteId": 1, "designIsLocked": 0,
	"createdTimestamp": ctTs, "updatedTimestamp": ctTs,
	"config": map[string]interface{}{},
}

var ciConfigItem = map[string]interface{}{
	"revision": 3, "label": "ci-1", "typeId": 42,
	"deployType": "create", "deployStatus": "not_started",
	"diskSizeGB": 40, "ramGB": 16, "cpuCores": 4,
	"updatedTimestamp": ctTs,
}

var ciItem = map[string]interface{}{
	"label": "ci-1", "typeId": 42, "diskSizeGB": 40, "ramGB": 16, "cpuCores": 4,
	"updatedTimestamp": ctTs, "id": 5678, "revision": 9,
	"groupId": 77, "infrastructureId": 1234,
	"infrastructure": map[string]interface{}{"id": 1234},
	"serviceStatus":  "active", "instanceType": "container",
	"createdTimestamp": ctTs,
	"config":           ciConfigItem,
	"meta":             map[string]interface{}{},
}

var cigConfigItem = map[string]interface{}{
	"revision": 4, "label": "cig-1", "instanceCount": 2,
	"deployType": "create", "deployStatus": "not_started",
	"updatedTimestamp": ctTs,
}

var cigItem = map[string]interface{}{
	"label": "cig-1", "updatedTimestamp": ctTs, "id": 77, "revision": 11,
	"infrastructureId": 1234, "infrastructure": map[string]interface{}{"id": 1234},
	"serviceStatus": "active", "instanceGroupType": "container",
	"diskSizeGB": 40, "instanceCount": 2,
	"createdTimestamp": ctTs,
	"config":           cigConfigItem,
	"meta":             map[string]interface{}{},
}

var cigConnectionItem = map[string]interface{}{
	"id": "5", "tagged": true, "accessMode": "l2",
}

// containerTestServer registers every route exercised by the container command
// tests below.
func containerTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	base := "/api/v2/infrastructures/1234"
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/container-types", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusCreated, ctItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(ctItem))
		})
		mux.HandleFunc("/api/v2/container-types/42", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, ctItem)
		})
		mux.HandleFunc("/api/v2/container-types/42/containers", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(ctContainerItem))
		})

		mux.HandleFunc("/api/v2/containers", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(ctContainerItem))
		})
		mux.HandleFunc("/api/v2/containers/100", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ctContainerItem)
		})
		mux.HandleFunc("/api/v2/containers/100/power-status", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, "on")
		})
		mux.HandleFunc("/api/v2/containers/100/reboot", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(ctInfraItem))
		})
		mux.HandleFunc(base+"/container-instances", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(ciItem))
		})
		mux.HandleFunc(base+"/container-instances/5678", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ciItem)
		})
		mux.HandleFunc(base+"/container-instances/5678/config", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ciConfigItem)
		})
		mux.HandleFunc(base+"/container-instances/5678/actions/apply-type/42", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ciItem)
		})

		mux.HandleFunc(base+"/container-instance-groups", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(cigItem))
		})
		mux.HandleFunc(base+"/container-instance-groups/77", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, cigItem)
		})
		mux.HandleFunc(base+"/container-instance-groups/77/config/networking/connections", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"data": []interface{}{cigConnectionItem}})
		})
	})
	return httptest.NewServer(mux)
}

func TestContainerTypeCommands(t *testing.T) {
	srv := containerTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "container-type", "list")
	if err != nil {
		t.Fatalf("container-type list failed: %v", err)
	}
	if !strings.Contains(out, "small-ct") {
		t.Fatalf("expected the container type label in the output, got %s", out)
	}

	out, err = runCLI(t, srv, "ct", "get", "42")
	if err != nil {
		t.Fatalf("container-type get failed: %v", err)
	}
	if !strings.Contains(out, "Small") {
		t.Fatalf("expected the display name in the output, got %s", out)
	}

	out, err = runCLI(t, srv, "container-type", "containers", "42")
	if err != nil {
		t.Fatalf("container-type containers failed: %v", err)
	}
	if !strings.Contains(out, "container-100") {
		t.Fatalf("expected the container name in the output, got %s", out)
	}

	if _, err := runCLI(t, srv, "container-type", "delete", "42"); err != nil {
		t.Fatalf("container-type delete failed: %v", err)
	}

	out, err = runCLI(t, srv, "container-type", "config-example")
	if err != nil {
		t.Fatalf("container-type config-example failed: %v", err)
	}
	if !strings.Contains(out, "cpuCores") {
		t.Fatalf("expected the example to contain cpuCores, got %s", out)
	}
}

func TestContainerTypeCreateRequiresConfigSource(t *testing.T) {
	srv := containerTestServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "container-type", "create"); err == nil {
		t.Fatal("expected an error when --config-source is missing")
	}
}

func TestContainerCommands(t *testing.T) {
	srv := containerTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "container", "list")
	if err != nil {
		t.Fatalf("container list failed: %v", err)
	}
	if !strings.Contains(out, "container-100") {
		t.Fatalf("expected the container name in the output, got %s", out)
	}

	if _, err := runCLI(t, srv, "containers", "get", "100"); err != nil {
		t.Fatalf("container get failed: %v", err)
	}

	out, err = runCLI(t, srv, "container", "power-status", "100")
	if err != nil {
		t.Fatalf("container power-status failed: %v", err)
	}
	if !strings.Contains(out, "on") {
		t.Fatalf("expected the power state in the output, got %s", out)
	}

	if _, err := runCLI(t, srv, "container", "reboot", "100"); err != nil {
		t.Fatalf("container reboot failed: %v", err)
	}
}

func TestContainerInstanceCommands(t *testing.T) {
	srv := containerTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "container-instance", "list", "test-infra")
	if err != nil {
		t.Fatalf("container-instance list failed: %v", err)
	}
	if !strings.Contains(out, "ci-1") {
		t.Fatalf("expected the instance label in the output, got %s", out)
	}

	if _, err := runCLI(t, srv, "ci", "get", "1234", "5678"); err != nil {
		t.Fatalf("container-instance get failed: %v", err)
	}

	if _, err := runCLI(t, srv, "container-instance", "config", "1234", "5678"); err != nil {
		t.Fatalf("container-instance config failed: %v", err)
	}

	if _, err := runCLI(t, srv, "container-instance", "apply-type", "1234", "5678", "42"); err != nil {
		t.Fatalf("container-instance apply-type failed: %v", err)
	}
}

func TestContainerInstanceCreateRequiresFlags(t *testing.T) {
	srv := containerTestServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "container-instance", "create", "1234"); err == nil {
		t.Fatal("expected an error when neither --config-source nor --container-type-id is given")
	}
}

func TestContainerInstanceGroupCommands(t *testing.T) {
	srv := containerTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "container-instance-group", "list", "test-infra")
	if err != nil {
		t.Fatalf("container-instance-group list failed: %v", err)
	}
	if !strings.Contains(out, "cig-1") {
		t.Fatalf("expected the group label in the output, got %s", out)
	}

	if _, err := runCLI(t, srv, "cig", "get", "1234", "77"); err != nil {
		t.Fatalf("container-instance-group get failed: %v", err)
	}

	out, err = runCLI(t, srv, "cig", "network", "list", "1234", "77")
	if err != nil {
		t.Fatalf("container-instance-group network list failed: %v", err)
	}
	if !strings.Contains(out, "l2") {
		t.Fatalf("expected the access mode in the output, got %s", out)
	}
}

func TestContainerInstanceGroupNetworkUpdateRequiresFlag(t *testing.T) {
	srv := containerTestServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "cig", "network", "update", "1234", "77", "5"); err == nil {
		t.Fatal("expected an error when no update flag is given")
	}
}
