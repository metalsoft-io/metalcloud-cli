package server_instance

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// infraListResponse returns a minimal infrastructure list JSON with ID=123 and the given label.
func infraListResponse(label string) string {
	return fmt.Sprintf(`{"data":[{"id":123,"revision":1,"label":%q,"serviceStatus":"active","datacenterName":"dc1","siteId":1,"designIsLocked":0,"config":{},"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","links":[]}],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`, label)
}

func TestServerInstanceGroupList(t *testing.T) {
	const sigListPage1 = `{
		"data": [
			{"id": 10, "revision": 1, "label": "sig-1", "infrastructureId": 123,
			 "instanceCount": 1, "defaultServerTypeId": 5,
			 "ipAllocateAuto": 1, "ipv4SubnetCreateAuto": 1,
			 "processorCount": 2, "processorCoreCount": 8, "processorCoreMhz": 2400,
			 "diskCount": 2, "diskSizeMbytes": 102400, "diskTypes": [],
			 "virtualInterfacesEnabled": 0,
			 "serviceStatus": "active", "isVmGroup": 0, "isEndpointInstanceGroup": 0,
			 "meta": {},
			 "createdTimestamp": "2024-01-01T00:00:00Z",
			 "updatedTimestamp": "2024-01-01T00:00:00Z",
			 "links": []}
		],
		"meta": {"currentPage": 1, "totalPages": 1, "itemsPerPage": 100}
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/infrastructures": func(w http.ResponseWriter, r *http.Request) {
				testutils.RawHandler(http.StatusOK, infraListResponse("infra-1"))(w, r)
			},
			"/api/v2/infrastructures/123/server-instance-groups": testutils.RawHandler(http.StatusOK, sigListPage1),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerInstanceGroupList(ctx, "123"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("InfraNotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/infrastructures": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerInstanceGroupList(ctx, "999"); err == nil {
			t.Error("expected error when infrastructure not found, got nil")
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/infrastructures": func(w http.ResponseWriter, r *http.Request) {
				testutils.RawHandler(http.StatusOK, infraListResponse("infra-1"))(w, r)
			},
			"/api/v2/infrastructures/123/server-instance-groups": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerInstanceGroupList(ctx, "123"); err == nil {
			t.Error("expected error for HTTP 500, got nil")
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
				s += fmt.Sprintf(`{"id":%d,"revision":1,"label":"sig-%d","infrastructureId":123,"instanceCount":1,"defaultServerTypeId":5,"ipAllocateAuto":1,"ipv4SubnetCreateAuto":1,"processorCount":2,"processorCoreCount":8,"processorCoreMhz":2400,"diskCount":2,"diskSizeMbytes":102400,"diskTypes":[],"virtualInterfacesEnabled":0,"serviceStatus":"active","isVmGroup":0,"isEndpointInstanceGroup":0,"meta":{},"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","links":[]}`, id, id)
			}
			s += `]`
			return s
		}

		var call atomic.Int32
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/infrastructures": func(w http.ResponseWriter, r *http.Request) {
				testutils.RawHandler(http.StatusOK, infraListResponse("infra-1"))(w, r)
			},
			"/api/v2/infrastructures/123/server-instance-groups": func(w http.ResponseWriter, r *http.Request) {
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
		if err := ServerInstanceGroupList(ctx, "123"); err != nil {
			t.Errorf("expected nil error for pagination, got: %v", err)
		}
		if got := int(call.Load()); got < 3 {
			t.Errorf("expected at least 3 page fetches for SIG list, got %d", got)
		}
	})
}

func TestServerInstanceGroupGet(t *testing.T) {
	const sigSingle = `{
		"id": 10, "revision": 1, "label": "sig-1", "infrastructureId": 123,
		"instanceCount": 1, "defaultServerTypeId": 5,
		"ipAllocateAuto": 1, "ipv4SubnetCreateAuto": 1,
		"processorCount": 2, "processorCoreCount": 8, "processorCoreMhz": 2400,
		"diskCount": 2, "diskSizeMbytes": 102400, "diskTypes": [],
		"virtualInterfacesEnabled": 0,
		"serviceStatus": "active", "isVmGroup": 0, "isEndpointInstanceGroup": 0,
		"meta": {},
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"links": []
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-instance-groups/10": testutils.RawHandler(http.StatusOK, sigSingle),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerInstanceGroupGet(ctx, "10"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/server-instance-groups/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := ServerInstanceGroupGet(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})
}
