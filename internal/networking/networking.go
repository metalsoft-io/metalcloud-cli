package networking

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"net/http"
	urlVerifier "github.com/davidmytton/url-verifier"

	"golang.org/x/crypto/ssh"

	kh "golang.org/x/crypto/ssh/knownhosts"
)

func CheckValidUrl(rawUrl string) bool {
	verifier := urlVerifier.NewVerifier()
	ret, err := verifier.Verify(rawUrl)

	if err != nil {
		return false
	}

	return ret.IsURL
}

func CheckRemoteFileExists(remoteURL, fileName string) (bool, error) {
	resp, err := http.Get(remoteURL)

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
		return fmt.Errorf("hosts file not found at path %s.", knownHostsFilePath)
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

	fmt.Printf("added key %s to known_hosts file %s.", SerializeSSHKey(publicKey), knownHostsFilePath)
	return err
}

func DownloadFile(url, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", filepath, err)
	}
	defer out.Close()

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
