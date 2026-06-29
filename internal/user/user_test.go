package user

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// userCoreFields returns the set of fields required by User and UserConfiguration.
func userCoreFields() map[string]any {
	return map[string]any{
		"displayName": "Alice", "email": "alice@example.com",
		"emailStatus": "active", "language": "en",
		"brand": "default", "accessLevel": "admin",
		"isArchived": false, "isBlocked": false,
		"isBillable": false, "isTestingMode": false,
		"isSuspended": false, "authenticatorEnabled": false,
		"passwordChangeRequired": false, "authenticatorMustChange": false,
		"authenticatorCreatedTimestamp": "",
		"excludeFromReports": false, "isTestAccount": false,
		"isDatastorePublisher": false, "isBrandManager": false,
		"provider": "local", "franchise": "default",
		"planType": "default", "lastLoginType": "password",
		"lastLoginTimestamp": "2024-01-01T00:00:00Z",
		"passwordLastChangedTimestamp": "2024-01-01T00:00:00Z",
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"revision": float64(1),
	}
}

func makeUser(id int) map[string]any {
	u := userCoreFields()
	u["id"] = id
	// UserConfiguration requires the same fields as User.
	cfg := userCoreFields()
	cfg["id"] = id
	u["config"] = cfg
	u["meta"] = map[string]any{}
	u["links"] = []any{}
	return u
}

func TestList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeUser(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := List(ctx, false, "", "", "", "", "", "", "", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := List(ctx, false, "", "", "", "", "", "", "", ""); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := List(ctx, false, "", "", "", "", "", "", "", ""); err != nil {
		t.Fatalf("expected nil error on empty list, got: %v", err)
	}
}

func TestList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeUser(i + 1)
		page2[i] = makeUser(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeUser(i + 201)
	}

	ts := testutils.MultiPageServer("/api/v2/users", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := List(ctx, false, "", "", "", "", "", "", "", ""); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1": testutils.JSONHandler(200, makeUser(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Get(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Get(ctx, "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Get(ctx, "not-a-number"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

// --- Create ---

func TestCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users": testutils.JSONHandler(201, makeUser(42)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"email":"bob@example.com","displayName":"Bob","accessLevel":"user","password":"secret"}`)
	if err := Create(ctx, config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCreate_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"email":"bob@example.com","displayName":"Bob","accessLevel":"user","password":"secret"}`)
	if err := Create(ctx, config); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestCreate_InvalidConfig(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Create(ctx, []byte("not-json")); err == nil {
		t.Fatal("expected error on invalid config, got nil")
	}
}

// --- Archive / Unarchive ---

// makeUserRoute registers a GET /api/v2/users/{id} returning revision 1 and
// an action route, to satisfy the getUserIdAndRevision fetch.
func makeUserAndActionRoutes(userId, actionPath string, actionHandler http.HandlerFunc) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/api/v2/users/" + userId: testutils.JSONHandler(200, makeUser(1)),
		actionPath:                actionHandler,
	}
}

func TestArchive_HappyPath(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/archive",
		testutils.JSONHandler(200, makeUser(1)))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Archive(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestArchive_500(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/archive",
		testutils.ErrorHandler(500, "internal server error"))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Archive(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestArchive_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Archive(ctx, "bad"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

func TestUnarchive_HappyPath(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/unarchive",
		testutils.JSONHandler(200, makeUser(1)))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Unarchive(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestUnarchive_500(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/unarchive",
		testutils.ErrorHandler(500, "internal server error"))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Unarchive(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// --- Suspend / Unsuspend ---

func makeSuspendReason() map[string]any {
	return map[string]any{
		"id":               float64(1),
		"userId":           float64(1),
		"type":             "admin",
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"publicComment":    "suspended for testing",
	}
}

func TestSuspend_HappyPath(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/suspend",
		testutils.JSONHandler(200, makeSuspendReason()))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Suspend(ctx, "1", "testing"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestSuspend_500(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/suspend",
		testutils.ErrorHandler(500, "internal server error"))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Suspend(ctx, "1", "testing"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestSuspend_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Suspend(ctx, "bad", "reason"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

func TestUnsuspend_HappyPath(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/unsuspend",
		testutils.JSONHandler(200, nil))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Unsuspend(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestUnsuspend_500(t *testing.T) {
	routes := makeUserAndActionRoutes("1", "/api/v2/users/1/actions/unsuspend",
		testutils.ErrorHandler(500, "internal server error"))
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Unsuspend(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// --- GetLimits / UpdateLimits ---

// makeUserLimits returns a complete QuotaLimitsBreakdown payload. The
// "effective" QuotaProfileLimits object must include every SDK-required field
// or typed unmarshalling fails.
func makeUserLimits() map[string]any {
	effective := map[string]any{
		"infrastructureServerGroupMaxCount":     float64(10),
		"infrastructureDriveMaxCount":           float64(10),
		"infrastructureFileShareMaxCount":       float64(10),
		"infrastructureBucketMaxCount":          float64(10),
		"infrastructureVmInstanceGroupMaxCount": float64(10),
		"serverGroupInstancesMaxCount":          float64(10),
		"serverGroupInstancesMinCount":          float64(1),
		"vmInstanceGroupVmInstancesMaxCount":    float64(10),
		"vmInstanceMaxDiskSizeMbytes":           float64(1048576),
		"driveMaxSizeMbytes":                    float64(1048576),
		"driveMinSizeMbytes":                    float64(1024),
		"fileShareMinSizeGb":                    float64(1),
		"fileShareMaxSizeGb":                    float64(1024),
		"bucketMinSizeGb":                       float64(1),
		"bucketMaxSizeGb":                       float64(1024),
		"showOperatingSystemImagesTab":          true,
		"showTemplateAssetsView":                true,
		"userResourceServerTypeNameToMaxCount":  map[string]any{},
		"userSshKeysCountMax":                   float64(10),
		"showLegacyPages":                       false,
		"showEliChatBot":                        false,
		"enableCustomRaidConfiguration":         true,
		"enableInfrastructureVmInstance":        true,
		"enableInfrastructureExtensions":        true,
		"allowedInfrastructureExtensions":       []any{},
		"allowedServerTypes":                    []any{},
		"allowedSites":                          []any{},
		"allowedLogicalNetworkProfiles":         []any{},
		"allowedPreCreatedLogicalNetworks":      []any{},
	}
	return map[string]any{
		"effective": effective,
	}
}

func TestGetLimits_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/quota-limits-breakdown": testutils.JSONHandler(200, makeUserLimits()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetLimits(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGetLimits_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/quota-limits-breakdown": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetLimits(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestGetLimits_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetLimits(ctx, "bad"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

// --- SetPassword ---

func TestSetPassword_HappyPath(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/users/1":                        testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/set-password":   testutils.JSONHandler(200, makeUser(1)),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SetPassword(ctx, "1", "newpassword"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestSetPassword_500(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/users/1":                       testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/set-password":  testutils.ErrorHandler(500, "internal server error"),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SetPassword(ctx, "1", "newpassword"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestSetPassword_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SetPassword(ctx, "bad", "password"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

// --- ChangeAccount ---

func TestChangeAccount_HappyPath(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/users/1":                         testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/change-account":  testutils.JSONHandler(200, makeUser(1)),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ChangeAccount(ctx, "1", 99); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestChangeAccount_500(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/users/1":                        testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/change-account": testutils.ErrorHandler(500, "internal server error"),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ChangeAccount(ctx, "1", 99); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestChangeAccount_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ChangeAccount(ctx, "bad", 99); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

// --- GetSSHKeys / AddSSHKey / DeleteSSHKey ---

func makeSSHKeysList() map[string]any {
	return map[string]any{
		"data": []any{
			map[string]any{
				"id":               float64(1),
				"userId":           float64(1),
				"sshKey":           "ssh-rsa AAAA...",
				"status":           "active",
				"createdTimestamp": "2024-01-01T00:00:00Z",
			},
		},
	}
}

func makeSSHKey() map[string]any {
	return map[string]any{
		"id":               float64(1),
		"userId":           float64(1),
		"sshKey":           "ssh-rsa AAAA...",
		"status":           "active",
		"createdTimestamp": "2024-01-01T00:00:00Z",
	}
}

func TestGetSSHKeys_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/ssh-keys": testutils.JSONHandler(200, makeSSHKeysList()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetSSHKeys(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGetSSHKeys_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/ssh-keys": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetSSHKeys(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestGetSSHKeys_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetSSHKeys(ctx, "bad"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

func TestAddSSHKey_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/ssh-keys": testutils.JSONHandler(201, makeSSHKey()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AddSSHKey(ctx, "1", "ssh-rsa AAAA..."); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestAddSSHKey_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/ssh-keys": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AddSSHKey(ctx, "1", "ssh-rsa AAAA..."); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestAddSSHKey_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AddSSHKey(ctx, "bad", "ssh-rsa AAAA..."); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

func TestDeleteSSHKey_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/ssh-keys/5": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeleteSSHKey(ctx, "1", "5"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDeleteSSHKey_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/ssh-keys/5": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeleteSSHKey(ctx, "1", "5"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestDeleteSSHKey_InvalidUserId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeleteSSHKey(ctx, "bad", "5"); err == nil {
		t.Fatal("expected error on invalid user ID, got nil")
	}
}

func TestDeleteSSHKey_InvalidKeyId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeleteSSHKey(ctx, "1", "bad"); err == nil {
		t.Fatal("expected error on invalid key ID, got nil")
	}
}
