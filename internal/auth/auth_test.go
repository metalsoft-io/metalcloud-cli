package auth

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// authConfigWithMappings returns a minimal /api/v2/configuration response that
// contains an LDAP group mapping so AuthLdapMappingList has data to print.
func authConfigWithMappings() map[string]any {
	return map[string]any{
		"auth": map[string]any{
			"ldap": map[string]any{
				"groupsMapping": []any{
					map[string]any{
						"groupName":              "cn=admins,dc=example,dc=com",
						"roleName":               "admin",
						"priority":               float64(1),
						"userExternalIdentifier": "objectGUID",
						"username":               "sAMAccountName",
						"email":                  "mail",
					},
				},
			},
		},
	}
}

func TestAuthLdapMappingList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.JSONHandler(200, authConfigWithMappings()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AuthLdapMappingList(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestAuthLdapMappingList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AuthLdapMappingList(ctx); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestAuthLdapMappingList_EmptyMappings(t *testing.T) {
	body := map[string]any{
		"auth": map[string]any{
			"ldap": map[string]any{
				"groupsMapping": []any{},
			},
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// AuthLdapMappingList builds the list locally — should return nil even when empty.
	if err := AuthLdapMappingList(ctx); err != nil {
		t.Fatalf("expected nil error on empty mappings, got: %v", err)
	}
}

// AuthLdapMappingAdd — GET then PATCH /api/v2/config/auth

func TestAuthLdapMappingAdd_HappyPath(t *testing.T) {
	// After adding, the updated config is returned from the PATCH.
	updatedConfig := map[string]any{
		"ldap": map[string]any{
			"groupsMapping": []any{
				map[string]any{
					"groupName":              "cn=newgroup,dc=example,dc=com",
					"roleName":               "operator",
					"priority":               float64(2),
					"userExternalIdentifier": "objectGUID",
					"username":               "sAMAccountName",
					"email":                  "mail",
				},
			},
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		// GET /api/v2/config?filter=auth
		"/api/v2/config": testutils.JSONHandler(200, authConfigWithMappings()),
		// PATCH /api/v2/config/auth
		"/api/v2/config/auth": testutils.JSONHandler(200, updatedConfig),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	opts := AuthLdapMappingOptions{RoleName: "operator", Priority: 2}
	if err := AuthLdapMappingAdd(ctx, "cn=newgroup,dc=example,dc=com", opts); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestAuthLdapMappingAdd_DuplicateGroup(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.JSONHandler(200, authConfigWithMappings()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// Use the group name that already exists in authConfigWithMappings.
	opts := AuthLdapMappingOptions{RoleName: "admin", Priority: 1}
	if err := AuthLdapMappingAdd(ctx, "cn=admins,dc=example,dc=com", opts); err == nil {
		t.Fatal("expected error for duplicate group, got nil")
	}
}

func TestAuthLdapMappingAdd_GetError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	opts := AuthLdapMappingOptions{RoleName: "operator", Priority: 2}
	if err := AuthLdapMappingAdd(ctx, "cn=newgroup,dc=example,dc=com", opts); err == nil {
		t.Fatal("expected error when GET fails, got nil")
	}
}

// AuthLdapMappingUpdate

func TestAuthLdapMappingUpdate_HappyPath(t *testing.T) {
	updatedConfig := map[string]any{
		"ldap": map[string]any{
			"groupsMapping": []any{
				map[string]any{
					"groupName":              "cn=admins,dc=example,dc=com",
					"roleName":               "superadmin",
					"priority":               float64(5),
					"userExternalIdentifier": "objectGUID",
					"username":               "sAMAccountName",
					"email":                  "mail",
				},
			},
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config":      testutils.JSONHandler(200, authConfigWithMappings()),
		"/api/v2/config/auth": testutils.JSONHandler(200, updatedConfig),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	opts := AuthLdapMappingOptions{RoleName: "superadmin", Priority: 5}
	if err := AuthLdapMappingUpdate(ctx, "cn=admins,dc=example,dc=com", opts); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestAuthLdapMappingUpdate_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.JSONHandler(200, authConfigWithMappings()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	opts := AuthLdapMappingOptions{RoleName: "admin", Priority: 1}
	if err := AuthLdapMappingUpdate(ctx, "cn=nonexistent,dc=example,dc=com", opts); err == nil {
		t.Fatal("expected error for missing group, got nil")
	}
}

func TestAuthLdapMappingUpdate_GetError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	opts := AuthLdapMappingOptions{RoleName: "admin", Priority: 1}
	if err := AuthLdapMappingUpdate(ctx, "cn=admins,dc=example,dc=com", opts); err == nil {
		t.Fatal("expected error when GET fails, got nil")
	}
}

// AuthLdapMappingRemove

func TestAuthLdapMappingRemove_HappyPath(t *testing.T) {
	updatedConfig := map[string]any{
		"ldap": map[string]any{
			"groupsMapping": []any{},
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config":      testutils.JSONHandler(200, authConfigWithMappings()),
		"/api/v2/config/auth": testutils.JSONHandler(200, updatedConfig),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AuthLdapMappingRemove(ctx, "cn=admins,dc=example,dc=com"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestAuthLdapMappingRemove_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.JSONHandler(200, authConfigWithMappings()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AuthLdapMappingRemove(ctx, "cn=nonexistent,dc=example,dc=com"); err == nil {
		t.Fatal("expected error for missing group, got nil")
	}
}

func TestAuthLdapMappingRemove_GetError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := AuthLdapMappingRemove(ctx, "cn=admins,dc=example,dc=com"); err == nil {
		t.Fatal("expected error when GET fails, got nil")
	}
}
