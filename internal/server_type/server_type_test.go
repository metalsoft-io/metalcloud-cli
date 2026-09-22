package server_type

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

const serverTypeSingleJSON = `{
	"id": 1, "name": "type-1", "label": "type-1",
	"processorCount": 2, "processorCoreMhz": 2400, "processorCoreCount": 8,
	"processorNames": ["Intel"],
	"ramGbytes": 32, "networkInterfaceCount": 2, "networkTotalCapacityMbps": 10000,
	"networkInterfaceSpeeds": [], "diskCount": 2, "serverClass": "M",
	"links": []
}`

// stypeRecorder records the requests the mock server receives.
type stypeRecorder struct {
	Method string
	Path   string
	Body   string
}

func (r *stypeRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.Method = req.Method
		r.Path = req.URL.Path
		r.Body = string(body)
		next(w, req)
	}
}

func TestServerTypeCreate(t *testing.T) {
	rec := &stypeRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-types": rec.wrap(testutils.RawHandler(http.StatusOK, serverTypeSingleJSON)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	config := []byte(`{
		"name":"type-1","label":"type-1","ramGbytes":32,"processorCount":2,
		"processorCoreMhz":2400,"processorCoreCount":8,"processorNames":["Intel"],
		"networkInterfaceSpeeds":[10000],"diskCount":2,"serverClass":"M"
	}`)

	if err := ServerTypeCreate(ctx, config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if rec.Method != http.MethodPost {
		t.Errorf("create should be a POST, got %s", rec.Method)
	}
	if !strings.Contains(rec.Body, `"serverClass":"M"`) {
		t.Errorf("expected the create body to be sent, got %q", rec.Body)
	}
}

func TestServerTypeCreateInvalidConfig(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := ServerTypeCreate(ctx, []byte("not json")); err == nil {
		t.Fatal("expected an error for an invalid configuration")
	}
}

func TestServerTypeUpdate(t *testing.T) {
	rec := &stypeRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-types/1": rec.wrap(testutils.RawHandler(http.StatusOK, serverTypeSingleJSON)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := ServerTypeUpdate(ctx, "1", []byte(`{"label":"type-1","description":"updated"}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !strings.Contains(rec.Body, `"description":"updated"`) {
		t.Errorf("expected the update body to be sent, got %q", rec.Body)
	}
}

func TestServerTypeDelete(t *testing.T) {
	rec := &stypeRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-types/1": rec.wrap(testutils.NoContentHandler()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := ServerTypeDelete(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if rec.Method != http.MethodDelete {
		t.Errorf("delete should be a DELETE, got %s", rec.Method)
	}
}

func TestServerTypeDeleteInvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := ServerTypeDelete(ctx, "abc"); err == nil {
		t.Fatal("expected an error for an invalid server type ID")
	}
}

func TestServerTypeCleanUnused(t *testing.T) {
	rec := &stypeRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-types/actions/clean-unused": rec.wrap(testutils.NoContentHandler()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := ServerTypeCleanUnused(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if rec.Method != http.MethodPost {
		t.Errorf("clean-unused should be a POST, got %s", rec.Method)
	}
}

func TestServerTypeStatistics(t *testing.T) {
	const statisticsJSON = `{
		"serverTypeIdToServerCount": {"1": 3, "2": 0},
		"serverTypeIdToServerInformation": {},
		"utilizationReport": {
			"groupByServerRamGb": {},
			"groupByServerTypeName": {},
			"groupByServerProductName": {},
			"groupByUserId": {}
		}
	}`

	rec := &stypeRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-types/statistics": rec.wrap(testutils.RawHandler(http.StatusOK, statisticsJSON)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerTypeStatistics(ctx, 1, []string{"1", "2"}, 5, 10, 7); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "serverTypeIdToServerCount") {
		t.Errorf("expected the counts in the output, got: %s", out)
	}

	// The body is always sent; an omitted body would serialize as a literal
	// null, which the API rejects with 400.
	for _, want := range []string{`"siteId":1`, `"serverTypeIds":[1,2]`, `"userIdOwner":5`,
		`"maximumResultsPerServerType":10`, `"instanceArrayId":7`} {
		if !strings.Contains(rec.Body, want) {
			t.Errorf("expected %s in the request body, got %q", want, rec.Body)
		}
	}
}

func TestServerTypeStatisticsInvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := ServerTypeStatistics(ctx, 1, []string{"abc"}, 0, 0, 0); err == nil {
		t.Fatal("expected an error for an invalid server type ID")
	}
}

func TestServerTypeConfigExamples(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")

	out := testutils.CaptureStdout(t, func() {
		if err := ServerTypeConfigExample(ctx); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})
	if !strings.Contains(out, "serverClass") {
		t.Errorf("expected a create example, got: %s", out)
	}

	out = testutils.CaptureStdout(t, func() {
		if err := ServerTypeUpdateConfigExample(ctx); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})
	if !strings.Contains(out, "label") {
		t.Errorf("expected an update example, got: %s", out)
	}
}
