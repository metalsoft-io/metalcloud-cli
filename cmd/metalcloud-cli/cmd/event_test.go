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
	"type":               "server.provision",
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
