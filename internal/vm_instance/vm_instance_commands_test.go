package vm_instance

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// vm_instance_commands_test.go covers the VM instance commands added on top of
// the original get/list/config/power set: config-example, create, delete,
// update-config, update-meta, apply-type, variables and os-installation-data.

func TestVMInstanceConfigExample(t *testing.T) {
	if err := VMInstanceConfigExample(t.Context()); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// VMInstanceCreate

func TestVMInstanceCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                  infraHandler(),
		"/api/v2/infrastructures/123/vm-instances": testutils.JSONHandler(201, makeVMInstance(2)),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceCreate(ctx, "123", []byte(`{"typeId":5,"groupId":10}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceCreate_InvalidConfig(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceCreate(ctx, "123", []byte(`not a config`)); err == nil {
		t.Fatal("expected error for invalid configuration, got nil")
	}
}

func TestVMInstanceCreateFromFlags_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                  infraHandler(),
		"/api/v2/infrastructures/123/vm-instances": testutils.JSONHandler(201, makeVMInstance(2)),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceCreateFromFlags(ctx, "123", "5", "10", "40", []string{"tag1"}); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceCreateFromFlags_InvalidTypeId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceCreateFromFlags(ctx, "123", "bad", "10", "", nil); err == nil {
		t.Fatal("expected error for invalid VM type id, got nil")
	}
}

// VMInstanceDelete

func TestVMInstanceDelete_HappyPath(t *testing.T) {
	// The revision is read first with a GET on the path the DELETE then uses.
	getHandler := testutils.JSONHandler(200, makeVMInstance(1))
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				if r.Header.Get("If-Match") != "1" {
					t.Errorf("expected If-Match 1 on delete, got %q", r.Header.Get("If-Match"))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			getHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceDelete(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceDelete_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                    infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceDelete(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

// VMInstanceUpdateConfig

func TestVMInstanceUpdateConfig_HappyPath(t *testing.T) {
	config := map[string]any{
		"revision":         float64(7),
		"label":            "vm-1",
		"typeId":           float64(5),
		"diskSizeGB":       float64(100),
		"ramGB":            float64(16),
		"cpuCores":         float64(8),
		"deployType":       "deploy",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	configHandler := testutils.JSONHandler(200, config)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1/config": func(w http.ResponseWriter, r *http.Request) {
			// The config sub-resource is guarded by the config revision, not by
			// the VM instance revision.
			if r.Method == http.MethodPatch && r.Header.Get("If-Match") != "7" {
				t.Errorf("expected If-Match 7 on update, got %q", r.Header.Get("If-Match"))
			}
			configHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceUpdateConfig(ctx, "123", "1", []byte(`{"label":"vm-1"}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceUpdateConfig_InvalidConfig(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceUpdateConfig(ctx, "123", "1", []byte(`not a config`)); err == nil {
		t.Fatal("expected error for invalid configuration, got nil")
	}
}

// VMInstanceUpdateMeta

func TestVMInstanceUpdateMeta_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                         infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1/meta": testutils.JSONHandler(200, makeVMInstance(1)),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceUpdateMeta(ctx, "123", "1", []byte(`{"tags":["prod"]}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceUpdateMeta_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                         infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1/meta": testutils.ErrorHandler(500, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceUpdateMeta(ctx, "123", "1", []byte(`{"tags":["prod"]}`)); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// VMInstanceApplyType

func TestVMInstanceApplyType_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                    infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1": testutils.JSONHandler(200, makeVMInstance(1)),
		"/api/v2/infrastructures/123/vm-instances/1/actions/apply-type/5": func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("If-Match") != "1" {
				t.Errorf("expected If-Match 1, got %q", r.Header.Get("If-Match"))
			}
			testutils.JSONHandler(200, makeVMInstance(1))(w, r)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceApplyType(ctx, "123", "1", "5"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceApplyType_InvalidTypeId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceApplyType(ctx, "123", "1", "bad"); err == nil {
		t.Fatal("expected error for invalid VM type id, got nil")
	}
}

// VMInstanceVariables / VMInstanceOSInstallationData
//
// Both endpoints answer with "server": null and a 'network' object, which the
// strict SDK models reject, so the CLI reads the bodies raw.

const vmInstanceContextBody = `{
	"site": {"id": 1},
	"siteConfig": {},
	"server": null,
	"serverInstance": {"id": 1},
	"serverInstanceGroup": {"id": 10},
	"infrastructure": {"id": 123},
	"driveGroups": [],
	"drives": [],
	"fileShares": [],
	"buckets": [],
	"sharedDrives": [],
	"network": {"interfaces": []},
	"variables": {},
	"secrets": {},
	"userSSHKeys": [],
	"managementSSHKey": null
}`

func TestVMInstanceVariables_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1/variables": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("usage") != "AnsibleBundle" {
				t.Errorf("expected usage=AnsibleBundle, got %q", r.URL.RawQuery)
			}
			testutils.RawHandler(200, vmInstanceContextBody)(w, r)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceVariables(ctx, "123", "1", "AnsibleBundle"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceVariables_InvalidUsage(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceVariables(ctx, "123", "1", "NoSuchUsage"); err == nil {
		t.Fatal("expected error for invalid usage, got nil")
	}
}

func TestVMInstanceVariables_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                              infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1/variables": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceVariables(ctx, "123", "1", ""); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestVMInstanceOSInstallationData_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1/os-installation-data": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("removeEmpty") != "1" {
				t.Errorf("expected removeEmpty=1, got %q", r.URL.RawQuery)
			}
			testutils.RawHandler(200, vmInstanceContextBody)(w, r)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceOSInstallationData(ctx, "123", "1", "", true); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceOSInstallationData_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instances/1/os-installation-data": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceOSInstallationData(ctx, "123", "1", "", false); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}
