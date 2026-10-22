package server_cleanup_policy

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bodyRecorder serves a cleanup policy and records the JSON body of each request.
type bodyRecorder struct {
	mu     sync.Mutex
	bodies []map[string]any
}

func (b *bodyRecorder) handler(method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
			return
		}

		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)

		b.mu.Lock()
		b.bodies = append(b.bodies, body)
		b.mu.Unlock()

		testutils.RawHandler(http.StatusOK, cleanupPolicyItem)(w, r)
	}
}

func TestCleanupPolicyCreate_SendsClearTpmAndResetRaidControllers(t *testing.T) {
	rec := &bodyRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/cleanup-policies": rec.handler(http.MethodPost),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CleanupPolicyCreate(ctx, "wipe", false, false, false, true, true, "raid1", "raid1", "raid5", "raid5", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if len(rec.bodies) != 1 {
		t.Fatalf("expected 1 request, got %d", len(rec.bodies))
	}
	if rec.bodies[0]["clearTpm"] != float64(1) || rec.bodies[0]["resetRaidControllers"] != float64(1) {
		t.Errorf("expected clearTpm and resetRaidControllers to be 1, got %v and %v",
			rec.bodies[0]["clearTpm"], rec.bodies[0]["resetRaidControllers"])
	}
}

func TestCleanupPolicyCreate_DefaultsToNoTpmClearOrControllerReset(t *testing.T) {
	rec := &bodyRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/cleanup-policies": rec.handler(http.MethodPost),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := CleanupPolicyCreate(ctx, "wipe", false, false, false, false, false, "raid1", "raid1", "raid5", "raid5", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if rec.bodies[0]["clearTpm"] != float64(0) || rec.bodies[0]["resetRaidControllers"] != float64(0) {
		t.Errorf("expected both to be 0, got %v and %v", rec.bodies[0]["clearTpm"], rec.bodies[0]["resetRaidControllers"])
	}
}

// updateCommand builds a command with the flags CleanupPolicyUpdate inspects.
func updateCommand(args ...string) *cobra.Command {
	cmd := &cobra.Command{Use: "update"}
	for _, name := range []string{"clear-tpm", "reset-raid-controllers", "label"} {
		if name == "label" {
			cmd.Flags().String(name, "", "")
		} else {
			cmd.Flags().Bool(name, false, "")
		}
	}
	_ = cmd.Flags().Parse(args)
	return cmd
}

func TestCleanupPolicyUpdate_SetsOnlyTheFlagsGiven(t *testing.T) {
	rec := &bodyRecorder{}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/cleanup-policies/1": rec.handler(http.MethodPatch),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := CleanupPolicyUpdate(ctx, "1", "", false, false, false, true, false, "", "", "", "", "", updateCommand("--clear-tpm")); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got := rec.bodies[0]["clearTpm"]; got != float64(1) {
		t.Errorf("clearTpm: expected 1, got %v", got)
	}
	if _, sent := rec.bodies[0]["resetRaidControllers"]; sent {
		t.Error("resetRaidControllers must not be sent when its flag was not given")
	}

	if err := CleanupPolicyUpdate(ctx, "1", "", false, false, false, false, true, "", "", "", "", "", updateCommand("--reset-raid-controllers", "--clear-tpm=false")); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if rec.bodies[1]["resetRaidControllers"] != float64(1) || rec.bodies[1]["clearTpm"] != float64(0) {
		t.Errorf("expected resetRaidControllers=1 and clearTpm=0 (explicitly turned off), got %v and %v",
			rec.bodies[1]["resetRaidControllers"], rec.bodies[1]["clearTpm"])
	}
}

func TestCleanupPolicyGet_ShowsClearTpmAndResetRaidControllers(t *testing.T) {
	const item = `{"id":1,"label":"wipe","cleanupDrivesForOobEnabledServer":0,"clearTpm":1,"recreateRaid":0,
		"resetRaidControllers":0,"disableEmbeddedNics":0,"raidOneDrive":"raid1","raidTwoDrives":"raid1",
		"raidEvenNumberMoreThanTwoDrives":"raid5","raidOddNumberMoreThanOneDrive":"raid5",
		"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","skipRaidActions":[]}`

	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/cleanup-policies/1": testutils.RawHandler(http.StatusOK, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	viper.Set(formatter.ConfigFormat, "csv")
	defer viper.Set(formatter.ConfigFormat, "json")

	out := testutils.CaptureStdout(t, func() {
		if err := CleanupPolicyGet(ctx, "1"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	records, err := csv.NewReader(strings.NewReader(out)).ReadAll()
	if err != nil || len(records) < 2 {
		t.Fatalf("unexpected CSV output %q (%v)", out, err)
	}
	row := map[string]string{}
	for i, name := range records[0] {
		row[name] = records[1][i]
	}

	if row["Clear TPM"] != "true" {
		t.Errorf("Clear TPM: expected true, got %q", row["Clear TPM"])
	}
	if row["Reset RAID Controllers"] != "false" {
		t.Errorf("Reset RAID Controllers: expected false, got %q", row["Reset RAID Controllers"])
	}
}
