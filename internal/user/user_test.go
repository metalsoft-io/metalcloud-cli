package user

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
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
		"excludeFromReports":            false, "isTestAccount": false,
		"isDatastorePublisher": false, "isBrandManager": false,
		"provider": "local", "franchise": "default",
		"planType": "default", "lastLoginType": "password",
		"lastLoginTimestamp":           "2024-01-01T00:00:00Z",
		"passwordLastChangedTimestamp": "2024-01-01T00:00:00Z",
		"createdTimestamp":             "2024-01-01T00:00:00Z",
		"revision":                     float64(1),
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
		"infrastructureServerGroupMaxCount":                float64(10),
		"infrastructureDriveMaxCount":                      float64(10),
		"infrastructureFileShareMaxCount":                  float64(10),
		"infrastructureBucketMaxCount":                     float64(10),
		"infrastructureVmInstanceGroupMaxCount":            float64(10),
		"infrastructureContainerInstanceGroupMaxCount":     float64(10),
		"serverGroupInstancesMaxCount":                     float64(10),
		"serverGroupInstancesMinCount":                     float64(1),
		"vmInstanceGroupVmInstancesMaxCount":               float64(10),
		"containerInstanceGroupContainerInstancesMaxCount": float64(10),
		"vmInstanceMaxDiskSizeMbytes":                      float64(1048576),
		"containerInstanceMaxDiskSizeMbytes":               float64(1048576),
		"driveMaxSizeMbytes":                               float64(1048576),
		"driveMinSizeMbytes":                               float64(1024),
		"fileShareMinSizeGb":                               float64(1),
		"fileShareMaxSizeGb":                               float64(1024),
		"bucketMinSizeGb":                                  float64(1),
		"bucketMaxSizeGb":                                  float64(1024),
		"showOperatingSystemImagesTab":                     true,
		"showTemplateAssetsView":                           true,
		"userResourceServerTypeNameToMaxCount":             map[string]any{},
		"userSshKeysCountMax":                              float64(10),
		"showLegacyPages":                                  false,
		"showEliChatBot":                                   false,
		"enableCustomRaidConfiguration":                    true,
		"enableInfrastructureVmInstance":                   true,
		"enableInfrastructureContainerInstance":            true,
		"enableInfrastructureExtensions":                   true,
		"allowedInfrastructureExtensions":                  []any{},
		"allowedServerTypes":                               []any{},
		"allowedSites":                                     []any{},
		"allowedLogicalNetworkProfiles":                    []any{},
		"allowedPreCreatedLogicalNetworks":                 []any{},
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
		"/api/v2/users/1":                      testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/set-password": testutils.JSONHandler(200, makeUser(1)),
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
		"/api/v2/users/1":                      testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/set-password": testutils.ErrorHandler(500, "internal server error"),
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
		"/api/v2/users/1":                        testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/change-account": testutils.JSONHandler(200, makeUser(1)),
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

// ---------------------------------------------------------------------------
// Request recorder
// ---------------------------------------------------------------------------

type recordedRequest struct {
	Method string
	Path   string
	Query  url.Values
	Header http.Header
	Body   string
}

type requestRecorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

func (r *requestRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, recordedRequest{
			Method: req.Method,
			Path:   req.URL.Path,
			Query:  req.URL.Query(),
			Header: req.Header.Clone(),
			Body:   string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *requestRecorder) last() recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.reqs) == 0 {
		return recordedRequest{}
	}
	return r.reqs[len(r.reqs)-1]
}

func makeUserInfo(id int) map[string]any {
	return map[string]any{
		"id": id, "revision": float64(1), "displayName": "Alice",
		"email": "alice@example.com", "emailStatus": "active",
		"language": "en", "accessLevel": "admin", "isArchived": false,
		"lastLoginTimestamp": "2024-01-01T00:00:00Z",
		"createdTimestamp":   "2024-01-01T00:00:00Z",
	}
}

func makeSshKey(id int) map[string]any {
	return map[string]any{
		"id": id, "userId": 1, "sshKey": "ssh-rsa AAAA...",
		"status": "active", "createdTimestamp": "2024-01-01T00:00:00Z",
	}
}

// ---------------------------------------------------------------------------
// GetConfiguration
// ---------------------------------------------------------------------------

func TestGetConfiguration_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/config": testutils.JSONHandler(200, userCoreFields()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetConfiguration(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGetConfiguration_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetConfiguration(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestGetConfiguration_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetConfiguration(ctx, "bad"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

// ---------------------------------------------------------------------------
// UpdateMeta
// ---------------------------------------------------------------------------

func TestUpdateMeta_HappyPath(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/meta": rec.wrap(testutils.JSONHandler(200, map[string]any{"guiSettings": map[string]any{}})),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := UpdateMeta(ctx, "1", []byte(`{"guiSettings":{}}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	req := rec.last()
	if req.Method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", req.Method)
	}
	if !strings.Contains(req.Body, "guiSettings") {
		t.Errorf("expected the metadata in the body, got %q", req.Body)
	}
}

func TestUpdateMeta_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/meta": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := UpdateMeta(ctx, "1", []byte(`{"guiSettings":{}}`)); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// ---------------------------------------------------------------------------
// SSH key / suspend reasons
// ---------------------------------------------------------------------------

func TestGetSSHKey_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/ssh-keys/5": testutils.JSONHandler(200, makeSshKey(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetSSHKey(ctx, "1", "5"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGetSSHKey_InvalidKeyId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetSSHKey(ctx, "1", "bad"); err == nil {
		t.Fatal("expected error on invalid key ID, got nil")
	}
}

func TestGetSuspendReasons_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/suspend-reasons": testutils.JSONHandler(200, map[string]any{
			"data": []any{makeSuspendReason()},
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetSuspendReasons(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGetSuspendReasons_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/suspend-reasons": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetSuspendReasons(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// ---------------------------------------------------------------------------
// Delegates
// ---------------------------------------------------------------------------

func TestAddDelegate_HappyPath(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/actions/add-delegate/2": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AddDelegate(ctx, "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); req.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", req.Method)
	}
}

func TestAddDelegate_InvalidDelegateId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AddDelegate(ctx, "1", "bad"); err == nil {
		t.Fatal("expected error on invalid delegate ID, got nil")
	}
}

func TestRemoveDelegate_HappyPath(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/actions/remove-delegate/2": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := RemoveDelegate(ctx, "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); req.Path != "/api/v2/users/1/actions/remove-delegate/2" {
		t.Errorf("unexpected path %s", req.Path)
	}
}

func TestGetParentDelegates_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/parent-delegates": testutils.JSONHandler(200, map[string]any{
			"data": []any{makeUserInfo(2)},
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetParentDelegates(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGetChildDelegates_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/child-delegates": testutils.JSONHandler(200, map[string]any{
			"data": []any{makeUserInfo(3)},
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetChildDelegates(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestGetChildDelegates_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/child-delegates": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetChildDelegates(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// ---------------------------------------------------------------------------
// Notification actions
// ---------------------------------------------------------------------------

// TestResendEmailVerification_SendsBody guards the SDK gotcha where an unset
// optional body is serialized as a literal `null`, which the API rejects.
func TestResendEmailVerification_SendsBody(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/actions/resend-email-verification": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResendEmailVerification(ctx, "1", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	req := rec.last()
	if strings.TrimSpace(req.Body) == "null" || strings.TrimSpace(req.Body) == "" {
		t.Errorf("expected a JSON object body, got %q", req.Body)
	}
}

func TestResendEmailVerification_RedirectUrl(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/actions/resend-email-verification": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResendEmailVerification(ctx, "1", "https://example.com/welcome"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); !strings.Contains(req.Body, "https://example.com/welcome") {
		t.Errorf("expected the redirect URL in the body, got %q", req.Body)
	}
}

func TestResendInvitation_HappyPath(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/actions/resend-user-invitation": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResendInvitation(ctx, "1", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); strings.TrimSpace(req.Body) == "null" {
		t.Errorf("expected a JSON object body, got %q", req.Body)
	}
}

func TestSendPasswordReset_HappyPath(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/actions/send-password-reset": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SendPasswordReset(ctx, "1", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); strings.TrimSpace(req.Body) == "null" {
		t.Errorf("expected a JSON object body, got %q", req.Body)
	}
}

func TestSendPasswordReset_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1/actions/send-password-reset": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SendPasswordReset(ctx, "1", ""); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete_SendsIfMatch(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1":                testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/delete": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Delete(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	req := rec.last()
	if req.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", req.Method)
	}
	if req.Header.Get("If-Match") != "1" {
		t.Errorf("expected If-Match 1, got %q", req.Header.Get("If-Match"))
	}
}

func TestDelete_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1":                testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/users/1/actions/delete": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Delete(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// ---------------------------------------------------------------------------
// Self-service operations
// ---------------------------------------------------------------------------

func TestGetPermissions_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/permissions": testutils.JSONHandler(200, map[string]any{
			"permissions": []any{"servers_read", "users_read"},
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetPermissions(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestGetPermissions_UnknownPermission covers the enum desync: the SDK models
// the permission keys as a strict enum, so a permission it does not know about
// must not break the command.
func TestGetPermissions_UnknownPermission(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/permissions": testutils.JSONHandler(200, map[string]any{
			"permissions": []any{"servers_read", "brand_new_permission_the_sdk_ignores"},
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	out := testutils.CaptureStdout(t, func() {
		if err := GetPermissions(ctx); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "brand_new_permission_the_sdk_ignores") {
		t.Errorf("expected the unknown permission in the output, got: %s", out)
	}
}

func TestGetPermissions_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/permissions": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := GetPermissions(ctx); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestChangePassword_SendsIfMatch(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1":                      testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/user/actions/change-password": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	ctx = api.SetUserId(ctx, "1")

	if err := ChangePassword(ctx, []byte(`{"newPassword":"new","oldPassword":"old"}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	req := rec.last()
	if req.Header.Get("If-Match") != "1" {
		t.Errorf("expected If-Match 1, got %q", req.Header.Get("If-Match"))
	}
	if !strings.Contains(req.Body, "newPassword") {
		t.Errorf("expected the new password in the body, got %q", req.Body)
	}
}

func TestInitiatePasswordReset_SendsEmail(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/actions/initiate-password-reset": rec.wrap(testutils.RawHandler(204, "")),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InitiatePasswordReset(ctx, "alice@example.com", "https://example.com/login"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	req := rec.last()
	if !strings.Contains(req.Body, "alice@example.com") {
		t.Errorf("expected the e-mail address in the body, got %q", req.Body)
	}
	if !strings.Contains(req.Body, "https://example.com/login") {
		t.Errorf("expected the redirect URL in the body, got %q", req.Body)
	}
}

func TestInitiateEmailChange_HappyPath(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/actions/initiate-email-change": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InitiateEmailChange(ctx, "new@example.com", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); !strings.Contains(req.Body, "new@example.com") {
		t.Errorf("expected the new address in the body, got %q", req.Body)
	}
}

func TestRegenerateJwtSalt_SendsIfMatch(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/users/1":                          testutils.JSONHandler(200, makeUser(1)),
		"/api/v2/user/actions/regenerate-jwt-salt": rec.wrap(testutils.JSONHandler(200, makeUser(1))),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	ctx = api.SetUserId(ctx, "1")

	if err := RegenerateJwtSalt(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); req.Header.Get("If-Match") != "1" {
		t.Errorf("expected If-Match 1, got %q", req.Header.Get("If-Match"))
	}
}

func TestVerifyEmail_SendsToken(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/actions/verify-email": rec.wrap(testutils.RawHandler(204, "")),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VerifyEmail(ctx, "token-123"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	req := rec.last()
	if req.Method != http.MethodGet {
		t.Errorf("expected GET, got %s", req.Method)
	}
	if req.Query.Get("token") != "token-123" {
		t.Errorf("expected token in the query string, got %q", req.Query.Encode())
	}
}

func TestResetPassword_SendsToken(t *testing.T) {
	rec := &requestRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/actions/reset-password": rec.wrap(testutils.RawHandler(204, "")),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResetPassword(ctx, "token-456"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if req := rec.last(); req.Query.Get("token") != "token-456" {
		t.Errorf("expected token in the query string, got %q", req.Query.Encode())
	}
}

func TestResetPassword_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/user/actions/reset-password": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResetPassword(ctx, "token-456"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}
