package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newAuthTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/config", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"auth": map[string]interface{}{
					"ldap": map[string]interface{}{
						"server":        map[string]interface{}{},
						"groupsMapping": []interface{}{},
						"profileMapping": map[string]interface{}{
							"userExternalIdentifier": map[string]interface{}{},
							"username":               map[string]interface{}{},
							"email":                  map[string]interface{}{},
							"role":                   map[string]interface{}{},
						},
					},
				},
			})
		})
	})
	return httptest.NewServer(mux)
}

func TestAuthHelp(t *testing.T) {
	out, err := runCLI(t, nil, "auth", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "auth") {
		t.Errorf("expected help output to contain 'auth', got: %s", out)
	}
}

func TestAuthLdapHelp(t *testing.T) {
	out, err := runCLI(t, nil, "auth", "ldap", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(out), "ldap") {
		t.Errorf("expected help output to contain 'ldap', got: %s", out)
	}
}

func TestAuthLdapMappingListAlias(t *testing.T) {
	srv := newAuthTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "authentication", "ldap", "mapping-ls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
