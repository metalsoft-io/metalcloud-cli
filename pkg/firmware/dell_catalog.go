package firmware

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"unicode/utf16"

	"golang.org/x/net/html/charset"
)

const (
	STOP_AFTER           int    = 3
)

type Manifest struct {
	XMLName    xml.Name            `xml:"Manifest"`
	Software   []SoftwareBundle    `xml:"SoftwareBundle"`
	Components []SoftwareComponent `xml:"SoftwareComponent"`
}

type SoftwareBundle struct {
	XMLName          xml.Name         `xml:"SoftwareBundle"`
	DateTime         string           `xml:"dateTime,attr"`
	Path             string           `xml:"path,attr"`
	RebootRequired   string           `xml:"rebootRequired,attr"`
	Name             Name             `xml:"Name"`
	SupportedDevices SupportedDevices `xml:"SupportedDevices"`
	SupportedSystems SupportedSystems `xml:"SupportedSystems"`
}

type Name struct {
	Display string `xml:"Display"`
}

type SupportedDevices struct {
	Devices []Device `xml:"Device"`
}

type Device struct {
	ComponentID string `xml:"componentID,attr"`
	Display     string `xml:"Display"`
}

type SupportedSystems struct {
	Brands []Brand `xml:"Brand"`
}

type Brand struct {
	Models []Model `xml:"Model"`
}

type Model struct {
	SystemID     string `xml:"systemID,attr"`
	SystemIDType string `xml:"systemIDType,attr"`
	Display      string `xml:"Display"`
}

type SoftwareComponent struct {
	XMLName          xml.Name         `xml:"SoftwareComponent"`
	DateTime         string           `xml:"dateTime,attr"`
	Path             string           `xml:"path,attr"`
	RebootRequired   string           `xml:"rebootRequired,attr"`
	Name             Name             `xml:"Name"`
	SupportedDevices SupportedDevices `xml:"SupportedDevices"`
	SupportedSystems SupportedSystems `xml:"SupportedSystems"`
}

func downloadCatalog(catalogURL string, catalogFilePath string) error {
	req, err := http.NewRequest(http.MethodGet, catalogURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	file, err := os.Create(catalogFilePath)
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

// Is this used?
func utf16ToUTF8(data []byte) []byte {
	d := make([]uint16, (len(data)/2)+1)
	for i := 0; i < len(data); i += 2 {
		d[i/2] = uint16(data[i]) + (uint16(data[i+1]) << 8)
	}
	return []byte(string(utf16.Decode(d)))
}

func BypassReader(label string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func parseDellCatalog(configFile rawConfigFile) {
	if configFile.DownloadCatalog {
		err := downloadCatalog(configFile.CatalogUrl, configFile.CatalogPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	xmlFile, err := os.Open(configFile.CatalogPath)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	reader, err := charset.NewReader(xmlFile, "utf-16")
	if err != nil {
		return
	}
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = BypassReader

	var manifest Manifest

	err = decoder.Decode(&manifest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\n====>Software Bundles:")
	for idx, component := range manifest.Software {
		fmt.Println("DateTime:", component.DateTime)
		fmt.Println("Path:", component.Path)
		fmt.Println("RebootRequired:", component.RebootRequired)
		fmt.Println("Name:", component.Name.Display)
		fmt.Println("SupportedDevices:")
		for _, device := range component.SupportedDevices.Devices {
			fmt.Println("  ComponentID:", device.ComponentID)
			fmt.Println("  Display:", device.Display)
		}
		fmt.Println("SupportedSystems:")
		for _, brand := range component.SupportedSystems.Brands {
			for _, model := range brand.Models {
				fmt.Println("  SystemID:", model.SystemID)
				fmt.Println("  SystemIDType:", model.SystemIDType)
				fmt.Println("  Display:", model.Display)
			}
		}
		fmt.Println("------------------------")

		if idx == STOP_AFTER {
			break
		}
	}

	fmt.Println("\n\n====> Software Components:")
	for idx, component := range manifest.Components {
		fmt.Println("DateTime:", component.DateTime)
		fmt.Println("Path:", component.Path)
		fmt.Println("RebootRequired:", component.RebootRequired)
		fmt.Println("Name:", component.Name.Display)
		fmt.Println("SupportedDevices:")
		for _, device := range component.SupportedDevices.Devices {
			fmt.Println("  ComponentID:", device.ComponentID)
			fmt.Println("  Display:", device.Display)
		}
		fmt.Println("SupportedSystems:")
		for _, brand := range component.SupportedSystems.Brands {
			for _, model := range brand.Models {
				fmt.Println("  SystemID:", model.SystemID)
				fmt.Println("  SystemIDType:", model.SystemIDType)
				fmt.Println("  Display:", model.Display)
			}
		}
		fmt.Println("------------------------")

		if idx == STOP_AFTER {
			break
		}
	}
}
