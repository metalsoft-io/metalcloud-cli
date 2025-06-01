package firmware_catalog

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

type VendorCatalog struct {
	CatalogInfo             sdk.FirmwareCatalog
	Binaries                []*sdk.FirmwareBinary
	ServerTypesFilter       []string
	VendorSystemsFilter     []string
	VendorSystemsFilterEx   map[string]string
	VendorToken             string
	VendorLocalCatalogPath  string
	VendorLocalBinariesPath string
	DownloadBinaries        bool
	UploadBinaries          bool
	RepoBaseUrl             string
	RepoSshHost             string
	RepoSshUser             string
	UserPrivateKeyPath      string
	KnownHostsPath          string
	IgnoreHostKeyCheck      bool
}

func NewVendorCatalogFromCreateOptions(options FirmwareCatalogCreateOptions) (*VendorCatalog, error) {
	vendor := sdk.ServerFirmwareCatalogVendor(options.Vendor)
	if !vendor.IsValid() {
		return nil, fmt.Errorf("invalid vendor flag %s - valid options are: %v", options.Vendor, sdk.AllowedServerFirmwareCatalogVendorEnumValues)
	}

	updateType := sdk.CatalogUpdateType(options.UpdateType)
	if !updateType.IsValid() {
		return nil, fmt.Errorf("invalid update type flag %s - valid options are: %v", options.UpdateType, sdk.AllowedCatalogUpdateTypeEnumValues)
	}

	return &VendorCatalog{
		CatalogInfo: sdk.FirmwareCatalog{
			Name:        options.Name,
			Description: sdk.PtrString(options.Description),
			Vendor:      options.Vendor,
			UpdateType:  options.UpdateType,
			VendorUrl:   sdk.PtrString(options.VendorUrl),
		},
		Binaries:                []*sdk.FirmwareBinary{},
		ServerTypesFilter:       options.ServerTypesFilter,
		VendorSystemsFilter:     options.VendorSystemsFilter,
		VendorToken:             options.VendorToken,
		VendorLocalCatalogPath:  options.VendorLocalCatalogPath,
		VendorLocalBinariesPath: options.VendorLocalBinariesPath,
		DownloadBinaries:        options.DownloadBinaries,
		UploadBinaries:          options.UploadBinaries,
		RepoBaseUrl:             options.RepoBaseUrl,
		RepoSshHost:             options.RepoSshHost,
		RepoSshUser:             options.RepoSshUser,
		UserPrivateKeyPath:      options.UserPrivateKeyPath,
		KnownHostsPath:          options.KnownHostsPath,
		IgnoreHostKeyCheck:      options.IgnoreHostKeyCheck,
	}, nil
}

func (vc *VendorCatalog) ProcessVendorCatalog(ctx context.Context) error {
	if len(vc.ServerTypesFilter) > 0 {
		// Lookup the system models for the requested server types
		// and add them to the vendor systems filter
		systemModels, systemModelsEx, err := vc.getFilteredSystemModels(ctx)
		if err != nil {
			return err
		}

		if len(vc.VendorSystemsFilter) > 0 {
			vc.VendorSystemsFilter = append(vc.VendorSystemsFilter, systemModels...)
		} else {
			vc.VendorSystemsFilter = systemModels
		}
		vc.VendorSystemsFilterEx = systemModelsEx
	}

	switch sdk.ServerFirmwareCatalogVendor(vc.CatalogInfo.Vendor) {
	case sdk.SERVERFIRMWARECATALOGVENDOR_DELL:
		return vc.processDellCatalog(ctx)
	case sdk.SERVERFIRMWARECATALOGVENDOR_HP:
		return vc.processHpeCatalog(ctx)
	case sdk.SERVERFIRMWARECATALOGVENDOR_LENOVO:
		return vc.processLenovoCatalog(ctx)
	default:
		return fmt.Errorf("unsupported vendor %s", vc.CatalogInfo.Vendor)
	}
}

func (vc *VendorCatalog) CreateMetalsoftCatalog(ctx context.Context) error {
	firmwareCatalogCreate := sdk.CreateFirmwareCatalog{
		Name:                          vc.CatalogInfo.Name,
		Description:                   vc.CatalogInfo.Description,
		Vendor:                        sdk.ServerFirmwareCatalogVendor(vc.CatalogInfo.Vendor),
		UpdateType:                    sdk.CatalogUpdateType(vc.CatalogInfo.UpdateType),
		VendorUrl:                     vc.CatalogInfo.VendorUrl,
		MetalsoftServerTypesSupported: vc.CatalogInfo.MetalsoftServerTypesSupported,
		VendorServerTypesSupported:    vc.CatalogInfo.VendorServerTypesSupported,
		VendorId:                      vc.CatalogInfo.VendorId,
		VendorReleaseTimestamp:        vc.CatalogInfo.VendorReleaseTimestamp,
	}

	client := api.GetApiClient(ctx)

	firmwareCatalog, httpRes, err := client.FirmwareCatalogAPI.
		CreateFirmwareCatalogs(ctx).
		CreateFirmwareCatalog(firmwareCatalogCreate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	vc.CatalogInfo.Id = firmwareCatalog.Id

	var sftpClient *sftp.Client
	if vc.UploadBinaries {
		authMethod, err := privateKeyFile(vc.UserPrivateKeyPath)
		if err != nil {
			return fmt.Errorf("failed to parse private key '%s': %v", vc.UserPrivateKeyPath, err)
		}

		var hostKeyCallback ssh.HostKeyCallback
		if vc.IgnoreHostKeyCheck {
			hostKeyCallback = ssh.InsecureIgnoreHostKey()
		} else {
			hostKeyCallback, err = knownhosts.New(vc.KnownHostsPath)
			if err != nil {
				return fmt.Errorf("could not create host key callback: %v", err)
			}
		}

		sshConfig := &ssh.ClientConfig{
			User: vc.RepoSshUser,
			Auth: []ssh.AuthMethod{
				authMethod,
			},
			HostKeyCallback: hostKeyCallback,
		}

		sshClient, err := ssh.Dial("tcp", vc.RepoSshHost, sshConfig)
		if err != nil {
			return err
		}
		defer sshClient.Close()

		sftpClient, err := sftp.NewClient(sshClient)
		if err != nil {
			return err
		}
		defer sftpClient.Close()
	}

	var repoUrl *url.URL
	if vc.RepoBaseUrl != "" {
		repoUrl, err = url.Parse(vc.RepoBaseUrl)
		if err != nil {
			return fmt.Errorf("unable to parse repo catalog URL: %v", err)
		}
	}

	for _, binary := range vc.Binaries {
		// Determine the local path to the binary
		var localPath string
		if vc.DownloadBinaries {
			// Download the binary from the vendor catalog
			localPath, err = vc.downloadBinary(binary)
			if err != nil {
				return err
			}
		} else {
			if vc.VendorLocalBinariesPath != "" {
				// Use the local path provided
				localPath, err = filepath.Abs(filepath.Join(vc.VendorLocalBinariesPath, *binary.ExternalId))
				if err != nil {
					return fmt.Errorf("error getting download binary absolute path: %v", err)
				}
			}
		}

		// Set the binary URL in the repo
		if repoUrl != nil {
			binary.CacheDownloadUrl = sdk.PtrString(repoUrl.JoinPath(*binary.ExternalId).String())
		} else {
			binary.CacheDownloadUrl = sdk.PtrString(*binary.ExternalId)
		}

		// Upload the binary to the repository if needed
		if vc.UploadBinaries {
			err = vc.uploadBinaryToRepository(sftpClient, localPath, *binary.ExternalId)
			if err != nil {
				return fmt.Errorf("error uploading binary to repository: %v", err)
			}
		}

		binaryCreate := sdk.CreateFirmwareBinary{
			CatalogId:              firmwareCatalog.Id,
			ExternalId:             binary.ExternalId,
			VendorInfoUrl:          binary.VendorInfoUrl,
			VendorDownloadUrl:      binary.VendorDownloadUrl,
			CacheDownloadUrl:       binary.CacheDownloadUrl,
			Name:                   binary.Name,
			PackageId:              binary.PackageId,
			PackageVersion:         binary.PackageVersion,
			RebootRequired:         binary.RebootRequired,
			UpdateSeverity:         binary.UpdateSeverity,
			VendorSupportedDevices: binary.VendorSupportedDevices,
			VendorSupportedSystems: binary.VendorSupportedSystems,
			VendorReleaseTimestamp: binary.VendorReleaseTimestamp,
			Vendor:                 binary.Vendor,
		}

		newBinary, httpRes, err := client.FirmwareBinaryAPI.CreateFirmwareBinary(ctx).
			CreateFirmwareBinary(binaryCreate).
			Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		binary.Id = newBinary.Id
	}

	return nil
}

// Returns a list of vendor models corresponding to the requested list of server types.
func (vc *VendorCatalog) getFilteredSystemModels(ctx context.Context) ([]string, map[string]string, error) {
	systemModels := []string{}
	systemModelsEx := map[string]string{}

	client := api.GetApiClient(ctx)

	serverTypes, httpRes, err := client.ServerTypeAPI.GetServerTypes(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, nil, err
	}

	for _, serverType := range serverTypes.Data {
		serverTypeIdentifier := serverType.Name
		if serverType.Label != nil {
			serverTypeIdentifier = *serverType.Label
		}

		if slices.Contains(vc.ServerTypesFilter, serverTypeIdentifier) {
			servers, httpRes, err := client.ServerAPI.GetServers(ctx).FilterServerTypeId([]string{fmt.Sprintf("%d", int(serverType.Id))}).Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return nil, nil, err
			}

			for _, server := range servers.Data {
				if server.Vendor == nil || server.Model == nil {
					continue
				}
				if strings.ToLower(*server.Vendor) != vc.CatalogInfo.Vendor {
					continue
				}
				if slices.Contains(systemModels, *server.Model) {
					continue
				}

				systemModels = append(systemModels, *server.Model)
				systemModelsEx[*server.Model] = *server.VendorSkuId // TODO: What if we have multiple serial numbers for the same model?
			}
		}
	}

	return systemModels, systemModelsEx, nil
}

// Downloads a binary from the vendor catalog
func (vc *VendorCatalog) downloadBinary(binary *sdk.FirmwareBinary) (string, error) {
	if binary.VendorDownloadUrl == "" {
		return "", fmt.Errorf("no vendor download URL provided for binary %s", *binary.ExternalId)
	}

	var err error
	localPath := ""
	if vc.VendorLocalBinariesPath != "" {
		localPath, err = filepath.Abs(filepath.Join(vc.VendorLocalBinariesPath, *binary.ExternalId))
		if err != nil {
			return "", fmt.Errorf("error getting download binary absolute path: %v", err)
		}
	} else {
		// Create a temporary file to download the binary
		tempFile, err := os.CreateTemp("", "binary_*.bin")
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		localPath = tempFile.Name()
	}

	// Download binary from vendor
	resp, err := http.Get(binary.VendorDownloadUrl)
	if err != nil {
		return "", fmt.Errorf("error downloading binary from vendor URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response when downloading binary: %d", resp.StatusCode)
	}

	// Create the file
	outFile, err := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open output file: %v", err)
	}
	defer outFile.Close()

	// Copy the body to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save binary to file: %v", err)
	}

	resp.Body.Close()
	outFile.Close()

	return localPath, nil
}

func (vc *VendorCatalog) uploadBinaryToRepository(sftpClient *sftp.Client, binaryPath string, remotePath string) error {
	srcFile, err := os.Open(binaryPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	err = sftpClient.MkdirAll(path.Dir(remotePath))
	if err != nil {
		return fmt.Errorf("failed to create remote directory: %v", err)
	}

	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// Reads a PEM private key file
func privateKeyFile(file string) (ssh.AuthMethod, error) {
	key, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

func downloadGzipCatalog(url string, filePath string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error downloading catalog: %v", resp.Status)
	}
	defer resp.Body.Close()

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, gzReader)
	if err != nil {
		return err
	}

	return nil
}

func downloadCatalog(url string, filePath string, authToken string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	if authToken != "" {
		encodedData := base64.StdEncoding.EncodeToString([]byte(authToken + ":null"))
		req.Header.Set("Authorization", "Basic "+encodedData)
	}

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error downloading catalog: %v", resp.Status)
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
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

func readXmlDocument(filePath string, target any) error {
	// Read file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read catalog file: %v", err)
	}

	// Try to convert from UTF-16 to UTF-8 if needed
	if len(fileContent) >= 2 && ((fileContent[0] == 0xFF && fileContent[1] == 0xFE) || (fileContent[0] == 0xFE && fileContent[1] == 0xFF)) {
		// Detect endianness
		endian := unicode.LittleEndian
		if fileContent[0] == 0xFE && fileContent[1] == 0xFF {
			endian = unicode.BigEndian
		}

		// Convert to UTF-8
		decoder := unicode.UTF16(endian, unicode.UseBOM)
		utf8Content, err := decoder.NewDecoder().Bytes(fileContent)
		if err == nil {
			fileContent = utf8Content
		} else {
			fmt.Printf("Warning: Failed to convert UTF-16 to UTF-8: %v\n", err)
		}
	}

	// Try to detect BOM and remove it if present
	fileContent = removeBOM(fileContent)

	// Try to clean and normalize XML before parsing
	fileContent = cleanXML(fileContent)

	// Create a reader from the sanitized content
	reader := bytes.NewReader(fileContent)

	// Create XML decoder with charset support
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = AdaptiveReader
	decoder.Strict = false // Be more lenient with XML parsing

	// Decode the XML into target struct
	err = decoder.Decode(target)
	if err != nil {
		return fmt.Errorf("failed to decode XML: %v", err)
	}

	return nil
}

// removeBOM removes the Byte Order Mark from the beginning of the file if present
func removeBOM(data []byte) []byte {
	// Check for UTF-8 BOM (EF BB BF)
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	// Check for UTF-16 LE BOM (FF FE)
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		return data[2:]
	}
	// Check for UTF-16 BE BOM (FE FF)
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		return data[2:]
	}
	return data
}

// cleanXML tries to fix common XML issues
func cleanXML(data []byte) []byte {
	// Convert to string for easier manipulation
	content := string(data)

	// Replace any character that's not valid in XML
	// Using double quotes instead of backticks for the regex pattern
	re := regexp.MustCompile("[^\x09\x0A\x0D\x20-\uD7FF\uE000-\uFFFD\U00010000-\U0010FFFF]")
	content = re.ReplaceAllString(content, "")

	// Ensure XML declaration is correct
	if strings.HasPrefix(content, "<?xml version=\"1.0\" encoding=\"utf-16\"?>") {
		content = strings.Replace(content, "<?xml version=\"1.0\" encoding=\"utf-16\"?>",
			"<?xml version=\"1.0\" encoding=\"utf-8\"?>", 1)
	}

	// Remove incomplete tags
	if strings.Contains(content, "]]></Display>") {
		re := regexp.MustCompile(`Model[^<]*]]></Display>`)
		content = re.ReplaceAllString(content, "")
	}

	return []byte(content)
}

func AdaptiveReader(charset string, input io.Reader) (io.Reader, error) {
	switch {
	case strings.EqualFold(charset, "ISO-8859-1"):
		return charmap.ISO8859_1.NewDecoder().Reader(input), nil
	case strings.EqualFold(charset, "windows-1252"):
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	case strings.EqualFold(charset, "utf-16"):
		// Try to handle UTF-16 encoded content
		data, err := io.ReadAll(input)
		if err != nil {
			return nil, err
		}

		// Use the correct UTF-16 decoder
		decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		decoded, err := decoder.NewDecoder().Bytes(data)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(decoded), nil
	default:
		// For UTF-8 and other encodings, return as is
		return input, nil
	}
}
