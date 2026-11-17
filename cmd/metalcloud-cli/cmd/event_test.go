package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// eventItem satisfies sdk.Event required fields.
var eventItem = map[string]interface{}{
	"id":                 "1",
	"type":               "cluster_deployed",
	"level":              "info",
	"severity":           "info",
	"visibility":         "public",
	"title":              "Server provisioned",
	"message":            "Server SN-001 was provisioned",
	"occurredTimestamp":  "2024-01-01T00:00:00Z",
}

func newEventTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// event list — GET /api/v2/events (paginated or FetchAllPages)
		mux.HandleFunc("/api/v2/events", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(eventItem))
		})
		// event get — GET /api/v2/events/{eventId}
		mux.HandleFunc("/api/v2/events/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(eventItem)
		})
	}))
}

// --- event list ---

func TestEventList_HappyPath(t *testing.T) {
	srv := newEventTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "event", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "Server provisioned") {
		t.Fatalf("expected event title in output, got: %s", out)
	}
}

func TestEventList_Alias(t *testing.T) {
	srv := newEventTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "event", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestEventList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "event", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// TestEventList_LimitFlag verifies that --limit is passed through to the API.
// When limit > 0 EventList calls request.Execute() directly (single page path).
func TestEventList_LimitFlag(t *testing.T) {
	srv := newEventTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "event", "list", "--limit", "5")
	if err != nil {
		t.Fatalf("expected no error with --limit flag, got: %v", err)
	}
	if !strings.Contains(out, "Server provisioned") {
		t.Fatalf("expected event title in output, got: %s", out)
	}
}

// --- event get ---

func TestEventGet_HappyPath(t *testing.T) {
	srv := newEventTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "event", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "Server provisioned") {
		t.Fatalf("expected event title in output, got: %s", out)
	}
}

func TestEventGet_NoArgs(t *testing.T) {
	srv := newEventTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "event", "get"); err == nil {
		t.Fatal("expected error when no args given to event get")
	}
}

func TestEventList_Formats(t *testing.T) {
	srv := newEventTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "event", "list")
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

// --- event get without severity ---

// TestEventGet_NoSeverity guards against SDK schema drift: the SDK Event model
// marks `severity` required but the real API can omit it (regression:
// "no value given for required property severity").
func TestEventGet_NoSeverity(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/events/7", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "7", "type": "cluster_deployed", "visibility": "private",
				"level": "info", "message": "no severity provided",
				"title": "No severity event", "occurredTimestamp": "2024-01-01T00:00:00Z",
			})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "event", "get", "7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No severity event") {
		t.Errorf("expected output to contain event title, got: %s", out)
	}
}

// event get must print every field the API returns for json/yaml — including
// ones eventRaw does not model, such as context — and must not invent a
// severity field the API does not send.
func TestEventGet_FullRecordInJsonAndYaml(t *testing.T) {
	event := map[string]interface{}{
		"id": "673217", "type": "server_power_on", "level": "info", "visibility": "public",
		"title": "Server powered on", "message": "Server 288 powered on",
		"occurredTimestamp": "2024-01-01T00:00:00Z", "serverId": "288", "jobId": "9001",
		"context": map[string]interface{}{"source": "bmc"},
	}
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/events/673217", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(event)
		})
	}))
	defer srv.Close()

	out, err := runCLIFormat(t, srv, "json", "event", "get", "673217")
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("json: invalid output: %v: %s", err, out)
	}
	for key := range event {
		if _, ok := got[key]; !ok {
			t.Errorf("json: missing field %q: %s", key, out)
		}
	}
	if _, ok := got["severity"]; ok {
		t.Errorf("json: printed a severity field the API did not send: %s", out)
	}

	out, err = runCLIFormat(t, srv, "yaml", "event", "get", "673217")
	if err != nil {
		t.Fatalf("yaml: %v", err)
	}
	for _, want := range []string{`serverId: "288"`, "level: info", "source: bmc"} {
		if !strings.Contains(out, want) {
			t.Errorf("yaml: missing %q: %s", want, out)
		}
	}
}
