package server_registration_profile

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

const registrationProfileItem = `{
	"id": 1, "name": "profile-1", "revision": "1",
	"isDefault": false,
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"settings": {
		"registerCredentials": "default",
		"minimumNumberOfConnectedInterfaces": 1,
		"alwaysDiscoverInterfacesWithBDK": false,
		"enableTpm": false,
		"enableIntelTxt": false,
		"enableSyslogMonitoring": false,
		"disableTpmAfterRegistration": false,
		"defaultVirtualMediaProtocol": "ipmi",
		"resetRaidControllers": false,
		"cleanupDrives": false,
		"recreateRaid": false,
		"disableEmbeddedNics": false,
		"raidOneDrive": "",
		"raidTwoDrives": "",
		"raidEvenNumberMoreThanTwoDrives": "",
		"raidOddNumberMoreThanOneDrive": ""
	},
	"links": []
}`

func TestRegistrationProfileList(t *testing.T) {
	listPage1 := fmt.Sprintf(`{"data":[%s],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`, registrationProfileItem)

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles": testutils.RawHandler(http.StatusOK, listPage1),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileList(ctx); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileList(ctx); err == nil {
			t.Error("expected error for HTTP 500, got nil")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileList(ctx); err != nil {
			t.Errorf("expected nil error for empty list, got: %v", err)
		}
	})

	t.Run("Pagination3Pages", func(t *testing.T) {
		makeItems := func(start, count int) string {
			s := `[`
			for i := 0; i < count; i++ {
				if i > 0 {
					s += ","
				}
				id := start + i
				s += fmt.Sprintf(`{"id":%d,"name":"profile-%d","revision":"1","isDefault":false,"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","settings":{"registerCredentials":"default","minimumNumberOfConnectedInterfaces":1,"alwaysDiscoverInterfacesWithBDK":false,"enableTpm":false,"enableIntelTxt":false,"enableSyslogMonitoring":false,"disableTpmAfterRegistration":false,"defaultVirtualMediaProtocol":"ipmi","resetRaidControllers":false,"cleanupDrives":false,"recreateRaid":false,"disableEmbeddedNics":false,"raidOneDrive":"","raidTwoDrives":"","raidEvenNumberMoreThanTwoDrives":"","raidOddNumberMoreThanOneDrive":""},"links":[]}`, id, id)
			}
			s += `]`
			return s
		}

		var call atomic.Int32
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles": func(w http.ResponseWriter, r *http.Request) {
				n := int(call.Add(1))
				var page, total int
				var items string
				switch n {
				case 1:
					page, total = 1, 3
					items = makeItems(1, 100)
				case 2:
					page, total = 2, 3
					items = makeItems(101, 100)
				default:
					page, total = 3, 3
					items = makeItems(201, 5)
				}
				body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`, items, page, total)
				testutils.RawHandler(http.StatusOK, body)(w, r)
			},
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileList(ctx); err != nil {
			t.Errorf("expected nil error for pagination, got: %v", err)
		}
		if got := int(call.Load()); got < 3 {
			t.Errorf("expected at least 3 page fetches, got %d", got)
		}
	})
}

func TestRegistrationProfileGet(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles/1": testutils.RawHandler(http.StatusOK, registrationProfileItem),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileGet(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileGet(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})
}

func TestRegistrationProfileSearch(t *testing.T) {
	t.Run("WithSiteId", func(t *testing.T) {
		var query string
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles/search": func(w http.ResponseWriter, r *http.Request) {
				query = r.URL.RawQuery
				testutils.RawHandler(http.StatusOK, registrationProfileItem)(w, r)
			},
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)

		out := testutils.CaptureStdout(t, func() {
			if err := RegistrationProfileSearch(ctx, 7); err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
		})

		if !strings.Contains(out, "profile-1") {
			t.Errorf("expected the profile in the output, got: %s", out)
		}
		if !strings.Contains(query, "siteId=7") {
			t.Errorf("expected the site filter to be sent, got query: %s", query)
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles/search": testutils.ErrorHandler(http.StatusInternalServerError, "boom"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileSearch(ctx, 1); err == nil {
			t.Fatal("expected an error for HTTP 500")
		}
	})
}

func TestRegistrationProfileForServer(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles/search/for-server/12": testutils.RawHandler(http.StatusOK, registrationProfileItem),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)

		out := testutils.CaptureStdout(t, func() {
			if err := RegistrationProfileForServer(ctx, "12"); err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
		})
		if !strings.Contains(out, "profile-1") {
			t.Errorf("expected the profile in the output, got: %s", out)
		}
	})

	t.Run("InvalidServerId", func(t *testing.T) {
		ctx := testutils.SetupTestContext("http://127.0.0.1:1")
		if err := RegistrationProfileForServer(ctx, "abc"); err == nil {
			t.Fatal("expected an error for an invalid server ID")
		}
	})
}

func TestRegistrationProfileSystemDefaults(t *testing.T) {
	const settingsJSON = `{
		"registerCredentials": "user",
		"minimumNumberOfConnectedInterfaces": 1,
		"enableTpm": true,
		"defaultVirtualMediaProtocol": "HTTPS",
		"raidOneDrive": "RAID0",
		"dpuMode": "dpu"
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles/system-defaults": testutils.RawHandler(http.StatusOK, settingsJSON),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)

		out := testutils.CaptureStdout(t, func() {
			if err := RegistrationProfileSystemDefaults(ctx); err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
		})
		if !strings.Contains(out, "HTTPS") {
			t.Errorf("expected the default settings in the output, got: %s", out)
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/servers/registration-profiles/system-defaults": testutils.ErrorHandler(http.StatusInternalServerError, "boom"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := RegistrationProfileSystemDefaults(ctx); err == nil {
			t.Fatal("expected an error for HTTP 500")
		}
	})
}
