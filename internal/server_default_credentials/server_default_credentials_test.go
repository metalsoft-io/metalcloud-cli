package server_default_credentials

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

const credentialsItem = `{
	"id": 1, "siteId": 10,
	"serverSerialNumber": "SN123", "serverMacAddress": "AA:BB:CC:DD:EE:FF",
	"defaultUsername": "admin",
	"defaultRackName": "rack-1", "defaultRackPositionLowerUnit": "1", "defaultRackPositionUpperUnit": "2",
	"defaultInventoryId": "inv-1", "defaultUuid": "uuid-1",
	"links": []
}`

func TestServerDefaultCredentialsList(t *testing.T) {
	listPage1 := fmt.Sprintf(`{"data":[%s],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`, credentialsItem)

	t.Run("FetchAll_HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials": testutils.RawHandler(http.StatusOK, listPage1),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		// page=0, limit=0 → fetch-all path
		if err := ServerDefaultCredentialsList(ctx, 0, 0); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("FetchAll_HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsList(ctx, 0, 0); err == nil {
			t.Error("expected error for HTTP 500, got nil")
		}
	})

	t.Run("FetchAll_EmptyList", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsList(ctx, 0, 0); err != nil {
			t.Errorf("expected nil error for empty list, got: %v", err)
		}
	})

	t.Run("FetchAll_Pagination3Pages", func(t *testing.T) {
		makeItems := func(start, count int) string {
			s := `[`
			for i := 0; i < count; i++ {
				if i > 0 {
					s += ","
				}
				id := start + i
				s += fmt.Sprintf(`{"id":%d,"siteId":10,"serverSerialNumber":"SN%d","serverMacAddress":"AA:BB:CC:DD:EE:FF","defaultUsername":"admin","defaultRackName":"rack-1","defaultRackPositionLowerUnit":"1","defaultRackPositionUpperUnit":"2","defaultInventoryId":"inv-1","defaultUuid":"uuid-1","links":[]}`, id, id)
			}
			s += `]`
			return s
		}

		var call atomic.Int32
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials": func(w http.ResponseWriter, r *http.Request) {
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
		if err := ServerDefaultCredentialsList(ctx, 0, 0); err != nil {
			t.Errorf("expected nil error for pagination, got: %v", err)
		}
		if got := int(call.Load()); got < 3 {
			t.Errorf("expected at least 3 page fetches, got %d", got)
		}
	})

	t.Run("SinglePage_page1_limit5", func(t *testing.T) {
		// page>0 || limit>0 → manual pagination branch
		items := fmt.Sprintf(`[%s]`, credentialsItem)
		body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":5}}`, items)
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials": testutils.RawHandler(http.StatusOK, body),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsList(ctx, 1, 5); err != nil {
			t.Errorf("expected nil error for page=1 limit=5, got: %v", err)
		}
	})

	t.Run("SinglePage_HttpError", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsList(ctx, 1, 5); err == nil {
			t.Error("expected error for HTTP 500 in single-page path, got nil")
		}
	})
}

func TestServerDefaultCredentialsGet(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials/1": testutils.RawHandler(http.StatusOK, credentialsItem),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsGet(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsGet(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})
}

func TestServerDefaultCredentialsCreate(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials": func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					testutils.RawHandler(http.StatusOK, credentialsItem)(w, r)
					return
				}
				testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`)(w, r)
			},
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsCreate(ctx, 10, "SN123", "AA:BB:CC:DD:EE:FF", "admin", "pass", "rack-1", "1", "2", "inv-1", "uuid-1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}

func TestServerDefaultCredentialsDelete(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/default-credentials/1": func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				testutils.RawHandler(http.StatusOK, credentialsItem)(w, r)
			},
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerDefaultCredentialsDelete(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}
