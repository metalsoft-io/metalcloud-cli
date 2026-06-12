package role

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func roleItem(id, name, label string) map[string]interface{} {
	return map[string]interface{}{
		"id":            id,
		"name":          name,
		"label":         label,
		"type":          "custom",
		"description":   "test role",
		"usersWithRole": 0,
		"permissions":   []interface{}{},
	}
}

func TestRoleList_HappyPath(t *testing.T) {
	resp := map[string]interface{}{
		"roles": []interface{}{
			roleItem("1", "admin", "Admin"),
			roleItem("2", "user", "User"),
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/roles": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := List(ctx); err != nil {
		t.Errorf("RoleList: expected nil error, got: %v", err)
	}
}

func TestRoleList_Empty(t *testing.T) {
	resp := map[string]interface{}{
		"roles": []interface{}{},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/roles": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := List(ctx); err != nil {
		t.Errorf("RoleList empty: expected nil error, got: %v", err)
	}
}

func TestRoleList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/roles": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := List(ctx); err == nil {
		t.Error("RoleList with 500: expected error, got nil")
	}
}

func TestRoleGet_Success(t *testing.T) {
	item := roleItem("1", "admin", "Admin")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/roles/admin": testutils.JSONHandler(http.StatusOK, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Get(ctx, "admin"); err != nil {
		t.Errorf("RoleGet: expected nil error, got: %v", err)
	}
}

func TestRoleGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/roles/nonexistent": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Get(ctx, "nonexistent"); err == nil {
		t.Error("RoleGet not-found: expected error, got nil")
	}
}

func TestRoleDelete_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/roles/myrole": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.NotFound(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Delete(ctx, "myrole"); err != nil {
		t.Errorf("RoleDelete: expected nil error, got: %v", err)
	}
}

func TestRoleDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/roles/ghost": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := Delete(ctx, "ghost"); err == nil {
		t.Error("RoleDelete not-found: expected error, got nil")
	}
}
