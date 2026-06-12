package job

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestJobGroupList(t *testing.T) {
	const groupListPage1 = `{
		"data": [
			{"id": 1, "type": "deploy", "description": "group 1",
			 "createdTimestamp": "2024-01-01T00:00:00Z", "links": []}
		],
		"meta": {"currentPage": 1, "totalPages": 1, "itemsPerPage": 100}
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/job-groups": testutils.RawHandler(http.StatusOK, groupListPage1),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGroupList(ctx, GroupListFlags{}); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/job-groups": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGroupList(ctx, GroupListFlags{}); err == nil {
			t.Error("expected error for HTTP 500, got nil")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/job-groups": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGroupList(ctx, GroupListFlags{}); err != nil {
			t.Errorf("expected nil error for empty list, got: %v", err)
		}
	})

	t.Run("Pagination3Pages", func(t *testing.T) {
		// 100 items page 1, 100 items page 2, 5 items page 3 = 205 total
		makeItems := func(start, count int) string {
			s := `[`
			for i := 0; i < count; i++ {
				if i > 0 {
					s += ","
				}
				id := start + i
				s += fmt.Sprintf(`{"id":%d,"type":"deploy","description":"group %d","createdTimestamp":"2024-01-01T00:00:00Z","links":[]}`, id, id)
			}
			s += `]`
			return s
		}

		var call atomic.Int32
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/job-groups": func(w http.ResponseWriter, r *http.Request) {
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
		if err := JobGroupList(ctx, GroupListFlags{}); err != nil {
			t.Errorf("expected nil error for pagination, got: %v", err)
		}
		if got := int(call.Load()); got < 3 {
			t.Errorf("expected at least 3 page fetches, got %d", got)
		}
	})
}

func TestJobGroupGet(t *testing.T) {
	const groupSingle = `{
		"id": 1, "type": "deploy", "description": "group 1",
		"createdTimestamp": "2024-01-01T00:00:00Z", "links": []
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/job-groups/1": testutils.RawHandler(http.StatusOK, groupSingle),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGroupGet(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/job-groups/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGroupGet(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGroupGet(ctx, "not-a-number"); err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})
}
