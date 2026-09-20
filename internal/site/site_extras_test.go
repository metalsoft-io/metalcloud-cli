package site

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func siteStatisticsFixture() map[string]interface{} {
	return map[string]interface{}{
		"siteControllerSeenAliveStatus": map[string]interface{}{
			"may_be_offline":                   1,
			"offline":                          2,
			"was_seen_connected_very_recently": 5,
		},
		"sitesTotalCount":  8,
		"sitesActiveCount": 6,
		// The live API returns per-site resource counters here, not the
		// site-controller status objects the typed SDK model declares. This is
		// why SiteStatistics parses the body raw.
		"sitesResourceCount": []interface{}{
			map[string]interface{}{
				"id":                   1,
				"serversCount":         4,
				"infrastructuresCount": 36,
				"storagesCount":        0,
				"networksCount":        0,
			},
		},
	}
}

func TestSiteStatistics(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/sites/statistics": testutils.JSONHandler(http.StatusOK, siteStatisticsFixture()),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := SiteStatistics(ctx); err != nil {
			t.Errorf("SiteStatistics: expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/sites/statistics": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := SiteStatistics(ctx); err == nil {
			t.Error("expected error for HTTP 500, got nil")
		}
	})
}

func TestSiteRegistryUrls(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/sites/controllers/actions/get/registry-urls": testutils.JSONHandler(http.StatusOK,
				[]interface{}{"registry.metalsoft.dev", "registry.metalsoft.io"}),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := SiteRegistryUrls(ctx); err != nil {
			t.Errorf("SiteRegistryUrls: expected nil error, got: %v", err)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/sites/controllers/actions/get/registry-urls": testutils.JSONHandler(http.StatusOK, []interface{}{}),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := SiteRegistryUrls(ctx); err != nil {
			t.Errorf("SiteRegistryUrls empty: expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/sites/controllers/actions/get/registry-urls": testutils.ErrorHandler(http.StatusForbidden, "forbidden"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := SiteRegistryUrls(ctx); err == nil {
			t.Error("expected error for HTTP 403, got nil")
		}
	})
}
