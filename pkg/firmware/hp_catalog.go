package firmware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
)

func all(slice StringOrSlice, condition func(any) bool) bool {
	for _, v := range slice {
		if !condition(v) {
			return false
		}
	}
	return true
}

type hpCatalogTemplate struct {
	Date                 string        `json:"date"`
	Description          string        `json:"description"`
	DeviceClass          string        `json:"deviceclass"`
	MinimumActiveVersion string        `json:"minimum_active_version"`
	RebootRequired       string        `json:"reboot_required"`
	Target               StringOrSlice `json:"target"`
	Version              string        `json:"version"`
}

type StringOrSlice []string

func (sos *StringOrSlice) UnmarshalJSON(data []byte) error {

	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*sos = slice
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*sos = []string{single}
		return nil
	}

	return fmt.Errorf("target can be either a string or an array of strings, but was neither")
}

func parseHpCatalog(configFile rawConfigFile, client metalcloud.MetalCloudClient, filter string, uploadToRepo, downloadBinaries bool, repoConfig repoConfiguration) (firmwareCatalog, []*firmwareBinary, error) {
	catalog, err := generateHpCatalog(configFile)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	binaries, err := parseHpBinaryInventory(configFile, uploadToRepo, downloadBinaries, repoConfig)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	return catalog, binaries, nil
}

func generateHpCatalog(configFile rawConfigFile) (firmwareCatalog, error) {
	catalog := firmwareCatalog{
		Name:                          configFile.Name,
		Description:                   configFile.Description,
		Vendor:                        configFile.Vendor,
		VendorID:                      configFile.Vendor,
		VendorURL:                     configFile.CatalogUrl,
		VendorReleaseTimestamp:        "",
		UpdateType:                    getUpdateType(configFile),
		MetalSoftServerTypesSupported: []string{},
		ServerTypesSupported:          []string{},
		Configuration:                 map[string]any{},
		CreatedTimestamp:              time.Now().Format(time.RFC3339),
	}
	return catalog, nil
}

func parseHpBinaryInventory(configFile rawConfigFile, uploadToRepo, downloadBinaries bool, repoConfig repoConfiguration) ([]*firmwareBinary, error) {
	hpSupportToken := os.Getenv("METALCLOUD_HP_SUPPORT_TOKEN")

	if configFile.CatalogUrl != "" {
		err := downloadHpCatalog(configFile.CatalogUrl, configFile.LocalCatalogPath, hpSupportToken)
		if err != nil {
			return []*firmwareBinary{}, err
		}
	}

	repositoryURL := repoConfig.HttpUrl
	if repositoryURL == "" {
		var err error
		repositoryURL, err = configuration.GetFirmwareRepositoryURL()
		if uploadToRepo && err != nil {
			return nil, fmt.Errorf("Error getting firmware repository URL: %v", err)
		}
	}

	jsonFile, err := os.Open(configFile.LocalCatalogPath)
	if err != nil {
		return []*firmwareBinary{}, err
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var packages map[string]hpCatalogTemplate
	json.Unmarshal(byteValue, &packages)

	binaries := []*firmwareBinary{}
	for key, value := range packages {
		if strings.HasSuffix(key, "fwpkg") &&
			value.DeviceClass != "" && value.DeviceClass != "null" &&
			value.Target != nil && len(value.Target) > 0 &&
			all(value.Target, func(v any) bool { return v != "" && v != "null" }) {

			downloadURL, err := getDownloadURL(configFile.CatalogUrl, key)
			if err != nil {
				//ignore invalid urls
				continue
			}

			componentRepoUrl, err := url.JoinPath(repositoryURL, key)
			if err != nil {
				return nil, err
			}

			localPath := ""
			if configFile.LocalFirmwarePath != "" && downloadBinaries {
				localPath, err = filepath.Abs(filepath.Join(configFile.LocalFirmwarePath, key))

				if err != nil {
					return nil, fmt.Errorf("error getting download binary absolute path: %v", err)
				}
			}

			supportedDevices := []map[string]string{}

			for _, target := range value.Target {
				supportedDevices = append(supportedDevices, map[string]string{
					"DeviceClass":            value.DeviceClass,
					"Target":                 target,
					"MinimumVersionRequired": value.MinimumActiveVersion,
				})
			}

			binaries = append(binaries, &firmwareBinary{
				ExternalId:             key,
				FileName:               key,
				Name:                   key,
				Description:            value.Description,
				PackageId:              key,
				PackageVersion:         value.Version,
				RebootRequired:         value.RebootRequired == "yes",
				UpdateSeverity:         updateSeverityUnknown,
				Hash:                   "",
				HashingAlgorithm:       "",
				SupportedDevices:       supportedDevices,
				SupportedSystems:       []map[string]string{},
				VendorProperties:       map[string]any{},
				VendorReleaseTimestamp: "",
				CreatedTimestamp:       time.Now().Format(time.RFC3339),
				DownloadURL:            downloadURL,
				RepoURL:                componentRepoUrl,
				LocalPath:              localPath,
			})
		}
	}
	return binaries, nil
}

func downloadHpCatalog(catalogURL string, catalogFilePath string, hpSupportToken string) error {
	req, err := http.NewRequest(http.MethodGet, catalogURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	if hpSupportToken != "" {
		encodedData := base64.StdEncoding.EncodeToString([]byte(hpSupportToken + ":null"))
		req.Header.Set("Authorization", "Basic "+encodedData)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error downloading HP firmware inverntory: %v", resp.Status)
	}

	defer resp.Body.Close()

	file, err := os.Create(catalogFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getDownloadURL(catalogURL string, key string) (string, error) {
	u, err := url.Parse(catalogURL)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(path.Dir(path.Dir(u.Path)), key)
	return u.String(), nil
}
