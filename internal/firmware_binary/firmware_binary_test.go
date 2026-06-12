package firmware_binary

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func binaryFixture(id int) map[string]any {
	return map[string]any{
		"id":                     id,
		"name":                   "BIOS-2.15.0",
		"catalogId":              10,
		"vendorDownloadUrl":      "https://example.com/bios.bin",
		"rebootRequired":         true,
		"updateSeverity":         "recommended",
		"vendorSupportedDevices": []map[string]any{},
		"vendorSupportedSystems": []map[string]any{},
		"vendor":                 map[string]any{},
		"links":                  []any{},
	}
}

func TestFirmwareBinaryList_HappyPath(t *testing.T) {
	page1 := []map[string]any{
		binaryFixture(1),
		binaryFixture(2),
	}
	srv := testutils.MultiPageServer("/api/v2/firmware/binary", []any{page1})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryList(ctx); err != nil {
		t.Fatalf("FirmwareBinaryList() unexpected error: %v", err)
	}
}

func TestFirmwareBinaryList_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryList(ctx); err == nil {
		t.Fatal("FirmwareBinaryList() expected error, got nil")
	}
}

func TestFirmwareBinaryList_Empty(t *testing.T) {
	srv := testutils.MultiPageServer("/api/v2/firmware/binary", []any{[]map[string]any{}})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryList(ctx); err != nil {
		t.Fatalf("FirmwareBinaryList() unexpected error on empty: %v", err)
	}
}

func TestFirmwareBinaryList_Pagination(t *testing.T) {
	makeItems := func(start, count int) []map[string]any {
		items := make([]map[string]any, count)
		for i := range items {
			items[i] = binaryFixture(start + i)
		}
		return items
	}

	page1 := makeItems(1, 100)
	page2 := makeItems(101, 100)
	page3 := makeItems(201, 5)

	srv := testutils.MultiPageServer("/api/v2/firmware/binary", []any{page1, page2, page3})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryList(ctx); err != nil {
		t.Fatalf("FirmwareBinaryList() pagination error: %v", err)
	}
}

func TestFirmwareBinaryGet_HappyPath(t *testing.T) {
	binary := binaryFixture(8)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/8": testutils.JSONHandler(200, binary),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryGet(ctx, "8"); err != nil {
		t.Fatalf("FirmwareBinaryGet() unexpected error: %v", err)
	}
}

func TestFirmwareBinaryGet_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryGet(ctx, "99"); err == nil {
		t.Fatal("FirmwareBinaryGet() expected error for not found, got nil")
	}
}

func TestFirmwareBinaryGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwareBinaryGet(ctx, "not-a-number"); err == nil {
		t.Fatal("FirmwareBinaryGet() expected error for invalid ID, got nil")
	}
}

func TestFirmwareBinaryCreate_HappyPath(t *testing.T) {
	created := binaryFixture(20)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary": testutils.JSONHandler(201, created),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"name":"BIOS-2.15.0","catalogId":1,"vendorDownloadUrl":"https://example.com/bios.bin","rebootRequired":true,"updateSeverity":"recommended","vendorSupportedDevices":[],"vendorSupportedSystems":[]}`)
	if err := FirmwareBinaryCreate(ctx, config); err != nil {
		t.Fatalf("FirmwareBinaryCreate() unexpected error: %v", err)
	}
}

func TestFirmwareBinaryCreate_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary": testutils.ErrorHandler(400, "bad request"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"name":"BIOS-2.15.0"}`)
	if err := FirmwareBinaryCreate(ctx, config); err == nil {
		t.Fatal("FirmwareBinaryCreate() expected error, got nil")
	}
}

func TestFirmwareBinaryDelete_HappyPath(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/4": testutils.RawHandler(204, ""),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryDelete(ctx, "4"); err != nil {
		t.Fatalf("FirmwareBinaryDelete() unexpected error: %v", err)
	}
}

func TestFirmwareBinaryDelete_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/4": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBinaryDelete(ctx, "4"); err == nil {
		t.Fatal("FirmwareBinaryDelete() expected error, got nil")
	}
}

func TestFirmwareBinaryConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwareBinaryConfigExample(ctx); err != nil {
		t.Errorf("FirmwareBinaryConfigExample() unexpected error: %v", err)
	}
}
