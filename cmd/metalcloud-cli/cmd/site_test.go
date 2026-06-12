package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// siteTestItem satisfies sdk.Site required fields: id, revision, slug, name.
var siteTestItem = map[string]interface{}{
	"id":       1.0,
	"revision": 1.0,
	"slug":     "dc-01",
	"name":     "DC-01",
	"label":    "dc-01",
}

func newSiteTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// SiteList and SiteGet (GetSiteByIdOrLabel) both call GET /api/v2/sites.
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(siteTestItem))
		})
	}))
}

// --- site list ---

func TestSiteList_HappyPath(t *testing.T) {
	srv := newSiteTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "site", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "DC-01") {
		t.Fatalf("expected DC-01 in output, got: %s", out)
	}
}

func TestSiteList_Alias(t *testing.T) {
	srv := newSiteTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestSiteList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "site", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- site get ---

func TestSiteGet_HappyPath(t *testing.T) {
	srv := newSiteTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "site", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "DC-01") {
		t.Fatalf("expected DC-01 in output, got: %s", out)
	}
}

func TestSiteGet_NoArgs(t *testing.T) {
	srv := newSiteTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "get"); err == nil {
		t.Fatal("expected error when no args given to site get")
	}
}

func TestSiteCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, siteTestItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(siteTestItem))
		})
	}))
	defer srv.Close()

	_, err := runCLI(t, srv, "site", "create", "DC-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSiteUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(siteTestItem))
		})
		mux.HandleFunc("/api/v2/sites/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, siteTestItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(siteTestItem)
		})
	}))
	defer srv.Close()

	// site update site_id_or_name [new_label]
	_, err := runCLI(t, srv, "site", "update", "1", "DC-01-updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSiteList_Formats(t *testing.T) {
	srv := newSiteTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "site", "list")
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
