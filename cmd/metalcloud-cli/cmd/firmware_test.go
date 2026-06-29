package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// firmware_test.go covers:
//   firmware-catalog list / ls
//   firmware-baseline list / ls
//   firmware-binary list / ls

func firmwareCatalogFixture() map[string]interface{} {
	return map[string]interface{}{
		"id": 1, "name": "Dell Catalog", "vendor": "dell",
		"updateType": "online", "createdTimestamp": "2024-01-01T00:00:00Z",
		"links": []interface{}{},
	}
}

func firmwareBaselineFixture() map[string]interface{} {
	return map[string]interface{}{
		"id": 1, "name": "Baseline 1", "links": []interface{}{},
	}
}

// vendorSupportedDevices/Systems must be slices (not empty maps),
// vendor must be a map (not string), disabled must be float32 (0/1).
func firmwareBinaryFixture() map[string]interface{} {
	return map[string]interface{}{
		"id": 1, "name": "BIOS-2.15.0", "catalogId": 10,
		"vendorDownloadUrl":      "https://example.com/bios.bin",
		"rebootRequired":         true, "updateSeverity": "recommended",
		"vendorSupportedDevices": []interface{}{},
		"vendorSupportedSystems": []interface{}{},
		"vendor":                 map[string]interface{}{},
		"links":                  []interface{}{},
	}
}

func writePagedJSON(w http.ResponseWriter, item interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(paginatedList(item))
}

// --- firmware-catalog ---

func TestFirmwareCatalogList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/catalog", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareCatalogFixture())
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "firmware-catalog", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestFirmwareCatalogList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/catalog", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareCatalogFixture())
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "firmware-catalog", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestFirmwareCatalogList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "firmware-catalog", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- firmware-baseline ---

func TestFirmwareBaselineList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/baseline", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareBaselineFixture())
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "firmware-baseline", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestFirmwareBaselineList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/baseline", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareBaselineFixture())
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "firmware-baseline", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestFirmwareBaselineList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "firmware-baseline", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- firmware-binary ---

func TestFirmwareBinaryList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/binary", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareBinaryFixture())
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "firmware-binary", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestFirmwareBinaryList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/binary", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareBinaryFixture())
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "firmware-binary", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestFirmwareBinaryList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "firmware-binary", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- format tests ---

func TestFirmwareCatalogList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/catalog", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareCatalogFixture())
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "firmware-catalog", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

func TestFirmwareBaselineList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/baseline", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareBaselineFixture())
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "firmware-baseline", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

func TestFirmwareBinaryList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/binary", func(w http.ResponseWriter, r *http.Request) {
			writePagedJSON(w, firmwareBinaryFixture())
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "firmware-binary", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

// --- firmware-catalog create / delete ---

func newFirmwareCatalogWriteServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/catalog/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			writePagedJSON(w, firmwareCatalogFixture())
		})
		mux.HandleFunc("/api/v2/firmware/catalog", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(firmwareCatalogFixture())
				return
			}
			writePagedJSON(w, firmwareCatalogFixture())
		})
	}))
}

func TestFirmwareCatalogCreate_MissingVendor(t *testing.T) {
	srv := newFirmwareCatalogWriteServer()
	defer srv.Close()

	// create with invalid vendor fails validation before any API call
	_, err := runCLI(t, srv, "firmware-catalog", "create",
		"--name", "Test Catalog",
		"--vendor", "invalid-vendor",
		"--update-type", "online",
	)
	if err == nil {
		t.Fatal("expected error for invalid vendor, got nil")
	}
}

func TestFirmwareCatalogDelete(t *testing.T) {
	srv := newFirmwareCatalogWriteServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "firmware-catalog", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- firmware-baseline create / delete ---

func newFirmwareBaselineWriteServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/firmware/baseline/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			_ = json.NewEncoder(w).Encode(firmwareBaselineFixture())
		})
		mux.HandleFunc("/api/v2/firmware/baseline", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(firmwareBaselineFixture())
				return
			}
			writePagedJSON(w, firmwareBaselineFixture())
		})
	}))
}

func TestFirmwareBaselineCreate(t *testing.T) {
	srv := newFirmwareBaselineWriteServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "baseline-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"name":"Baseline 1","catalogIds":[1]}`)
	f.Close()

	_, err = runCLI(t, srv, "firmware-baseline", "create", "--config-source", f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFirmwareBaselineDelete(t *testing.T) {
	srv := newFirmwareBaselineWriteServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "firmware-baseline", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

