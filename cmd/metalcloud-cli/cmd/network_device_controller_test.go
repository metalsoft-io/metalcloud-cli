package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var ndcItem = map[string]interface{}{
	"id": "23", "revision": 3.0, "status": "active", "siteId": 1.0,
	"identifierString": "ndfc-controller-01", "description": nil,
	"datacenterName": "dc1", "managementAddress": "10.0.0.50", "managementPort": 443.0,
	"username": "admin", "driver": "cisco_ndfc", "tags": []interface{}{},
}

var ndcCredentialsItem = map[string]interface{}{
	"username": "admin", "password": "secret", "host": "10.0.0.50", "port": 443.0,
	"datacenter": "dc1", "driver": "cisco_ndfc", "hostname": "ndfc-controller-01",
}

func newNdcTestServer(t *testing.T, lastBody *string) *httptest.Server {
	t.Helper()
	capture := func(r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-device-controllers", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				jsonResponse(w, http.StatusCreated, ndcItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(ndcItem))
		})
		mux.HandleFunc("/api/v2/network-device-controllers/23", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method != http.MethodGet {
				capture(r)
			}
			jsonResponse(w, http.StatusOK, ndcItem)
		})
		mux.HandleFunc("/api/v2/network-device-controllers/23/credentials", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ndcCredentialsItem)
		})
		mux.HandleFunc("/api/v2/network-device-controllers/23/actions/deploy-confirm", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	})
	return httptest.NewServer(mux)
}

func TestNetworkDeviceControllerListAndAliases(t *testing.T) {
	var body string
	srv := newNdcTestServer(t, &body)
	defer srv.Close()

	for _, name := range []string{"network-device-controller", "nd-controller", "ndc"} {
		out, err := runCLI(t, srv, name, "list")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "ndfc-controller-01") {
			t.Errorf("%s list: expected output to contain the identifier, got: %s", name, out)
		}
	}
}

func TestNetworkDeviceControllerGetCmd(t *testing.T) {
	var body string
	srv := newNdcTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "ndc", "get", "23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"driver":"cisco_ndfc"`) {
		t.Errorf("expected driver in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "ndc", "get"); err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestNetworkDeviceControllerCreateFromFlags(t *testing.T) {
	var body string
	srv := newNdcTestServer(t, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "ndc", "create",
		"--management-address", "10.0.0.50", "--datacenter-name", "dc1",
		"--driver", "cisco_ndfc", "--username", "admin", "--management-password", "secret",
		"--site-id", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"managementAddress":"10.0.0.50"`, `"datacenterName":"dc1"`, `"driver":"cisco_ndfc"`, `"managementPort":443`, `"siteId":1`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body missing %s: %s", want, body)
		}
	}
	if strings.Contains(body, "description") {
		t.Errorf("unset optional flags must not be sent: %s", body)
	}
}

func TestNetworkDeviceControllerCreateFlagValidation(t *testing.T) {
	var body string
	srv := newNdcTestServer(t, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "ndc", "create"); err == nil {
		t.Error("expected error when no source flag given")
	}
	if _, err := runCLI(t, srv, "ndc", "create", "--management-address", "10.0.0.50"); err == nil {
		t.Error("expected error when required companion flags are missing")
	}
	if _, err := runCLI(t, srv, "ndc", "create",
		"--management-address", "10.0.0.50", "--datacenter-name", "dc1",
		"--driver", "bogus", "--username", "admin", "--management-password", "secret"); err == nil {
		t.Error("expected error for invalid --driver")
	}
}

func TestNetworkDeviceControllerCredentialsAndDeployConfirmCmd(t *testing.T) {
	var body string
	srv := newNdcTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "ndc", "get-credentials", "23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"password":"secret"`) {
		t.Errorf("credentials output wrong: %s", out)
	}

	if _, err := runCLI(t, srv, "ndc", "deploy-confirm", "23"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNetworkDeviceControllerHelpListsCommands(t *testing.T) {
	out, err := runCLI(t, nil, "network-device-controller", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "get", "create", "update", "delete", "config-example", "get-credentials", "deploy-confirm"} {
		if !strings.Contains(out, sub) {
			t.Errorf("help missing sub-command %q", sub)
		}
	}
}
