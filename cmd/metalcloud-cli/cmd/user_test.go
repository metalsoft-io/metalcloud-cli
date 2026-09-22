package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
		"excludeFromReports":            false, "isTestAccount": false,
		"isArchived": false, "isDatastorePublisher": false, "provider": "local",
		"passwordLastChangedTimestamp": "2024-01-01T00:00:00Z",
		"franchise":                    "default", "createdTimestamp": "2024-01-01T00:00:00Z",
		"planType": "default", "isSuspended": false, "authenticatorEnabled": false,
		"config": map[string]interface{}{
			"revision": 1, "displayName": "Example User", "emailStatus": "active",
			"language": "en", "brand": "default", "isBrandManager": false,
			"lastLoginTimestamp": "2024-01-01T00:00:00Z", "lastLoginType": "password",
			"isBlocked": false, "passwordChangeRequired": false, "accessLevel": "user",
			"isBillable": true, "isTestingMode": false, "authenticatorMustChange": false,
			"authenticatorCreatedTimestamp": "2024-01-01T00:00:00Z",
			"excludeFromReports":            false, "isTestAccount": false,
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

// --- user configuration, SSH key and suspend reasons ---

// usrSshKeyFixture returns a UserSSHKeys payload with all required fields.
func usrSshKeyFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "userId": 1, "sshKey": "ssh-rsa AAAA...",
		"status": "active", "createdTimestamp": "2024-01-01T00:00:00Z",
	}
}

// usrSuspendReasonFixture returns a UserSuspendReason payload with all
// required fields.
func usrSuspendReasonFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "userId": 1, "type": "admin",
		"publicComment": "suspended", "createdTimestamp": "2024-01-01T00:00:00Z",
	}
}

// usrUserInfoFixture returns a UserInfo payload with all required fields, as
// returned by the delegate endpoints.
func usrUserInfoFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "revision": 1, "displayName": "Example User",
		"email": "user@example.com", "emailStatus": "active",
		"language": "en", "accessLevel": "user", "isArchived": false,
		"lastLoginTimestamp": "2024-01-01T00:00:00Z",
		"createdTimestamp":   "2024-01-01T00:00:00Z",
	}
}

func TestUserConfigGet(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/config", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(1)["config"])
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "config", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserUpdateMeta(t *testing.T) {
	gotMethod := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/meta", func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			jsonResponse(w, http.StatusOK, map[string]interface{}{"guiSettings": map[string]interface{}{}})
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "meta-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"guiSettings":{}}`)
	f.Close()

	if _, err := runCLI(t, srv, "user", "update-meta", "1", "--config-source", f.Name()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", gotMethod)
	}
}

func TestUserSshKeyGet(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/ssh-keys/5", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, usrSshKeyFixture(5))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "ssh-key", "1", "5"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserSshKeyGet_RequiresTwoArgs(t *testing.T) {
	if _, err := runCLI(t, nil, "user", "ssh-key", "1"); err == nil {
		t.Fatal("expected error when the key ID is missing")
	}
}

func TestUserSuspendReasons(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/suspend-reasons", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"data": []interface{}{usrSuspendReasonFixture(3)},
			})
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "suspend-reasons", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// --- user delegates ---

func TestUserDelegateParents(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/parent-delegates", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"data": []interface{}{usrUserInfoFixture(2)},
			})
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "delegate", "parents", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserDelegateChildren(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/child-delegates", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"data": []interface{}{usrUserInfoFixture(3)},
			})
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "delegate", "children", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserDelegateAdd(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/actions/add-delegate/2", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "delegate", "add", "1", "2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserDelegateRemove(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/actions/remove-delegate/2", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "delegate", "rm", "1", "2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// --- user notifications ---

// TestUserResendEmailVerification also guards the SDK gotcha where an unset
// optional body is sent as a literal `null`, which the API rejects.
func TestUserResendEmailVerification(t *testing.T) {
	gotBody := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/actions/resend-email-verification", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			gotBody = strings.TrimSpace(string(body))
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "resend-email-verification", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gotBody == "" || gotBody == "null" {
		t.Errorf("expected a JSON object body, got %q", gotBody)
	}
}

func TestUserResendInvitation(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/actions/resend-user-invitation", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "resend-invitation", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserSendPasswordReset(t *testing.T) {
	gotRedirect := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1/actions/send-password-reset", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			gotRedirect = string(body)
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "send-password-reset", "1", "--redirect-url", "https://example.com/login"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(gotRedirect, "https://example.com/login") {
		t.Errorf("expected the redirect URL in the body, got %q", gotRedirect)
	}
}

// --- user delete ---

func TestUserDelete(t *testing.T) {
	gotIfMatch := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/2", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(2))
		})
		mux.HandleFunc("/api/v2/users/2/actions/delete", func(w http.ResponseWriter, r *http.Request) {
			gotIfMatch = r.Header.Get("If-Match")
			jsonResponse(w, http.StatusOK, userListFixture(2))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "delete", "2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gotIfMatch != "1" {
		t.Errorf("expected If-Match 1, got %q", gotIfMatch)
	}
}

// --- user self-service ---

func TestUserPermissions(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/user/permissions", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"permissions": []interface{}{"users_read", "servers_read"},
			})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "user", "permissions")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "users_read") {
		t.Errorf("expected the permissions in the output, got: %s", out)
	}
}

func TestUserChangePassword(t *testing.T) {
	gotBody := ""
	gotIfMatch := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
		mux.HandleFunc("/api/v2/user/actions/change-password", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			gotBody = string(body)
			gotIfMatch = r.Header.Get("If-Match")
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "change-password", "--current-password", "old", "--new-password", "new"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(gotBody, "newPassword") || !strings.Contains(gotBody, "oldPassword") {
		t.Errorf("expected both passwords in the body, got %q", gotBody)
	}
	if gotIfMatch != "1" {
		t.Errorf("expected If-Match 1, got %q", gotIfMatch)
	}
}

func TestUserChangePassword_RequiresFlags(t *testing.T) {
	if _, err := runCLI(t, nil, "user", "change-password"); err == nil {
		t.Fatal("expected error when neither --new-password nor --config-source is given")
	}
}

func TestUserInitiatePasswordReset(t *testing.T) {
	gotBody := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/user/actions/initiate-password-reset", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			gotBody = string(body)
			w.WriteHeader(http.StatusNoContent)
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "initiate-password-reset", "--email", "user@example.com"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(gotBody, "user@example.com") {
		t.Errorf("expected the e-mail address in the body, got %q", gotBody)
	}
}

func TestUserInitiateEmailChange(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/user/actions/initiate-email-change", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "initiate-email-change", "--email", "new@example.com"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUserRegenerateJwtSalt(t *testing.T) {
	gotIfMatch := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/users/1", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
		mux.HandleFunc("/api/v2/user/actions/regenerate-jwt-salt", func(w http.ResponseWriter, r *http.Request) {
			gotIfMatch = r.Header.Get("If-Match")
			jsonResponse(w, http.StatusOK, userListFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "regenerate-jwt-salt"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gotIfMatch != "1" {
		t.Errorf("expected If-Match 1, got %q", gotIfMatch)
	}
}

func TestUserVerifyEmail(t *testing.T) {
	gotToken := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/user/actions/verify-email", func(w http.ResponseWriter, r *http.Request) {
			gotToken = r.URL.Query().Get("token")
			w.WriteHeader(http.StatusNoContent)
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "verify-email", "--token", "tok-1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gotToken != "tok-1" {
		t.Errorf("expected the token in the query string, got %q", gotToken)
	}
}

func TestUserResetPassword(t *testing.T) {
	gotToken := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/user/actions/reset-password", func(w http.ResponseWriter, r *http.Request) {
			gotToken = r.URL.Query().Get("token")
			w.WriteHeader(http.StatusNoContent)
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "user", "reset-password", "--token", "tok-2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gotToken != "tok-2" {
		t.Errorf("expected the token in the query string, got %q", gotToken)
	}
}

func TestUserResetPassword_RequiresToken(t *testing.T) {
	if _, err := runCLI(t, nil, "user", "reset-password"); err == nil {
		t.Fatal("expected error when --token is missing")
	}
}
