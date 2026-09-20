package account

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	viper.Set("format", "json")
	m.Run()
}

func setupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

// accountJSON has all required fields: name, id, revision, limits, config.
const accountJSON = `{"id":1,"name":"acme","revision":1,"limits":{},"config":{"revision":1,"name":"acme"}}`

func accountListHandler(statusCode int, count, currentPage, totalPages int) http.HandlerFunc {
	items := make([]string, count)
	for i := range items {
		items[i] = accountJSON
	}
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func accountPage(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{
			"id": i + 1, "name": fmt.Sprintf("acme-%d", i+1), "revision": 1,
			"limits": map[string]any{},
			"config": map[string]any{"revision": 1, "name": fmt.Sprintf("acme-%d", i+1)},
		}
	}
	return items
}

func TestAccountList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": accountListHandler(http.StatusOK, 2, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestAccountList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": accountListHandler(http.StatusOK, 0, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestAccountList_Pagination(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/accounts", []any{accountPage(100), accountPage(100), accountPage(5)})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestAccountGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1": testutils.RawHandler(http.StatusOK, accountJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestAccountGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestAccountCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, accountJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"name":"acme","code":"acme-code"}`)
	if err := AccountCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountCreate_BadRequest(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{}`)
	if err := AccountCreate(ctx, config); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

func TestAccountArchive_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1": testutils.RawHandler(http.StatusOK, accountJSON),
		"/api/v2/accounts/1/actions/archive": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, accountJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountArchive(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountArchive_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountArchive(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

// ---------------------------------------------------------------------------
// Unarchive / config / quota breakdown
// ---------------------------------------------------------------------------

// accountQuotaLimits returns a QuotaProfileLimits payload with every required
// property populated.
func accountQuotaLimits(serverGroups float64) map[string]any {
	return map[string]any{
		"infrastructureServerGroupMaxCount":                serverGroups,
		"infrastructureDriveMaxCount":                      10,
		"infrastructureFileShareMaxCount":                  10,
		"infrastructureBucketMaxCount":                     10,
		"infrastructureVmInstanceGroupMaxCount":            10,
		"infrastructureContainerInstanceGroupMaxCount":     10,
		"serverGroupInstancesMaxCount":                     10,
		"serverGroupInstancesMinCount":                     1,
		"vmInstanceGroupVmInstancesMaxCount":               10,
		"containerInstanceGroupContainerInstancesMaxCount": 10,
		"vmInstanceMaxDiskSizeMbytes":                      1024,
		"containerInstanceMaxDiskSizeMbytes":               1024,
		"driveMaxSizeMbytes":                               1024,
		"driveMinSizeMbytes":                               1,
		"fileShareMinSizeGb":                               1,
		"fileShareMaxSizeGb":                               100,
		"bucketMinSizeGb":                                  1,
		"bucketMaxSizeGb":                                  100,
		"showOperatingSystemImagesTab":                     true,
		"showTemplateAssetsView":                           true,
		"userResourceServerTypeNameToMaxCount":             map[string]any{"M.8.8.2": 4},
		"userSshKeysCountMax":                              5,
		"showLegacyPages":                                  false,
		"showEliChatBot":                                   false,
		"enableCustomRaidConfiguration":                    false,
		"enableInfrastructureVmInstance":                   true,
		"enableInfrastructureContainerInstance":            true,
		"enableInfrastructureExtensions":                   true,
		"allowedInfrastructureExtensions":                  []any{},
		"allowedServerTypes":                               []any{"M.8.8.2"},
		"allowedSites":                                     []any{},
		"allowedLogicalNetworkProfiles":                    []any{},
		"allowedPreCreatedLogicalNetworks":                 []any{},
	}
}

func TestAccountUnarchive_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1":                   testutils.RawHandler(http.StatusOK, accountJSON),
		"/api/v2/accounts/1/actions/unarchive": testutils.RawHandler(http.StatusOK, accountJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountUnarchive(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestAccountUnarchive_SendsIfMatch verifies the optimistic concurrency header
// is taken from the account fetched first.
func TestAccountUnarchive_SendsIfMatch(t *testing.T) {
	gotIfMatch := ""
	gotMethod := ""
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1": testutils.RawHandler(http.StatusOK, accountJSON),
		"/api/v2/accounts/1/actions/unarchive": func(w http.ResponseWriter, r *http.Request) {
			gotIfMatch = r.Header.Get("If-Match")
			gotMethod = r.Method
			testutils.RawHandler(http.StatusOK, accountJSON)(w, r)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountUnarchive(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if gotIfMatch != "1" {
		t.Errorf("expected If-Match 1, got %q", gotIfMatch)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
}

func TestAccountUnarchive_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1":                   testutils.RawHandler(http.StatusOK, accountJSON),
		"/api/v2/accounts/1/actions/unarchive": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountUnarchive(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestAccountUnarchive_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountUnarchive(ctx, "bad"); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

func TestAccountGetConfig_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1/config": testutils.RawHandler(http.StatusOK, `{"revision":1,"name":"acme","code":"ACME"}`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGetConfig(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestAccountGetConfig_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1/config": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGetConfig(ctx, "1"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestAccountQuotaBreakdown_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1/quota-limits-breakdown": testutils.JSONHandler(http.StatusOK, map[string]any{
			"effective": accountQuotaLimits(5),
			"account":   accountQuotaLimits(7),
		}),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountQuotaBreakdown(ctx, "1", false); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestAccountQuotaBreakdown_IncludeUsage checks that the flag reaches the API
// as the documented query parameter.
func TestAccountQuotaBreakdown_IncludeUsage(t *testing.T) {
	gotIncludeUsage := ""
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1/quota-limits-breakdown": func(w http.ResponseWriter, r *http.Request) {
			gotIncludeUsage = r.URL.Query().Get("includeUsage")
			testutils.JSONHandler(http.StatusOK, map[string]any{
				"effective":    accountQuotaLimits(5),
				"accountUsage": map[string]any{"infrastructureDriveMaxCount": 2},
			})(w, r)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)

	if err := AccountQuotaBreakdown(ctx, "1", false); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if gotIncludeUsage != "" {
		t.Errorf("expected no includeUsage parameter, got %q", gotIncludeUsage)
	}

	if err := AccountQuotaBreakdown(ctx, "1", true); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if gotIncludeUsage != "true" {
		t.Errorf("expected includeUsage=true, got %q", gotIncludeUsage)
	}
}

// TestAccountQuotaBreakdown_Flattened verifies the text formats get one row per
// limit instead of an unreadable nested object.
func TestAccountQuotaBreakdown_Flattened(t *testing.T) {
	viper.Set("format", "csv")
	defer viper.Set("format", "json")

	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1/quota-limits-breakdown": testutils.JSONHandler(http.StatusOK, map[string]any{
			"effective": accountQuotaLimits(5),
			"account":   accountQuotaLimits(7),
		}),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	out := testutils.CaptureStdout(t, func() {
		if err := AccountQuotaBreakdown(ctx, "1", false); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "infrastructureServerGroupMaxCount") {
		t.Errorf("expected one row per limit, got: %s", out)
	}
	// The effective value wins over the less restrictive account value.
	if !strings.Contains(out, "5,7") {
		t.Errorf("expected the effective and account columns, got: %s", out)
	}
}

// TestAccountQuotaBreakdown_NullScope covers the scopes the API returns as null
// when no quota profile applies.
func TestAccountQuotaBreakdown_NullScope(t *testing.T) {
	viper.Set("format", "csv")
	defer viper.Set("format", "json")

	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1/quota-limits-breakdown": testutils.RawHandler(http.StatusOK,
			`{"effective":`+quotaLimitsJSON()+`,"account":null,"parentAccount":null}`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	out := testutils.CaptureStdout(t, func() {
		if err := AccountQuotaBreakdown(ctx, "1", false); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "-") {
		t.Errorf("expected unset scopes to render as '-', got: %s", out)
	}
}

func TestAccountQuotaBreakdown_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1/quota-limits-breakdown": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountQuotaBreakdown(ctx, "1", false); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestAccountQuotaBreakdown_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountQuotaBreakdown(ctx, "bad", false); err == nil {
		t.Fatal("expected error on invalid ID, got nil")
	}
}

func quotaLimitsJSON() string {
	encoded, err := json.Marshal(accountQuotaLimits(5))
	if err != nil {
		panic(err)
	}
	return string(encoded)
}
