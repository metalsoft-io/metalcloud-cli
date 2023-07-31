package firmware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
)

type hpCatalongTemplate struct {
	Date                 string `json:"date"`
	Description          string `json:"description"`
	DeviceClass          string `json:"deviceclass"`
	MinimumActiveVersion string `json:"minimum_active_version"`
	RebootRequired       string `json:"reboot_required"`
	Target               string `json:"target"`
	Version              string `json:"version"`
}

func parseHpCatalog(configFile rawConfigFile, client metalcloud.MetalCloudClient, filter string, uploadToRepo bool, downloadBinaries bool) (firmwareCatalog, []*firmwareBinary, error) {

	catalog, err := generateHpCatalog(configFile)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	binaries, err := parseBinaryInventory(configFile, uploadToRepo, downloadBinaries)
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
		VendorID:                      configFile.Vendor, // TODO: check if this is correct
		VendorURL:                     configFile.CatalogUrl,
		VendorReleaseTimestamp:        time.Now().Format(time.RFC3339),
		UpdateType:                    getUpdateType(configFile),
		MetalSoftServerTypesSupported: []string{},
		ServerTypesSupported:          []string{},
		Configuration:                 map[string]string{},
		CreatedTimestamp:              time.Now().Format(time.RFC3339),
	}
	return catalog, nil
}

func parseBinaryInventory(configFile rawConfigFile, uploadToRepo bool, downloadBinaries bool) ([]*firmwareBinary, error) {

	hpSupportToken := os.Getenv("HP_SUPPORT_TOKEN")

	if configFile.CatalogUrl != "" {
		err := downloadHpCatalog(configFile.CatalogUrl, configFile.LocalCatalogPath, hpSupportToken)
		if err != nil {
			return []*firmwareBinary{}, err
		}
	}

	repositoryURL, err := configuration.GetFirmwareRepositoryURL()
	if uploadToRepo && err != nil {
		return nil, fmt.Errorf("Error getting firmware repository URL: %v", err)
	}
	jsonFile, err := os.Open(configFile.LocalCatalogPath)
	if err != nil {
		return []*firmwareBinary{}, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var packages map[string]hpCatalongTemplate
	json.Unmarshal(byteValue, &packages)

	binaries := []*firmwareBinary{}
	for key, value := range packages {
		if strings.HasSuffix(key, "fwpkg") &&
			value.DeviceClass != "" && value.DeviceClass != "null" &&
			value.Target != "" && value.Target != "null" {

			downloadURL, err := getDownloadURL(configFile.CatalogUrl, key)
			if err != nil {
				//ignore invalid urls
				continue
			}
			componentRepoUrl := path.Join(repositoryURL, key)
			localPath := ""
			if configFile.LocalFirmwarePath != "" && downloadBinaries {
				localPath, err = filepath.Abs(filepath.Join(configFile.LocalFirmwarePath, key))

				if err != nil {
					return nil, fmt.Errorf("error getting download binary absolute path: %v", err)
				}
			}
			binaries = append(binaries, &firmwareBinary{
				ExternalId:       key,
				FileName:         key,
				Name:             key,
				Description:      value.Description,
				PackageId:        key,
				PackageVersion:   value.Version,
				RebootRequired:   value.RebootRequired == "yes",
				UpdateSeverity:   updateSeverityUnknown,
				Hash:             "",
				HashingAlgorithm: "",
				SupportedDevices: []map[string]string{},
				SupportedSystems: []map[string]string{},
				VendorProperties: map[string]string{
					"DeviceClass":            value.DeviceClass,
					"Target":                 value.Target,
					"MinimumVersionRequired": value.MinimumActiveVersion,
				},
				VendorReleaseTimestamp: time.Now().Format(time.RFC3339),
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
