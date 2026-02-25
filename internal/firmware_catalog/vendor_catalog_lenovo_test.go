package firmware_catalog

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// mockLenovoCatalog creates a lenovoCatalog with configurable updates for testing
func mockLenovoCatalog(updates []*lenovoSoftwareUpdate) lenovoCatalog {
	return lenovoCatalog{Data: updates}
}

// writeMockLenovoCatalogFile writes a JSON catalog file to disk and returns the path
func writeMockLenovoCatalogFile(t *testing.T, dir string, filename string, catalog lenovoCatalog) string {
	t.Helper()

	data, err := json.Marshal(catalog)
	if err != nil {
		t.Fatalf("Failed to marshal mock catalog: %v", err)
	}

	filePath := filepath.Join(dir, filename)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock catalog file: %v", err)
	}

	return filePath
}

func TestProcessLenovoCatalog_NoVendorSystemsFilter(t *testing.T) {
	ctx := context.Background()

	vc := &VendorCatalog{}

	err := vc.processLenovoCatalog(ctx)

	if err == nil {
		t.Fatal("Expected error when no vendor systems filter is provided, got nil")
	}
	if !strings.Contains(err.Error(), "no vendor systems filter provided") {
		t.Errorf("Expected 'no vendor systems filter' error, got: %v", err)
	}
}

func TestProcessLenovoCatalog_PopulatesFilterExFromBasicFilter(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	catalog := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "lnvo-xcc-01",
			ComponentID: lenovoSoftwareUpdateComponentXcc,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/xcc.bin", Description: "XCC Fix"},
			},
			Version:   "1.0",
			UpdateKey: "xcc-key",
		},
	})

	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V.json", catalog)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilter:    []string{"7D2V"},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	// VendorSystemsFilterEx should have been populated from VendorSystemsFilter
	if len(vc.VendorSystemsFilterEx) != 1 {
		t.Errorf("Expected VendorSystemsFilterEx to have 1 entry, got %d", len(vc.VendorSystemsFilterEx))
	}
	if serial, ok := vc.VendorSystemsFilterEx["7D2V"]; !ok || serial != "" {
		t.Errorf("Expected VendorSystemsFilterEx[\"7D2V\"] = \"\", got %q (exists=%v)", serial, ok)
	}
}

func TestProcessLenovoCatalog_FiltersComponentIDs(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	catalog := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "xcc-update",
			ComponentID: lenovoSoftwareUpdateComponentXcc,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/xcc.bin"},
				{Type: lenovoSoftwareUpdateTypeInstallXML, Description: "XCC Firmware"},
			},
			Version:   "2.0",
			UpdateKey: "xcc-key",
		},
		{
			FixID:       "uefi-update",
			ComponentID: lenovoSoftwareUpdateComponentUefi,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/uefi.bin"},
				{Type: lenovoSoftwareUpdateTypeInstallXML, Description: "UEFI Firmware"},
			},
			Version:   "3.0",
			UpdateKey: "uefi-key",
		},
		{
			FixID:       "lxpm-update",
			ComponentID: lenovoSoftwareUpdateComponentLxpm,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/lxpm.bin"},
				{Type: lenovoSoftwareUpdateTypeInstallXML, Description: "LXPM Firmware"},
			},
			Version:   "4.0",
			UpdateKey: "lxpm-key",
		},
		{
			FixID:       "other-update",
			ComponentID: "SomeOtherComponent",
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/other.bin"},
			},
			Version:   "1.0",
			UpdateKey: "other-key",
		},
	})

	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V.json", catalog)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx:  map[string]string{"7D2V": ""},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	// Should only include XCC, UEFI, LXPM (not "SomeOtherComponent")
	if len(vc.Binaries) != 3 {
		t.Fatalf("Expected 3 binaries (XCC, UEFI, LXPM), got %d", len(vc.Binaries))
	}

	expectedFixIDs := map[string]bool{"xcc-update": false, "uefi-update": false, "lxpm-update": false}
	for _, binary := range vc.Binaries {
		if _, ok := expectedFixIDs[*binary.ExternalId]; ok {
			expectedFixIDs[*binary.ExternalId] = true
		} else {
			t.Errorf("Unexpected binary with ExternalId %s", *binary.ExternalId)
		}
	}
	for fixID, found := range expectedFixIDs {
		if !found {
			t.Errorf("Expected binary with ExternalId %s not found", fixID)
		}
	}
}

func TestProcessLenovoCatalog_SkipsUpdatesWithoutDownloadURL(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	catalog := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "no-fix-url",
			ComponentID: lenovoSoftwareUpdateComponentXcc,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeReadMe, URL: "https://example.com/readme.html"},
			},
			Version: "1.0",
		},
		{
			FixID:       "has-fix-url",
			ComponentID: lenovoSoftwareUpdateComponentUefi,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/uefi.bin"},
			},
			Version:   "2.0",
			UpdateKey: "uefi-key",
		},
	})

	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V.json", catalog)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx:  map[string]string{"7D2V": ""},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	if len(vc.Binaries) != 1 {
		t.Fatalf("Expected 1 binary (skipping the one without Fix URL), got %d", len(vc.Binaries))
	}

	if *vc.Binaries[0].ExternalId != "has-fix-url" {
		t.Errorf("Expected binary ExternalId 'has-fix-url', got %s", *vc.Binaries[0].ExternalId)
	}
}

func TestProcessLenovoCatalog_UsesFixIDAsDescriptionFallback(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	catalog := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "lnvo-xcc-fallback",
			ComponentID: lenovoSoftwareUpdateComponentXcc,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/xcc.bin"},
				// No InstallXML file, so description will be empty
			},
			Version:   "1.0",
			UpdateKey: "xcc-key",
		},
	})

	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V.json", catalog)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx:  map[string]string{"7D2V": ""},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	if len(vc.Binaries) != 1 {
		t.Fatalf("Expected 1 binary, got %d", len(vc.Binaries))
	}

	// When no InstallXML description is available, Name should fall back to FixID
	if vc.Binaries[0].Name != "lnvo-xcc-fallback" {
		t.Errorf("Expected binary Name to be 'lnvo-xcc-fallback' (FixID fallback), got %s", vc.Binaries[0].Name)
	}
}

func TestProcessLenovoCatalog_BinaryFields(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	catalog := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "lnvo-xcc-42",
			ComponentID: lenovoSoftwareUpdateComponentXcc,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://download.lenovo.com/xcc42.bin"},
				{Type: lenovoSoftwareUpdateTypeInstallXML, Description: "XCC Firmware v42"},
				{Type: lenovoSoftwareUpdateTypeReadMe, URL: "https://support.lenovo.com/readme42"},
			},
			RequisitesFixIDs: []string{"prereq-1", "prereq-2"},
			Version:          "42.0",
			UpdateKey:        "xcc-update-key",
		},
	})

	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V_S1234.json", catalog)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx:  map[string]string{"7D2V": "S1234"},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	if len(vc.Binaries) != 1 {
		t.Fatalf("Expected 1 binary, got %d", len(vc.Binaries))
	}

	binary := vc.Binaries[0]

	if *binary.ExternalId != "lnvo-xcc-42" {
		t.Errorf("ExternalId = %s, want lnvo-xcc-42", *binary.ExternalId)
	}
	if binary.Name != "XCC Firmware v42" {
		t.Errorf("Name = %s, want 'XCC Firmware v42'", binary.Name)
	}
	if binary.VendorDownloadUrl != "https://download.lenovo.com/xcc42.bin" {
		t.Errorf("VendorDownloadUrl = %s, want 'https://download.lenovo.com/xcc42.bin'", binary.VendorDownloadUrl)
	}
	if binary.VendorInfoUrl == nil || *binary.VendorInfoUrl != "https://support.lenovo.com/readme42" {
		t.Errorf("VendorInfoUrl = %v, want 'https://support.lenovo.com/readme42'", binary.VendorInfoUrl)
	}
	if *binary.PackageId != "lnvo-xcc-42" {
		t.Errorf("PackageId = %s, want lnvo-xcc-42", *binary.PackageId)
	}
	if *binary.PackageVersion != "42.0" {
		t.Errorf("PackageVersion = %s, want 42.0", *binary.PackageVersion)
	}
	if !binary.RebootRequired {
		t.Error("RebootRequired should be true")
	}
	if binary.UpdateSeverity != sdk.FIRMWAREBINARYUPDATESEVERITY_UNKNOWN {
		t.Errorf("UpdateSeverity = %v, want UNKNOWN", binary.UpdateSeverity)
	}

	// Check vendor configuration (requisites)
	vendorConfig := binary.Vendor
	requisites, ok := vendorConfig["requires"].([]string)
	if !ok {
		t.Fatalf("Vendor['requires'] should be []string, got %T", vendorConfig["requires"])
	}
	if len(requisites) != 2 || requisites[0] != "prereq-1" || requisites[1] != "prereq-2" {
		t.Errorf("Vendor requisites = %v, want [prereq-1 prereq-2]", requisites)
	}

	// Check supported devices
	devices := binary.VendorSupportedDevices
	if len(devices) != 1 || devices[0]["type"] != "xcc-update-key" {
		t.Errorf("VendorSupportedDevices = %v, want [{type: xcc-update-key}]", devices)
	}

	// Check supported systems
	systems := binary.VendorSupportedSystems
	if len(systems) != 1 || systems[0]["machineType"] != "7D2V" || systems[0]["serialNumber"] != "S1234" {
		t.Errorf("VendorSupportedSystems = %v, want [{machineType: 7D2V, serialNumber: S1234}]", systems)
	}
}

func TestProcessLenovoCatalog_MultipleServerTypes(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	catalog1 := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "update-type1",
			ComponentID: lenovoSoftwareUpdateComponentXcc,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/type1.bin"},
			},
			Version:   "1.0",
			UpdateKey: "key1",
		},
	})

	catalog2 := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "update-type2",
			ComponentID: lenovoSoftwareUpdateComponentUefi,
			Files: []lenovoSoftwareUpdateFile{
				{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/type2.bin"},
			},
			Version:   "2.0",
			UpdateKey: "key2",
		},
	})

	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V.json", catalog1)
	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7Z73_SN001.json", catalog2)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx: map[string]string{
			"7D2V": "",
			"7Z73": "SN001",
		},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	if len(vc.Binaries) != 2 {
		t.Fatalf("Expected 2 binaries from two server types, got %d", len(vc.Binaries))
	}
}

func TestReadLenovoCatalog_EmptyMachineType(t *testing.T) {
	vc := &VendorCatalog{}

	_, err := vc.readLenovoCatalog("", "")

	if err == nil {
		t.Fatal("Expected error when machine type is empty, got nil")
	}
	if !strings.Contains(err.Error(), "machine type must be specified") {
		t.Errorf("Expected 'machine type must be specified' error, got: %v", err)
	}
}

func TestReadLenovoCatalog_LocalFileWithoutSerialNumber(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	expected := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "local-fix",
			ComponentID: "XCC",
			Version:     "1.0",
		},
	})

	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V.json", expected)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
	}

	catalog, err := vc.readLenovoCatalog("7D2V", "")

	if err != nil {
		t.Fatalf("readLenovoCatalog() returned error: %v", err)
	}

	if len(catalog.Data) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(catalog.Data))
	}
	if catalog.Data[0].FixID != "local-fix" {
		t.Errorf("FixID = %s, want 'local-fix'", catalog.Data[0].FixID)
	}
}

func TestReadLenovoCatalog_LocalFileWithSerialNumber(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	expected := mockLenovoCatalog([]*lenovoSoftwareUpdate{
		{
			FixID:       "serial-fix",
			ComponentID: "UEFI",
			Version:     "2.0",
		},
	})

	// With serial number, filename should be lenovo_{machineType}_{serialNumber}.json
	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V_S1234.json", expected)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
	}

	catalog, err := vc.readLenovoCatalog("7D2V", "S1234")

	if err != nil {
		t.Fatalf("readLenovoCatalog() returned error: %v", err)
	}

	if len(catalog.Data) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(catalog.Data))
	}
	if catalog.Data[0].FixID != "serial-fix" {
		t.Errorf("FixID = %s, want 'serial-fix'", catalog.Data[0].FixID)
	}
}

func TestReadLenovoCatalog_InvalidJSON(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write invalid JSON
	err = os.WriteFile(filepath.Join(tempDir, "lenovo_7D2V.json"), []byte("not valid json{{{"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
	}

	_, err = vc.readLenovoCatalog("7D2V", "")

	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestReadLenovoCatalog_DownloadsFallback(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// No local file exists, so it should try to download
	mockResponse := `{"Data":[{"FixID":"downloaded-fix","ComponentID":"XCC","Files":[],"Version":"5.0"}]}`

	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify it's calling the Lenovo content service
				if !strings.Contains(req.URL.String(), "SearchDrivers") {
					t.Errorf("Expected request to SearchDrivers, got %s", req.URL.String())
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(mockResponse)),
				}, nil
			},
		},
	}

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
	}

	catalog, err := vc.readLenovoCatalog("7D2V", "")

	if err != nil {
		t.Fatalf("readLenovoCatalog() returned error: %v", err)
	}

	if len(catalog.Data) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(catalog.Data))
	}
	if catalog.Data[0].FixID != "downloaded-fix" {
		t.Errorf("FixID = %s, want 'downloaded-fix'", catalog.Data[0].FixID)
	}
}

func TestDownloadLenovoFirmwareUpdates_Success(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	expectedResponse := `{"Data":[{"FixID":"fw-001","ComponentID":"XCC"}]}`

	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method
				if req.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", req.Method)
				}

				// Verify URL
				if req.URL.String() != lenovoContentServiceUrl+"SearchDrivers" {
					t.Errorf("Expected URL %s, got %s", lenovoContentServiceUrl+"SearchDrivers", req.URL.String())
				}

				// Verify content type
				if req.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got %s", req.Header.Get("Content-Type"))
				}

				// Verify request body contains expected params
				body, _ := io.ReadAll(req.Body)
				var params map[string]interface{}
				json.Unmarshal(body, &params)

				if params["QueryType"] != "SUP" {
					t.Errorf("Expected QueryType 'SUP', got %v", params["QueryType"])
				}
				if params["IsLatest"] != "true" {
					t.Errorf("Expected IsLatest 'true', got %v", params["IsLatest"])
				}

				targetInfos, ok := params["TargetInfos"].([]interface{})
				if !ok || len(targetInfos) != 1 {
					t.Fatalf("Expected TargetInfos array with 1 entry, got %v", params["TargetInfos"])
				}
				target := targetInfos[0].(map[string]interface{})
				if target["MachineType"] != "7D2V" {
					t.Errorf("Expected MachineType '7D2V', got %v", target["MachineType"])
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(expectedResponse)),
				}, nil
			},
		},
	}

	response, err := downloadLenovoFirmwareUpdates("7D2V", "")

	if err != nil {
		t.Fatalf("downloadLenovoFirmwareUpdates() returned error: %v", err)
	}
	if response != expectedResponse {
		t.Errorf("Response = %s, want %s", response, expectedResponse)
	}
}

func TestDownloadLenovoFirmwareUpdates_WithSerialNumber(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				body, _ := io.ReadAll(req.Body)
				var params map[string]interface{}
				json.Unmarshal(body, &params)

				targetInfos := params["TargetInfos"].([]interface{})
				target := targetInfos[0].(map[string]interface{})

				if target["MachineType"] != "7D2V" {
					t.Errorf("Expected MachineType '7D2V', got %v", target["MachineType"])
				}
				if target["SerialNumber"] != "S1234" {
					t.Errorf("Expected SerialNumber 'S1234', got %v", target["SerialNumber"])
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"Data":[]}`)),
				}, nil
			},
		},
	}

	_, err := downloadLenovoFirmwareUpdates("7D2V", "S1234")

	if err != nil {
		t.Fatalf("downloadLenovoFirmwareUpdates() returned error: %v", err)
	}
}

func TestDownloadLenovoFirmwareUpdates_WithoutSerialNumber(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				body, _ := io.ReadAll(req.Body)
				var params map[string]interface{}
				json.Unmarshal(body, &params)

				targetInfos := params["TargetInfos"].([]interface{})
				target := targetInfos[0].(map[string]interface{})

				// SerialNumber should not be present when empty
				if _, ok := target["SerialNumber"]; ok {
					t.Errorf("SerialNumber should not be present when empty, got %v", target["SerialNumber"])
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"Data":[]}`)),
				}, nil
			},
		},
	}

	_, err := downloadLenovoFirmwareUpdates("7D2V", "")

	if err != nil {
		t.Fatalf("downloadLenovoFirmwareUpdates() returned error: %v", err)
	}
}

func TestDownloadLenovoFirmwareUpdates_NetworkError(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, &mockNetworkError{message: "connection refused"}
			},
		},
	}

	_, err := downloadLenovoFirmwareUpdates("7D2V", "")

	if err == nil {
		t.Fatal("Expected network error, got nil")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("Expected 'connection refused' error, got: %v", err)
	}
}

func TestProcessLenovoCatalog_APIResultCodeError(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	resultCode := 106
	message := "Only POST method is supported"
	catalog := lenovoCatalog{
		ResultCode: &resultCode,
		Message:    &message,
		Data:       nil,
	}

	data, err := json.Marshal(catalog)
	if err != nil {
		t.Fatalf("Failed to marshal mock catalog: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "lenovo_7D2V.json"), data, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock catalog file: %v", err)
	}

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx:  map[string]string{"7D2V": ""},
	}

	err = vc.processLenovoCatalog(ctx)

	if err == nil {
		t.Fatal("Expected error when API returns non-zero ResultCode, got nil")
	}
	if !strings.Contains(err.Error(), "lenovo catalog API returned error") {
		t.Errorf("Expected 'lenovo catalog API returned error' message, got: %v", err)
	}
	if !strings.Contains(err.Error(), "106") {
		t.Errorf("Expected error to contain result code 106, got: %v", err)
	}
}

func TestDownloadLenovoFirmwareUpdates_HTTPError(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`{"error": "internal server error"}`)),
				}, nil
			},
		},
	}

	_, err := downloadLenovoFirmwareUpdates("7D2V", "")

	if err == nil {
		t.Fatal("Expected error for HTTP 500 response, got nil")
	}
	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Errorf("Expected error to contain 'HTTP 500', got: %v", err)
	}
}

func TestProcessLenovoCatalog_SuccessResultCodeZero(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	resultCode := 0
	catalog := lenovoCatalog{
		ResultCode: &resultCode,
		Data: []*lenovoSoftwareUpdate{
			{
				FixID:       "xcc-update",
				ComponentID: lenovoSoftwareUpdateComponentXcc,
				Files: []lenovoSoftwareUpdateFile{
					{Type: lenovoSoftwareUpdateTypeFix, URL: "https://example.com/xcc.bin"},
				},
				Version:   "1.0",
				UpdateKey: "xcc-key",
			},
		},
	}

	data, err := json.Marshal(catalog)
	if err != nil {
		t.Fatalf("Failed to marshal mock catalog: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "lenovo_7D2V.json"), data, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock catalog file: %v", err)
	}

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx:  map[string]string{"7D2V": ""},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	if len(vc.Binaries) != 1 {
		t.Fatalf("Expected 1 binary, got %d", len(vc.Binaries))
	}
}

func TestProcessLenovoCatalog_EmptyCatalog(t *testing.T) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "lenovo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	catalog := mockLenovoCatalog([]*lenovoSoftwareUpdate{})
	writeMockLenovoCatalogFile(t, tempDir, "lenovo_7D2V.json", catalog)

	vc := &VendorCatalog{
		VendorLocalCatalogPath: tempDir,
		VendorSystemsFilterEx:  map[string]string{"7D2V": ""},
	}

	err = vc.processLenovoCatalog(ctx)

	if err != nil {
		t.Fatalf("processLenovoCatalog() returned error: %v", err)
	}

	if len(vc.Binaries) != 0 {
		t.Errorf("Expected 0 binaries for empty catalog, got %d", len(vc.Binaries))
	}
}
