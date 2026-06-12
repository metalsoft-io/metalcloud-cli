package custom_iso

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func sdkCreateIso(label, accessURL string) sdk.CreateCustomIso {
	return sdk.CreateCustomIso{
		Label:     label,
		AccessUrl: accessURL,
	}
}

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeISO(id int) map[string]any {
	return map[string]any{
		"id": id, "label": "ubuntu-22", "name": "Ubuntu 22.04",
		"type": "iso", "isPublic": float64(0), "accessUrl": "https://example.com/ubuntu.iso",
		"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
		"links": []any{},
	}
}

func TestCustomIsoList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeISO(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoList(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCustomIsoList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoList(ctx); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestCustomIsoList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoList(ctx); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestCustomIsoList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeISO(i + 1)
		page2[i] = makeISO(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeISO(i + 201)
	}

	ts := testutils.MultiPageServer("/api/v2/custom-isos", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoList(ctx); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestCustomIsoGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/1": testutils.JSONHandler(200, makeISO(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoGet(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCustomIsoGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoGet(ctx, "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestCustomIsoDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/1": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoDelete(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCustomIsoCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos": testutils.JSONHandler(201, makeISO(3)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	payload := sdkCreateIso("new-iso", "https://example.com/new.iso")
	if err := CustomIsoCreate(ctx, payload); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCustomIsoCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	payload := sdkCreateIso("new-iso", "https://example.com/new.iso")
	if err := CustomIsoCreate(ctx, payload); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestCustomIsoUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/1": testutils.JSONHandler(200, makeISO(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"updated-iso"}`)
	if err := CustomIsoUpdate(ctx, "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCustomIsoUpdate_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"updated-iso"}`)
	if err := CustomIsoUpdate(ctx, "bad", config); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestCustomIsoUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/1": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"updated-iso"}`)
	if err := CustomIsoUpdate(ctx, "1", config); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestCustomIsoMakePublic_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/1/actions/make-public": testutils.JSONHandler(200, makeISO(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoMakePublic(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCustomIsoMakePublic_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoMakePublic(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestCustomIsoMakePublic_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/1/actions/make-public": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoMakePublic(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestCustomIsoBootServer_HappyPath(t *testing.T) {
	jobInfo := map[string]any{
		"jobId":      float64(100),
		"jobGroupId": float64(10),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/custom-isos/1/actions/boot-into-server/2": testutils.JSONHandler(200, jobInfo),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoBootServer(ctx, "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestCustomIsoBootServer_InvalidIsoId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoBootServer(ctx, "bad", "2"); err == nil {
		t.Fatal("expected error for invalid iso id, got nil")
	}
}

func TestCustomIsoBootServer_InvalidServerId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoBootServer(ctx, "1", "bad"); err == nil {
		t.Fatal("expected error for invalid server id, got nil")
	}
}

func TestCustomIsoConfigExample(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CustomIsoConfigExample(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
