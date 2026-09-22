package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
)

// aiCmdFixture satisfies both required properties of sdk.AIGenerateResponse.
func aiCmdFixture() map[string]interface{} {
	return map[string]interface{}{
		"result": "There are 3 available servers.",
		"steps":  "1. list servers",
	}
}

func aiCmdServer(t *testing.T, lastBody *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/ai/generate", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			*lastBody = string(body)
			jsonResponse(w, http.StatusOK, aiCmdFixture())
		})
	}))
}

// TestAIGenerate_PositionalPrompt checks the prompt may be given as positional
// text and that the body struct is always sent (a missing body would be the
// literal "null", which the API rejects).
func TestAIGenerate_PositionalPrompt(t *testing.T) {
	var body string
	srv := aiCmdServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "ai", "generate", "--datacenter", "dc1", "which", "servers", "are", "available?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var sent map[string]interface{}
	if err := json.Unmarshal([]byte(body), &sent); err != nil {
		t.Fatalf("request body should be a JSON object, got %q", body)
	}
	if sent["datacenter"] != "dc1" || sent["prompt"] != "which servers are available?" {
		t.Errorf("body should join the positional prompt, got %s", body)
	}
	if !strings.Contains(out, "There are 3 available servers.") {
		t.Errorf("output missing the AI result: %s", out)
	}
}

func TestAIGenerate_ConfigSource(t *testing.T) {
	var body string
	srv := aiCmdServer(t, &body)
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "prompt.json")
	if err := os.WriteFile(path, []byte(`{"datacenter":"dc2","prompt":"list infrastructures"}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if _, err := runCLI(t, srv, "ai", "generate", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"datacenter":"dc2"`) || !strings.Contains(body, `"prompt":"list infrastructures"`) {
		t.Errorf("config document should supply both fields, got %s", body)
	}
}

func TestAIGenerate_RequiresPrompt(t *testing.T) {
	var body string
	srv := aiCmdServer(t, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "ai", "generate", "--datacenter", "dc1"); err == nil {
		t.Fatal("generate without a prompt should fail")
	}
}

func TestAIGenerate_RequiresDatacenter(t *testing.T) {
	var body string
	srv := aiCmdServer(t, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "ai", "generate", "hello"); err == nil {
		t.Fatal("generate without --datacenter should fail")
	}
}

func TestAICommand_Wiring(t *testing.T) {
	found, _, err := aiCmd.Find([]string{"generate"})
	if err != nil || found == aiCmd {
		t.Fatal("ai generate is not registered")
	}
	if found.Annotations[system.REQUIRED_PERMISSION] != system.PERMISSION_AI_READ {
		t.Errorf("ai generate must require %s, got %q", system.PERMISSION_AI_READ, found.Annotations[system.REQUIRED_PERMISSION])
	}
	if !found.SilenceUsage {
		t.Error("ai generate must set SilenceUsage")
	}
}
