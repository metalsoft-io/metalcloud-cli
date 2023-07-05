package networking

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"

	netHTTP "net/http"

	"golang.org/x/crypto/ssh"

	kh "golang.org/x/crypto/ssh/knownhosts"
)

func RegexCheckIfUrl(result string) bool {
	m, err := regexp.MatchString(`(?m)(http:|https:).*`, result)
	if err != nil {
		//fmt.Println("your regex is faulty")
		// you should log it or throw an error
		//return err.Error()
		return false
	}
	if m {
		return true
	} else {
		return false
	}
}

func CheckRemoteFileExists(remoteURL string, fileName string) (bool, error) {
	resp, err := netHTTP.Get(remoteURL)

	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	responseBody := string(body)
	return strings.Contains(responseBody, fileName), nil
}

func SerializeSSHKey(key ssh.PublicKey) string {
	return key.Type() + " " + base64.StdEncoding.EncodeToString(key.Marshal())
}

// Add host key if host is not found in known_hosts.
// The return object is the error, if nil then connection proceeds, else connection stops.
func AddHostKey(knownHostsFilePath string, remoteAddress net.Addr, publicKey ssh.PublicKey) error {
	knownHostsFile, err := os.OpenFile(knownHostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Hosts file not found at path %s.", knownHostsFilePath)
	}
	defer knownHostsFile.Close()

	fileBytes, err := os.ReadFile(knownHostsFilePath)

	// We add an empty line if the file doesn't end in one and if it's not empty to begin with.
	if len(fileBytes) > 0 && string(fileBytes[len(fileBytes)-1]) != "\r" && string(fileBytes[len(fileBytes)-1]) != "\n" {
		_, err = knownHostsFile.WriteString("\n")

		if err != nil {
			return err
		}
	}

	knownHosts := kh.Normalize(remoteAddress.String())
	_, err = knownHostsFile.WriteString(kh.Line([]string{knownHosts}, publicKey))

	fmt.Printf("Added key %s to known_hosts file %s.", SerializeSSHKey(publicKey), knownHostsFilePath)
	return err
}
