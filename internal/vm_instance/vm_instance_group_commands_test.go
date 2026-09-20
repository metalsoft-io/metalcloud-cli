package vm_instance

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// vm_instance_group_commands_test.go covers the VM instance group commands
// added on top of the original list/get/create/update/delete/instances set:
// config, update-meta, apply-type, interfaces, interface and the network
// sub-commands.

func makeVMGroupConfig(revision int) map[string]any {
	return map[string]any{
		"revision":         float64(revision),
		"label":            "vmg-1",
		"instanceCount":    float64(2),
		"deployType":       "create",
		"deployStatus":     "finished",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
}

// makeVMGroupInterface mirrors the live payload: the label only exists inside
// 'config', which is why the interface endpoints are read raw.
func makeVMGroupInterface(id int) map[string]any {
	return map[string]any{
		"id":               float64(id),
		"revision":         float64(4),
		"serviceStatus":    "active",
		"groupId":          float64(1),
		"infrastructureId": float64(123),
		"index":            float64(id - 1),
		"networkId":        nil,
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"meta":             map[string]any{},
		"config": map[string]any{
			"revision":         float64(3),
			"label":            "if0",
			"index":            float64(id - 1),
			"networkId":        nil,
			"deployType":       "create",
			"deployStatus":     "finished",
			"updatedTimestamp": "2024-01-01T00:00:00Z",
		},
	}
}

var vmGroupNetworkEndpointGroup = map[string]any{
	"id":               "1019",
	"name":             "vmg-1",
	"revision":         "1",
	"siteId":           float64(1),
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

var vmGroupNetworkConnection = map[string]any{
	"id":                      "535",
	"tagged":                  true,
	"accessMode":              "l2",
	"mtu":                     float64(1500),
	"providesDefaultRoute":    false,
	"disableAutoIpAllocation": false,
}

// VMInstanceGroupGetConfig

func TestVMInstanceGroupGetConfig_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                                 infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config": testutils.JSONHandler(200, makeVMGroupConfig(3)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupGetConfig(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupGetConfig_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                                 infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupGetConfig(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

// VMInstanceGroupUpdate sends the CONFIG revision, not the group revision.

func TestVMInstanceGroupUpdate_SendsConfigRevision(t *testing.T) {
	configHandler := testutils.JSONHandler(200, makeVMGroupConfig(9))
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch && r.Header.Get("If-Match") != "9" {
				t.Errorf("expected If-Match 9 on update, got %q", r.Header.Get("If-Match"))
			}
			configHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupUpdate(ctx, "123", "1", "new-label", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// VMInstanceGroupUpdateMeta

func TestVMInstanceGroupUpdateMeta_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                               infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/meta": testutils.JSONHandler(200, makeVMGroup(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupUpdateMeta(ctx, "123", "1", []byte(`{"tags":["prod"]}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupUpdateMeta_InvalidConfig(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupUpdateMeta(ctx, "123", "1", []byte(`not a config`)); err == nil {
		t.Fatal("expected error for invalid configuration, got nil")
	}
}

// VMInstanceGroupApplyType

func TestVMInstanceGroupApplyType_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                          infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1": testutils.JSONHandler(200, makeVMGroup(1)),
		"/api/v2/infrastructures/123/vm-instance-groups/1/actions/apply-type/5": func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("If-Match") != "1" {
				t.Errorf("expected If-Match 1, got %q", r.Header.Get("If-Match"))
			}
			testutils.JSONHandler(200, makeVMGroup(1))(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupApplyType(ctx, "123", "1", "5"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupApplyType_InvalidTypeId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupApplyType(ctx, "123", "1", "bad"); err == nil {
		t.Fatal("expected error for invalid VM type id, got nil")
	}
}

// VMInstanceGroupInterfaces / VMInstanceGroupInterfaceGet

func TestVMInstanceGroupInterfaces_HappyPath(t *testing.T) {
	body := testutils.PaginatedResponse([]any{makeVMGroupInterface(1), makeVMGroupInterface(2)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                                     infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/interfaces": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupInterfaces(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupInterfaces_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                                     infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/interfaces": testutils.ErrorHandler(500, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupInterfaces(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMInstanceGroupInterfaceGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/interfaces/2": testutils.JSONHandler(200, makeVMGroupInterface(2)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupInterfaceGet(ctx, "123", "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupInterfaceGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupInterfaceGet(ctx, "123", "1", "bad"); err == nil {
		t.Fatal("expected error for invalid interface id, got nil")
	}
}

// Network configuration and connections

func TestVMInstanceGroupNetworkConfiguration_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking": testutils.JSONHandler(200, vmGroupNetworkEndpointGroup),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkConfiguration(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkConnections_HappyPath(t *testing.T) {
	body := map[string]any{"data": []any{vmGroupNetworkConnection}}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkConnections(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections/535": testutils.JSONHandler(200, vmGroupNetworkConnection),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkGet(ctx, "123", "1", "535"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkGet_InvalidConnectionId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkGet(ctx, "123", "1", "bad"); err == nil {
		t.Fatal("expected error for invalid connection id, got nil")
	}
}

func TestVMInstanceGroupNetworkConfigExample(t *testing.T) {
	if err := VMInstanceGroupNetworkConfigExample(t.Context()); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkConnect_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections": testutils.JSONHandler(201, vmGroupNetworkConnection),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"logicalNetworkId":"5","accessMode":"l2","tagged":true}`)
	if err := VMInstanceGroupNetworkConnect(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkConnect_InvalidConfig(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkConnect(ctx, "123", "1", []byte(`not a config`)); err == nil {
		t.Fatal("expected error for invalid configuration, got nil")
	}
}

func TestVMInstanceGroupNetworkConnectFromFlags_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections": testutils.JSONHandler(201, vmGroupNetworkConnection),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkConnectFromFlags(ctx, "123", "1", "5", "l2", "true", "active-backup", "1500"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkConnectFromFlags_InvalidTagged(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkConnectFromFlags(ctx, "123", "1", "5", "l2", "maybe", "", ""); err == nil {
		t.Fatal("expected error for invalid tagged value, got nil")
	}
}

func TestVMInstanceGroupNetworkUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections/535": testutils.JSONHandler(200, vmGroupNetworkConnection),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkUpdate(ctx, "123", "1", "535", []byte(`{"tagged":false}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkUpdateFromFlags_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections/535": testutils.JSONHandler(200, vmGroupNetworkConnection),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkUpdateFromFlags(ctx, "123", "1", "535", "l2", "true", "", "1500"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkUpdateFromFlags_InvalidMtu(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkUpdateFromFlags(ctx, "123", "1", "535", "", "", "", "huge"); err == nil {
		t.Fatal("expected error for invalid MTU, got nil")
	}
}

func TestVMInstanceGroupNetworkDisconnect_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections/535": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkDisconnect(ctx, "123", "1", "535"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupNetworkDisconnect_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config/networking/connections/535": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupNetworkDisconnect(ctx, "123", "1", "535"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

// VMInstanceGroupInstances

func TestVMInstanceGroupInstances_HappyPath(t *testing.T) {
	body := testutils.PaginatedResponse([]any{makeVMInstance(1)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1/vm-instances": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupInstances(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
