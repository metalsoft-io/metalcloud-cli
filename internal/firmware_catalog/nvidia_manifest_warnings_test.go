package firmware_catalog

import (
	"strings"
	"testing"
)

func TestNvidiaManifestWarnsAboutStepsVersionedByTheBundle(t *testing.T) {
	manifest := &NvidiaBundleManifest{
		BundleVersion: "25.06.1",
		Steps: []NvidiaBundleStep{
			{Id: "gpu", ComponentClass: "gpu", PackageVersion: "25.06.1"},
			{Id: "nvswitch", ComponentClass: "nvswitch", PackageVersion: "25.06.1", VerifyVersion: "96.10.6D.00.01"},
			{Id: "retimer", ComponentClass: "pcie_retimer", PackageVersion: "2.7.24"},
		},
	}

	warnings := manifest.Warnings()

	if len(warnings) != 1 {
		t.Fatalf("expected one warning, got %v", warnings)
	}
	if !strings.HasPrefix(warnings[0], `Step "gpu": packageVersion is the bundle version 25.06.1 and there is no verifyVersion`) {
		t.Fatalf("unexpected warning: %s", warnings[0])
	}
}
