// Package testutils provides shared test helpers for metalcloud-cli unit tests.
package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

// SetupTestFormat configures viper so formatter.PrintResult uses JSON output.
// Call this from TestMain in each test package.
func SetupTestFormat() {
	viper.Set(formatter.ConfigFormat, "json")
}

// SetupTestContext creates a context with an SDK API client pointed at serverURL.
func SetupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

// NewTestServer creates an httptest.Server that routes requests by URL path.
func NewTestServer(routes map[string]http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	return httptest.NewServer(mux)
}

// JSONHandler returns an http.HandlerFunc that writes statusCode and body as JSON.
func JSONHandler(statusCode int, body any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			panic(fmt.Sprintf("testutils.JSONHandler: encode failed: %v", err))
		}
	}
}

// RawHandler returns an http.HandlerFunc that writes a raw JSON string.
func RawHandler(statusCode int, rawJSON string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(rawJSON))
	}
}

// PaginatedMeta builds the standard pagination metadata envelope.
func PaginatedMeta(currentPage, totalPages, itemsPerPage int32) map[string]any {
	return map[string]any{
		"currentPage":  currentPage,
		"totalPages":   totalPages,
		"itemsPerPage": itemsPerPage,
	}
}

// PaginatedResponse wraps a data slice in the standard paginated envelope.
func PaginatedResponse(data any, currentPage, totalPages int32) map[string]any {
	return map[string]any{
		"data": data,
		"meta": PaginatedMeta(currentPage, totalPages, 100),
	}
}

// MultiPageServer creates an httptest.Server that serves sequential pages of items.
// Each call to the path returns the next page until exhausted, then repeats the last.
func MultiPageServer(path string, pages []any) *httptest.Server {
	var counter atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		idx := int(counter.Add(1)) - 1
		if idx >= len(pages) {
			idx = len(pages) - 1
		}
		total := int32(len(pages))
		current := int32(idx + 1)
		resp := PaginatedResponse(pages[idx], current, total)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
	return httptest.NewServer(mux)
}

// ErrorHandler returns an http.HandlerFunc that returns an API error response.
func ErrorHandler(statusCode int, message string) http.HandlerFunc {
	return RawHandler(statusCode, fmt.Sprintf(`{"message":%q,"statusCode":%d}`, message, statusCode))
}

// CaptureStdout runs fn and returns everything it wrote to os.Stdout. Use it to
// assert on formatter output, which is printed rather than returned.
func CaptureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	_ = w.Close()
	os.Stdout = old
	return <-done
}

// NoContentHandler returns an http.HandlerFunc that replies with an empty
// 204 response. Use it for action endpoints that return no body; a JSON
// encoder would panic writing a body on a 204.
func NoContentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// NetworkDeviceFixture returns a JSON-shaped network device that satisfies
// every required property of sdk.NetworkDevice, so the strict SDK unmarshaller
// accepts it. Callers may override or add keys on the returned map.
func NetworkDeviceFixture(id string, siteId int, identifierString string) map[string]any {
	return map[string]any{
		"id": id, "revision": 1, "status": "active", "vendorId": 1, "siteId": siteId,
		"identifierString": identifierString, "applyIdentifierAsHostnameOnNextDeploy": false,
		"description": "", "chassisIdentifier": "", "country": "", "city": "",
		"datacenterMeta": "", "datacenterRoom": "", "datacenterRack": "",
		"rackPositionUpperUnit": 1, "rackPositionLowerUnit": 1,
		"managementAddress": "10.0.0.100", "managementAddressPrefixLength": 24,
		"managementAddressGateway": "10.0.0.1", "managementPort": 22,
		"syslogEnabled": 0, "snmpServiceEnabled": false, "snmpMonitoringEnabled": false,
		"username": "admin", "managementMacAddress": "00:11:22:33:44:55", "serialNumber": "SN-1",
		"driver": "sonic_enterprise", "position": "leaf", "backupEnabled": false,
		"driftDetectionEnabled": false, "driftDetectionSyncStatus": "not_supported",
		"orderIndex": 0, "tags": []any{}, "tagsMap": map[string]any{},
		"readyForInitialConfiguration": 0, "bootstrapReadinessCheckInProgress": 0,
		"subnetOobId": 0, "subnetOobIndex": 0, "requiresOsInstall": false,
		"bootstrapExpectedPartnerHostname": "", "loopbackAddressIpv6": "", "asn": 0,
		"vtepAddressIpv6": "", "mlagSystemMac": "", "mlagDomainId": 0, "quarantineVlan": 0,
		"variablesMaterializedForOSAssets": map[string]any{}, "secretsMaterializedForOSAssets": map[string]any{},
		"bootstrapReadinessCheckResult": map[string]any{}, "isGateway": false, "portCount": 0,
	}
}
