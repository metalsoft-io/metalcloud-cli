package server_cleanup_policy

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

const cleanupPolicyItem = `{
	"id": 1, "label": "policy-1",
	"cleanupDrivesForOobEnabledServer": 1, "recreateRaid": 0,
	"clearTpm": 0,
	"resetRaidControllers": 0, "disableEmbeddedNics": 0,
	"raidOneDrive": "raid1", "raidTwoDrives": "raid1",
	"raidEvenNumberMoreThanTwoDrives": "raid5", "raidOddNumberMoreThanOneDrive": "raid5",
	"skipRaidActions": [],
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"links": []
}`

func TestCleanupPolicyList(t *testing.T) {
	listPage1 := fmt.Sprintf(`{"data":[%s],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`, cleanupPolicyItem)

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies": testutils.RawHandler(http.StatusOK, listPage1),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyList(ctx); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyList(ctx); err == nil {
			t.Error("expected error for HTTP 500, got nil")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyList(ctx); err != nil {
			t.Errorf("expected nil error for empty list, got: %v", err)
		}
	})

	t.Run("Pagination3Pages", func(t *testing.T) {
		makeItems := func(start, count int) string {
			s := `[`
			for i := 0; i < count; i++ {
				if i > 0 {
					s += ","
				}
				id := start + i
				s += fmt.Sprintf(`{"id":%d,"label":"policy-%d","cleanupDrivesForOobEnabledServer":0,"clearTpm":0,"recreateRaid":0,"resetRaidControllers":0,"disableEmbeddedNics":0,"raidOneDrive":"","raidTwoDrives":"","raidEvenNumberMoreThanTwoDrives":"","raidOddNumberMoreThanOneDrive":"","skipRaidActions":[],"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","links":[]}`, id, id)
			}
			s += `]`
			return s
		}

		var call atomic.Int32
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies": func(w http.ResponseWriter, r *http.Request) {
				n := int(call.Add(1))
				var page, total int
				var items string
				switch n {
				case 1:
					page, total = 1, 3
					items = makeItems(1, 100)
				case 2:
					page, total = 2, 3
					items = makeItems(101, 100)
				default:
					page, total = 3, 3
					items = makeItems(201, 5)
				}
				body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`, items, page, total)
				testutils.RawHandler(http.StatusOK, body)(w, r)
			},
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyList(ctx); err != nil {
			t.Errorf("expected nil error for pagination, got: %v", err)
		}
		if got := int(call.Load()); got < 3 {
			t.Errorf("expected at least 3 page fetches, got %d", got)
		}
	})
}

func TestCleanupPolicyGet(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies/1": testutils.RawHandler(http.StatusOK, cleanupPolicyItem),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyGet(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyGet(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})
}

func TestCleanupPolicyCreate(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies": func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					testutils.RawHandler(http.StatusOK, cleanupPolicyItem)(w, r)
					return
				}
				testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`)(w, r)
			},
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyCreate(ctx, "policy-1", true, false, false, "raid1", "raid1", "raid5", "raid5", ""); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}

func TestCleanupPolicyDelete(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/cleanup-policies/1": func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				testutils.RawHandler(http.StatusOK, cleanupPolicyItem)(w, r)
			},
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := CleanupPolicyDelete(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}
