package subnet

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// subnetJSON is a minimal valid Subnet JSON (all required fields present).
const subnetJSON = `{"id":1,"label":"subnet-1","name":"subnet-1","annotations":{},"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z","revision":1,"tags":{},"parentSubnetId":0,"ipVersion":"ipv4","networkAddress":"10.0.0.0","prefixLength":24,"netmask":"255.255.255.0","defaultGatewayAddress":"10.0.0.1","isPool":false,"allocationDenylist":[],"childOverlapAllowRules":[]}`

func subnetListHandler(statusCode int, count, currentPage, totalPages int) http.HandlerFunc {
	items := make([]string, count)
	for i := range items {
		items[i] = subnetJSON
	}
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func subnetGetHandler(statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if statusCode == http.StatusNotFound {
			w.WriteHeader(statusCode)
			fmt.Fprint(w, `{"message":"not found","statusCode":404}`)
			return
		}
		w.WriteHeader(statusCode)
		fmt.Fprint(w, subnetJSON)
	}
}

// subnetPage builds a page of n minimal subnet maps for MultiPageServer.
func subnetPage(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{
			"id": i + 1, "label": fmt.Sprintf("subnet-%d", i+1),
			"name": fmt.Sprintf("subnet-%d", i+1), "annotations": map[string]any{},
			"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
			"revision": 1, "tags": map[string]any{}, "parentSubnetId": 0,
			"ipVersion": "ipv4", "networkAddress": fmt.Sprintf("10.0.%d.0", i),
			"prefixLength": 24, "netmask": "255.255.255.0",
			"defaultGatewayAddress": "10.0.0.1", "isPool": false,
			"allocationDenylist": []any{}, "childOverlapAllowRules": []any{},
		}
	}
	return items
}

func TestSubnetList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets": subnetListHandler(http.StatusOK, 2, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetList(ctx); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSubnetList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetList(ctx); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestSubnetList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets": subnetListHandler(http.StatusOK, 0, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetList(ctx); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestSubnetList_Pagination(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/subnets", []any{subnetPage(100), subnetPage(100), subnetPage(5)})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetList(ctx); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestSubnetGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets/1": testutils.RawHandler(http.StatusOK, subnetJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSubnetGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets/99": subnetGetHandler(http.StatusNotFound),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestSubnetGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestSubnetCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, subnetJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"networkAddress":"10.0.0.0","prefixLength":24,"isPool":false}`)
	if err := SubnetCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSubnetCreate_BadRequest(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"networkAddress":"bad"}`)
	if err := SubnetCreate(ctx, config); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

func TestSubnetDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets/1": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodGet {
				fmt.Fprint(w, subnetJSON)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSubnetDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/subnets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}
