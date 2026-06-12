package server_type

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

func TestServerTypeList(t *testing.T) {
	const listPage1 = `{
		"data": [
			{"id": 1, "name": "type-1", "label": "type-1",
			 "processorCount": 2, "processorCoreMhz": 2400, "processorCoreCount": 8,
			 "processorNames": ["Intel"],
			 "ramGbytes": 32, "networkInterfaceCount": 2, "networkTotalCapacityMbps": 10000,
			 "networkInterfaceSpeeds": [], "diskCount": 2, "serverClass": "M",
			 "links": []}
		],
		"meta": {"currentPage": 1, "totalPages": 1, "itemsPerPage": 100}
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-types": testutils.RawHandler(http.StatusOK, listPage1),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerTypeList(ctx); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-types": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerTypeList(ctx); err == nil {
			t.Error("expected error for HTTP 500, got nil")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-types": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerTypeList(ctx); err != nil {
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
				s += fmt.Sprintf(`{"id":%d,"name":"type-%d","label":"type-%d","processorCount":2,"processorCoreMhz":2400,"processorCoreCount":8,"processorNames":["Intel"],"ramGbytes":32,"networkInterfaceCount":2,"networkTotalCapacityMbps":10000,"networkInterfaceSpeeds":[],"diskCount":2,"serverClass":"M","links":[]}`, id, id, id)
			}
			s += `]`
			return s
		}

		var call atomic.Int32
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-types": func(w http.ResponseWriter, r *http.Request) {
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
		if err := ServerTypeList(ctx); err != nil {
			t.Errorf("expected nil error for pagination, got: %v", err)
		}
		if got := int(call.Load()); got < 3 {
			t.Errorf("expected at least 3 page fetches, got %d", got)
		}
	})
}

func TestServerTypeGet(t *testing.T) {
	const typeSingle = `{
		"id": 1, "name": "type-1", "label": "type-1",
		"processorCount": 2, "processorCoreMhz": 2400, "processorCoreCount": 8,
		"processorNames": ["Intel"],
		"ramGbytes": 32, "networkInterfaceCount": 2, "networkTotalCapacityMbps": 10000,
		"networkInterfaceSpeeds": [], "diskCount": 2, "serverClass": "M",
		"links": []
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-types/1": testutils.RawHandler(http.StatusOK, typeSingle),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerTypeGet(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-types/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerTypeGet(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})
}
