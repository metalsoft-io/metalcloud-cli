package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var bucketItem = map[string]interface{}{
	"id":                 1.0,
	"revision":           1.0,
	"label":              "test-bucket",
	"subdomain":          "test-bucket.example.com",
	"subdomainPermanent": "test-bucket.example.com",
	"dnsSubdomainId":     0.0,
	"sizeGB":             10.0,
	"infrastructureId":   1.0,
	"infrastructure":     map[string]interface{}{"id": 1.0},
	"serviceStatus":      "active",
	"createdTimestamp":   "2024-01-01T00:00:00Z",
	"updatedTimestamp":   "2024-01-01T00:00:00Z",
	"config": map[string]interface{}{
		"revision":         1.0,
		"sizeGB":           10.0,
		"label":            "test-bucket",
		"subdomain":        "test-bucket.example.com",
		"deployType":       "soft",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	},
	"meta": map[string]interface{}{"name": "test-bucket"},
}

func newBucketTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		// Exact path for bucket get (must be registered before the prefix)
		mux.HandleFunc("/api/v2/infrastructures/1/buckets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(bucketItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/buckets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(bucketItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestBucketList(t *testing.T) {
	srv := newBucketTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "bucket", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-bucket") {
		t.Errorf("expected output to contain 'test-bucket', got: %s", out)
	}
}

func TestBucketGet(t *testing.T) {
	srv := newBucketTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "bucket", "get", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-bucket") {
		t.Errorf("expected output to contain 'test-bucket', got: %s", out)
	}
}

func TestBucketListRequiresArg(t *testing.T) {
	srv := newBucketTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "bucket", "list")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func newBucketWriteTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/buckets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			_ = json.NewEncoder(w).Encode(bucketItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/buckets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(bucketItem)
				return
			}
			_ = json.NewEncoder(w).Encode(paginatedList(bucketItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestBucketCreate(t *testing.T) {
	srv := newBucketWriteTestServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "bucket-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"label":"test-bucket","sizeGB":10,"storagePoolId":1}`)
	f.Close()

	_, err = runCLI(t, srv, "bucket", "create", "1", "--config-source", f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBucketDelete(t *testing.T) {
	srv := newBucketWriteTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "bucket", "delete", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBucketList_Formats(t *testing.T) {
	srv := newBucketTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "bucket", "list", "1")
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
