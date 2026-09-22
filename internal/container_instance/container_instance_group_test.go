package container_instance

import (
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

var containerInstanceGroupConfigItem = map[string]any{
	"revision": 4, "label": "cig-1", "instanceCount": 2,
	"deployType": "create", "deployStatus": "not_started",
	"updatedTimestamp": ts0,
}

var containerInstanceGroupItem = map[string]any{
	"label": "cig-1", "updatedTimestamp": ts0, "id": 77, "revision": 11,
	"infrastructureId": 1234, "infrastructure": map[string]any{"id": 1234},
	"serviceStatus": "active", "instanceGroupType": "container",
	"diskSizeGB": 40, "instanceCount": 2, "vmPoolId": 3,
	"createdTimestamp": ts0,
	"config":           containerInstanceGroupConfigItem,
	"meta":             map[string]any{},
}

// The live API carries the interface label only inside `config` (the SDK model
// requires a top-level `label`), so the fixture mirrors that shape.
var containerInstanceGroupInterfaceItem = map[string]any{
	"index": 0, "updatedTimestamp": ts0, "id": 3,
	"revision": 1, "serviceStatus": "active", "groupId": 77,
	"infrastructureId": 1234, "networkId": nil,
	"config": map[string]any{
		"revision": 1, "label": "if-1", "index": 0,
		"deployType": "create", "deployStatus": "not_started",
		"updatedTimestamp": ts0,
	},
	"meta":             map[string]any{},
	"createdTimestamp": ts0,
}

var networkConnectionItem = map[string]any{
	"id": "5", "tagged": true, "accessMode": "l2", "mtu": 1500,
	"providesDefaultRoute": true, "disableAutoIpAllocation": false,
}

var networkEndpointGroupItem = map[string]any{
	"id": "21", "name": "cig-1-neg", "revision": "1",
	"createdTimestamp": ts0, "updatedTimestamp": ts0,
}

func TestContainerInstanceGroupList(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupList(ctx, "test-infra") })

	if !strings.Contains(out, "cig-1") {
		t.Fatalf("expected the group label in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instance-groups" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstanceGroupGet(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupGet(ctx, "1234", "77") })

	if !strings.Contains(out, "cig-1") {
		t.Fatalf("expected the group label in the output, got %s", out)
	}
}

func TestContainerInstanceGroupGetInvalidId(t *testing.T) {
	if err := ContainerInstanceGroupGet(nil, "1234", "nope"); err == nil {
		t.Fatal("expected an error for a non numeric group ID")
	}
}

func TestContainerInstanceGroupConfigExample(t *testing.T) {
	out := run(t, func() error { return ContainerInstanceGroupConfigExample(nil) })

	if !strings.Contains(out, "osTemplateId") {
		t.Fatalf("expected the example to contain osTemplateId, got %s", out)
	}
}

func TestContainerInstanceGroupCreateFromConfig(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceGroupCreate(ctx, "test-infra",
			[]byte(`{"typeId":42,"osTemplateId":7,"vmPoolId":3,"diskSizeGB":40,"instanceCount":2}`))
	})

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/infrastructures/1234/container-instance-groups" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"osTemplateId":7`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerInstanceGroupCreateFromFlags(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceGroupCreateFromFlags(ctx, "1234", "42", "7", "3", "40", "2", []string{"prod"})
	})

	got := rec.last()
	if !strings.Contains(got.Body, `"vmPoolId":3`) || !strings.Contains(got.Body, `"instanceCount":2`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerInstanceGroupCreateFromFlagsInvalidNumber(t *testing.T) {
	if err := ContainerInstanceGroupCreateFromFlags(nil, "1234", "abc", "7", "3", "40", "2", nil); err == nil {
		t.Fatal("expected an error for a non numeric container type ID")
	}
}

func TestContainerInstanceGroupDeleteSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerInstanceGroupDelete(ctx, "1234", "77"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := rec.last()
	if got.Method != http.MethodDelete {
		t.Fatalf("unexpected request method %s", got.Method)
	}
	if got.Header.Get("If-Match") != "11" {
		t.Fatalf("expected the group revision as If-Match, got %q", got.Header.Get("If-Match"))
	}
}

func TestContainerInstanceGroupGetConfig(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupGetConfig(ctx, "1234", "77") })

	if !strings.Contains(out, "cig-1") {
		t.Fatalf("expected the configuration label in the output, got %s", out)
	}
}

func TestContainerInstanceGroupUpdateConfigUsesConfigRevision(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceGroupUpdateConfig(ctx, "1234", "77", []byte(`{"label":"renamed"}`))
	})

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/config" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	// The config sub-resource is guarded by the config revision (4), not by the
	// group revision (11).
	if got.Header.Get("If-Match") != "4" {
		t.Fatalf("expected the config revision as If-Match, got %q", got.Header.Get("If-Match"))
	}
}

func TestContainerInstanceGroupUpdateMeta(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceGroupUpdateMeta(ctx, "1234", "77", []byte(`{"tags":["prod"]}`))
	})

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/meta" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
}

func TestContainerInstanceGroupInstances(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupInstances(ctx, "1234", "77") })

	if !strings.Contains(out, "ci-1") {
		t.Fatalf("expected the instance label in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/container-instances" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstanceGroupInterfaces(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupInterfaces(ctx, "1234", "77") })

	if !strings.Contains(out, "if-1") {
		t.Fatalf("expected the interface label in the output, got %s", out)
	}
}

func TestContainerInstanceGroupInterfaceGet(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupInterfaceGet(ctx, "1234", "77", "3") })

	if !strings.Contains(out, "if-1") {
		t.Fatalf("expected the interface label in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/interfaces/3" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstanceGroupApplyTypeSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error { return ContainerInstanceGroupApplyType(ctx, "1234", "77", "42") })

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/actions/apply-type/42" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if got.Header.Get("If-Match") != "11" {
		t.Fatalf("expected the group revision as If-Match, got %q", got.Header.Get("If-Match"))
	}
}

func TestContainerInstanceGroupNetworkConfiguration(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupNetworkConfiguration(ctx, "1234", "77") })

	if !strings.Contains(out, "cig-1-neg") {
		t.Fatalf("expected the endpoint group name in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/config/networking" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstanceGroupNetworkList(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupNetworkList(ctx, "1234", "77") })

	if !strings.Contains(out, "l2") {
		t.Fatalf("expected the access mode in the output, got %s", out)
	}
	if got := rec.last(); got.Method != http.MethodGet || got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/config/networking/connections" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
}

func TestContainerInstanceGroupNetworkGet(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGroupNetworkGet(ctx, "1234", "77", "5") })

	if !strings.Contains(out, "l2") {
		t.Fatalf("expected the access mode in the output, got %s", out)
	}
}

func TestContainerInstanceGroupNetworkConnect(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceGroupNetworkConnect(ctx, "1234", "77", "9", "l2", "true", "active-backup")
	})

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/config/networking/connections" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"logicalNetworkId":"9"`) || !strings.Contains(got.Body, `"accessMode":"l2"`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
	if !strings.Contains(got.Body, `"mode":"active-backup"`) {
		t.Fatalf("expected the redundancy config in the body, got %s", got.Body)
	}
}

func TestContainerInstanceGroupNetworkConnectInvalidTagged(t *testing.T) {
	if err := ContainerInstanceGroupNetworkConnect(nil, "1234", "77", "9", "l2", "maybe", ""); err == nil {
		t.Fatal("expected an error for an invalid tagged value")
	}
}

func TestContainerInstanceGroupNetworkUpdate(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceGroupNetworkUpdate(ctx, "1234", "77", "5", "l2", "false", "")
	})

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/config/networking/connections/5" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"tagged":false`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerInstanceGroupNetworkDisconnect(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerInstanceGroupNetworkDisconnect(ctx, "1234", "77", "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := rec.last()
	if got.Method != http.MethodDelete || got.Path != "/api/v2/infrastructures/1234/container-instance-groups/77/config/networking/connections/5" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
}
