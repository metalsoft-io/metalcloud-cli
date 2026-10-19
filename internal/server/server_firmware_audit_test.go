package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// The audit endpoint answers with an array of per-server audits (MS-10695).
const firmwareAuditResponseJSON = `[
	{
		"serverId": 1,
		"baseline": {"id": 5, "name": "baseline-a"},
		"components": [
			{
				"componentId": 11, "name": "BIOS", "type": "bios",
				"currentVersion": "1.0", "latestVersion": "2.0",
				"status": "upgrade_available", "belowMinimumVersion": true,
				"upgrades": [{"version": "2.0", "binaryId": 3, "name": "bios-2.0", "rebootRequired": true}],
				"downgrades": []
			}
		]
	},
	{
		"serverId": 2,
		"error": {"statusCode": 404, "message": "no baseline assigned"}
	}
]`

func TestServerFirmwareGenerateAuditArrayResponse(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/firmware/actions/generate-upgrade-audit": testutils.RawHandler(http.StatusOK, firmwareAuditResponseJSON),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerFirmwareGenerateAudit(ctx, []byte(`{"serverIds":[1,2]}`)); err != nil {
			t.Fatalf("expected nil error on a 200 array response, got: %v", err)
		}
	})

	if !strings.Contains(out, `"serverId":1`) || !strings.Contains(out, "no baseline assigned") {
		t.Errorf("expected both audits in the output, got: %s", out)
	}
	if !strings.Contains(rec.last().Body, `"serverIds":[1,2]`) {
		t.Errorf("expected the server ids to be sent, got %q", rec.last().Body)
	}
}

func TestFlattenServerFirmwareAudit(t *testing.T) {
	var audits []sdk.ServerFirmwareAudit
	if err := json.Unmarshal([]byte(firmwareAuditResponseJSON), &audits); err != nil {
		t.Fatalf("decode: %v", err)
	}

	rows := flattenServerFirmwareAudit(audits)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}

	if rows[0].Component != "BIOS" || rows[0].Baseline != "baseline-a" || rows[0].BelowMinimum != "Yes" || rows[0].Upgrades != 1 {
		t.Errorf("unexpected component row: %+v", rows[0])
	}
	if rows[1].ServerId != 2 || rows[1].Error != "404: no baseline assigned" {
		t.Errorf("unexpected error row: %+v", rows[1])
	}
}
