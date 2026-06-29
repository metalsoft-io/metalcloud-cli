package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Required: id, revision, name, description, status, kind, definition
func extensionFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "revision": 1, "name": "My Extension",
		"description": "An extension", "status": "active",
		"kind": "workflow", "definition": map[string]interface{}{},
		"links": []interface{}{},
	}
}

// makeInfraSearchResponse mirrors internal/extension_instance infrastructure search fixture.
// Infrastructure model required: label, updatedTimestamp, id, revision, serviceStatus,
// datacenterName, siteId, createdTimestamp, designIsLocked, config
func makeInfraSearchResponse(infraId float64, label string) map[string]interface{} {
	return map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"id": infraId, "label": label, "revision": float64(1),
				"serviceStatus":    "active",
				"datacenterName":   "dc1",
				"siteId":           1,
				"createdTimestamp": "2024-01-01T00:00:00Z",
				"updatedTimestamp": "2024-01-01T00:00:00Z",
				"designIsLocked":   0,
				"config":           map[string]interface{}{},
			},
		},
		"meta": map[string]interface{}{"itemsPerPage": 100},
	}
}

// makeExtensionInstanceFixture mirrors internal/extension_instance test fixture.
// automaticManagement, waitForCompletion, disabled are float32 in Go → send as numbers.
func makeExtensionInstanceFixture(id float64, infraId float64) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "revision": float64(1), "label": "ext-instance",
		"automaticManagement": float64(1),
		"infrastructureId":    infraId,
		"infrastructure": map[string]interface{}{
			"id": infraId, "label": "my-infra",
		},
		"extensionId":      float64(5),
		"serviceStatus":    "active",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"inputVariables":   []interface{}{},
		"outputVariables":  []interface{}{},
		"config": map[string]interface{}{
			"revision":            float64(1),
			"label":               "ext-instance",
			"automaticManagement": float64(1),
			"deployType":          "none",
			"deployStatus":        "idle",
			"updatedTimestamp":    "2024-01-01T00:00:00Z",
		},
		"links": []interface{}{},
	}
}

// --- extension list ---

func TestExtensionList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/extensions", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(extensionFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestExtensionList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/extensions", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(extensionFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestExtensionList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "extension", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- extension-instance list ---
// GetInfrastructureByIdOrLabel calls GET /api/v2/infrastructures?search=<arg>,
// then ExtensionInstanceList calls GET /api/v2/infrastructures/{id}/extension-instances.

func TestExtensionInstanceList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(makeInfraSearchResponse(10, "my-infra"))
		})
		mux.HandleFunc("/api/v2/infrastructures/10/extension-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeExtensionInstanceFixture(1, 10)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "list", "my-infra"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestExtensionInstanceList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(makeInfraSearchResponse(10, "my-infra"))
		})
		mux.HandleFunc("/api/v2/infrastructures/10/extension-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeExtensionInstanceFixture(1, 10)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "ls", "my-infra"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestExtensionInstanceList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "extension-instance", "list", "my-infra"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// extensionDefinitionJSON is a minimal valid ExtensionDefinition JSON satisfying
// all SDK required properties.
const extensionDefinitionJSON = `{
  "kind":"workflow","schemaVersion":"1.0","name":"my-ext","label":"my-ext",
  "extensionType":"workflow","vendor":"test","extensionVersion":"1.0.0",
  "icon":"","dependencies":{"controllerVersion":"1.0"},"inputs":[],"outputs":[],"assets":[]
}`

// --- extension create ---

func TestExtensionCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/extensions", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(extensionFixtureWithDefinition(1))
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ext-def-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(extensionDefinitionJSON)
	f.Close()

	if _, execErr := runCLI(t, srv, "extension", "create", "my-ext", "workflow", "test desc", "--definition-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// extensionFixtureWithDefinition returns an extension fixture that includes a
// valid ExtensionDefinition so SDK unmarshalling succeeds in GetExtensionByIdOrLabel.
func extensionFixtureWithDefinition(id int) map[string]interface{} {
	f := extensionFixture(id)
	f["definition"] = map[string]interface{}{
		"kind": "workflow", "schemaVersion": "1.0", "name": "my-ext", "label": "my-ext",
		"extensionType": "workflow", "vendor": "test", "extensionVersion": "1.0.0",
		"icon": "", "dependencies": map[string]interface{}{"controllerVersion": "1.0"},
		"inputs": []interface{}{}, "outputs": []interface{}{}, "assets": []interface{}{},
	}
	return f
}

// --- extension list formats ---

func TestExtensionList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/extensions", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(extensionFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "extension", "list")
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

func TestExtensionInstanceList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(makeInfraSearchResponse(10, "my-infra"))
		})
		mux.HandleFunc("/api/v2/infrastructures/10/extension-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeExtensionInstanceFixture(1, 10)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "extension-instance", "list", "my-infra")
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

// --- extension update ---

func TestExtensionUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/extensions/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(extensionFixtureWithDefinition(1))
		})
	}))
	defer srv.Close()

	// Pass --definition-source to override any stale flag value from a prior test.
	f, err := os.CreateTemp(t.TempDir(), "ext-upd-def-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(extensionDefinitionJSON)
	f.Close()

	if _, execErr := runCLI(t, srv, "extension", "update", "1", "new-name", "--definition-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}
