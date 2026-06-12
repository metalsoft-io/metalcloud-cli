package bucket

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func infraHandler() http.HandlerFunc {
	return testutils.JSONHandler(200, map[string]any{
		"data": []any{
			map[string]any{
				"id": float64(123), "label": "test-infra",
				"serviceStatus": "active", "revision": float64(1),
				"datacenterName": "dc1", "siteId": float64(1),
				"designIsLocked": float64(0),
				"createdTimestamp": "2024-01-01T00:00:00Z",
				"updatedTimestamp": "2024-01-01T00:00:00Z",
				"config": map[string]any{}, "links": []any{},
			},
		},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	})
}

func makeBucket(id int) map[string]any {
	return map[string]any{
		"id": id, "label": "bucket-1", "sizeGB": float64(100),
		"infrastructureId": float64(123), "serviceStatus": "active",
		"revision": float64(1), "subdomain": "bucket-1.example.com",
		"subdomainPermanent": "bucket-1.example.com",
		"dnsSubdomainId": float64(1),
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"infrastructure": map[string]any{"id": float64(123)},
		"config": map[string]any{
			"revision": float64(1), "sizeGB": float64(100),
			"updatedTimestamp": "2024-01-01T00:00:00Z",
			"label": "bucket-1", "subdomain": "bucket-1.example.com",
			"deployType": "deploy", "deployStatus": "not_started",
		},
		"meta": map[string]any{"name": "bucket-1"},
		"links": []any{},
	}
}

func TestBucketList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeBucket(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":             infraHandler(),
		"/api/v2/infrastructures/123/buckets": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketList(ctx, "123", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":             infraHandler(),
		"/api/v2/infrastructures/123/buckets": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketList(ctx, "123", nil); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestBucketList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":             infraHandler(),
		"/api/v2/infrastructures/123/buckets": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketList(ctx, "123", nil); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestBucketList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeBucket(i + 1)
		page2[i] = makeBucket(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeBucket(i + 201)
	}

	callCount := 0
	pages := [][]any{page1, page2, page3}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/buckets": func(w http.ResponseWriter, r *http.Request) {
			idx := callCount
			if idx >= len(pages) {
				idx = len(pages) - 1
			}
			callCount++
			resp := testutils.PaginatedResponse(pages[idx], int32(idx+1), int32(len(pages)))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketList(ctx, "123", nil); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestBucketGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":               infraHandler(),
		"/api/v2/infrastructures/123/buckets/1": testutils.JSONHandler(200, makeBucket(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketGet(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                 infraHandler(),
		"/api/v2/infrastructures/123/buckets/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketGet(ctx, "123", "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestBucketCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":             infraHandler(),
		"/api/v2/infrastructures/123/buckets": testutils.JSONHandler(201, makeBucket(2)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"bucket-2","sizeGB":200,"storagePoolId":1}`)
	if err := BucketCreate(ctx, "123", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketDelete_HappyPath(t *testing.T) {
	// BucketDelete calls getBucketIdAndRevision (GET) then DELETE on the same path.
	getHandler := testutils.JSONHandler(200, makeBucket(1))
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/buckets/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			getHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketDelete(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketUpdateConfig_HappyPath(t *testing.T) {
	// BucketUpdateConfig calls getBucketIdAndRevision (GET on /buckets/1) then PATCH on /buckets/1/config.
	// BucketConfiguration requires: revision, label, sizeGB, deployType, deployStatus, updatedTimestamp.
	configBody := map[string]any{
		"revision": float64(1), "label": "bucket-updated", "sizeGB": float64(100),
		"subdomain":        "bucket-1.example.com",
		"deployType": "deploy", "deployStatus": "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                        infraHandler(),
		"/api/v2/infrastructures/123/buckets/1":          testutils.JSONHandler(200, makeBucket(1)),
		"/api/v2/infrastructures/123/buckets/1/config":   testutils.JSONHandler(200, configBody),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"bucket-updated"}`)
	if err := BucketUpdateConfig(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketUpdateConfig_GetError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":               infraHandler(),
		"/api/v2/infrastructures/123/buckets/1": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"bucket-updated"}`)
	if err := BucketUpdateConfig(ctx, "123", "1", config); err == nil {
		t.Fatal("expected error when GET fails, got nil")
	}
}

func TestBucketUpdateMeta_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                     infraHandler(),
		"/api/v2/infrastructures/123/buckets/1/meta":  testutils.JSONHandler(200, makeBucket(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"name":"new-name"}`)
	if err := BucketUpdateMeta(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketGetConfigInfo_HappyPath(t *testing.T) {
	// BucketConfiguration requires: revision, label, sizeGB, deployType, deployStatus, updatedTimestamp.
	configInfo := map[string]any{
		"revision":         float64(1),
		"label":            "bucket-1",
		"sizeGB":           float64(100),
		"subdomain":        "bucket-1.example.com",
		"deployType":       "deploy",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                       infraHandler(),
		"/api/v2/infrastructures/123/buckets/1/config":  testutils.JSONHandler(200, configInfo),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketGetConfigInfo(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketGetConfigInfo_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                      infraHandler(),
		"/api/v2/infrastructures/123/buckets/1/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketGetConfigInfo(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestBucketGetCredentials_HappyPath(t *testing.T) {
	creds := map[string]any{
		"accessKeyId": "AKIAIOSFODNN7EXAMPLE",
		"secretKey":   "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"endpoint":    "https://s3.example.com",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                            infraHandler(),
		"/api/v2/infrastructures/123/buckets/1/credentials": testutils.JSONHandler(200, creds),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketGetCredentials(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestBucketGetCredentials_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                            infraHandler(),
		"/api/v2/infrastructures/123/buckets/1/credentials": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := BucketGetCredentials(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}
