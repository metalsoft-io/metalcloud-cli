package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// userListFixture returns a User with all required fields populated.
// It reuses the same shape as fullUserJSON but as a Go map for use in
// paginatedList() responses from the /api/v2/users endpoint.
func userListFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "revision": 1, "email": "user@example.com",
		"displayName": "Example User", "emailStatus": "active",
		"language": "en", "brand": "default", "isBrandManager": false,
		"lastLoginTimestamp": "2024-01-01T00:00:00Z", "lastLoginType": "password",
		"isBlocked": false, "passwordChangeRequired": false, "accessLevel": "user",
		"isBillable": true, "isTestingMode": false, "authenticatorMustChange": false,
		"authenticatorCreatedTimestamp": "2024-01-01T00:00:00Z",
		"excludeFromReports": false, "isTestAccount": false,
		"isArchived": false, "isDatastorePublisher": false, "provider": "local",
		"passwordLastChangedTimestamp": "2024-01-01T00:00:00Z",
		"franchise": "default", "createdTimestamp": "2024-01-01T00:00:00Z",
		"planType": "default", "isSuspended": false, "authenticatorEnabled": false,
		"config": map[string]interface{}{
			"revision": 1, "displayName": "Example User", "emailStatus": "active",
			"language": "en", "brand": "default", "isBrandManager": false,
			"lastLoginTimestamp": "2024-01-01T00:00:00Z", "lastLoginType": "password",
			"isBlocked": false, "passwordChangeRequired": false, "accessLevel": "user",
			"isBillable": true, "isTestingMode": false, "authenticatorMustChange": false,
			"authenticatorCreatedTimestamp": "2024-01-01T00:00:00Z",
			"excludeFromReports": false, "isTestAccount": false,
			"isArchived": false, "isDatastorePublisher": false, "provider": "local",
			"passwordLastChangedTimestamp": "2024-01-01T00:00:00Z",
		},
		"meta":        map[string]interface{}{},
		"permissions": map[string]interface{}{"rolePermissions": []interface{}{}},
		"links":       []interface{}{},
	}
}

func TestUserList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(userListFixture(2)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(userListFixture(2)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestUserList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "user", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

func TestUserList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(userListFixture(2)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "user", "list")
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

func TestUserGet_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "get", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserGet_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "user", "get", "1"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- user create ---

func TestUserCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(userListFixture(2))
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "user", "create", "--email", "new@example.com", "--password", "secret123"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- user suspend ---

func TestUserSuspend(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/2", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(userListFixture(2))
		})
		mux.HandleFunc("/api/v2/users/2/actions/suspend", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id": 1, "userId": 2, "type": "suspend", "publicComment": "testing",
				"createdTimestamp": "2024-01-01T00:00:00Z",
			})
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "user", "suspend", "2", "--reason", "testing"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- user archive ---

func TestUserArchive(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/2", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(userListFixture(2))
		})
		mux.HandleFunc("/api/v2/users/2/actions/archive", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(userListFixture(2))
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "user", "archive", "2"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}
