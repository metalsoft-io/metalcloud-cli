package networking

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"net/http"
	"net/url"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	urlVerifier "github.com/davidmytton/url-verifier"

	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

const (
	HashingAlgorithmMD5  = "md5"
	HashingAlgorithmSHA1 = "sha1"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	responseBody := string(body)
	return strings.Contains(responseBody, fileName), nil
}

// Returns a list of files that do not exist on the remote URL.
func GetMissingRemoteFiles(remoteURL string, fileNames []string) ([]string, error) {
	var missingFiles []string
	resp, err := http.Get(remoteURL)

	if err != nil {
		return missingFiles, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return missingFiles, err
	}

	responseBody := string(body)

	for _, fileName := range fileNames {
		// TODO: how to check for MD5 hash? Add a file to the repo which holds the MD5 hashes?
		if !strings.Contains(responseBody, fileName) {
			missingFiles = append(missingFiles, fileName)
		}
	}

	return missingFiles, nil
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

func DownloadFile(url, path, hash, hashingAlgorithm, user, password string) error {
	ok := fileExists(path)
	if ok && hash != "" {
		localMD5, err := fileHash(path, hashingAlgorithm)

		if err != nil {
			return err
		}

		if localMD5 == hash {
			fmt.Printf("File %s already exists and is the same file as the one from %s. Skipping download.\n", path, url)
			return nil
		} else {
			fmt.Printf("File %s already exists but is not the same file as the one from %s. Downloading new file.\n", path, url)
		}
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	if user != "" && password != "" {
		encodedData := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
		req.Header.Set("Authorization", "Basic "+encodedData)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%d", http.StatusNotFound)
		}

		return fmt.Errorf("received bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded file '%s' from URL '%s' to path '%s'.\n", filepath.Base(url), url, path)
	return nil
}

func HandleKnownHostsFile() (ssh.HostKeyCallback, string, error) {
	knownHostsFilePath, err := configuration.GetKnownHostsPath()
	if err != nil {
		return nil, "", err
	}

	hostKeyCallback, err := kh.New(knownHostsFilePath)

	if err != nil {
		return nil, "", fmt.Errorf("Received following error when parsing the known_hosts file: %s.", err)
	}

	return hostKeyCallback, knownHostsFilePath, nil
}

func CreateSSHClientConfig(skipHostKeyChecking bool, sshUser, userPrivateSSHKeyPath string) (ssh.ClientConfig, error) {
	hostKeyCallback, knownHostsFilePath, err := HandleKnownHostsFile()

	if err != nil {
		return ssh.ClientConfig{}, err
	}

	// Use SSH key authentication from the auth package.
	clientConfig, err := auth.PrivateKey(
		sshUser,
		userPrivateSSHKeyPath,
		ssh.HostKeyCallback(func(hostname string, remoteAddress net.Addr, publicKey ssh.PublicKey) error {
			var keyError *kh.KeyError
			hostsError := hostKeyCallback(hostname, remoteAddress, publicKey)

			// Reference: https://www.godoc.org/golang.org/x/crypto/ssh/knownhosts#KeyError
			//if keyErr.Want is not empty and
			if errors.As(hostsError, &keyError) {
				if len(keyError.Want) > 0 {
					// If host is known then there is key mismatch and the connection is rejected.
					fmt.Printf(`
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The key sent by the remote host is
%s.
Please contact your system administrator.
Add correct host key in %s to get rid of this message.
Host key for %s has changed and you have requested strict checking.
Host key verification failed.
`,
						SerializeSSHKey(publicKey), knownHostsFilePath, hostname,
					)
					return keyError
				} else {
					// If keyErr.Want slice is empty then host is unknown.
					fmt.Printf(`
The authenticity of host '%s' can't be established.
SSH key is %s.
This key is not known by any other names.
It will be added to known_hosts file %s.
Are you sure you want to continue connecting (yes/no)?
`,
						hostname, SerializeSSHKey(publicKey), knownHostsFilePath,
					)

					if !skipHostKeyChecking {
						reader := bufio.NewReader(os.Stdin)
						input, err := reader.ReadString('\n')

						if err != nil {
							return err
						}

						// Remove \r and \n from input
						input = string(bytes.TrimSuffix([]byte(input), []byte("\r\n")))

						if input != "yes" {
							if input == "no" {
								fmt.Println("Aborting connection.")
							} else {
								fmt.Println("Invalid response given. Expecting 'yes' or 'no'. Aborting connection.")
							}

							return keyError
						}
					} else {
						fmt.Printf("Skipped manual check because 'skip-host-key-checking' is set to true.")
					}

					return AddHostKey(knownHostsFilePath, remoteAddress, publicKey)
				}
			}

			fmt.Printf("Public key exists for remote %s. Establishing connection.\n", hostname)
			return nil
		}),
	)

	if err != nil {
		return ssh.ClientConfig{}, fmt.Errorf("Could not create SSH client config. Received error: %s", err)
	}

	return clientConfig, nil
}

func CreateSSHConnection(skipHostKeyChecking bool, firmwareRepositoryURL, firmwareRepositorySSHPort, sshUser, userPrivateSSHKeyPath string) (scp.Client, *ssh.Client, error) {
	clientConfig, err := CreateSSHClientConfig(skipHostKeyChecking, sshUser, userPrivateSSHKeyPath)

	if err != nil {
		return scp.Client{}, &ssh.Client{}, err
	}

	URL, err := url.Parse(firmwareRepositoryURL)
	if err != nil {
		return scp.Client{}, &ssh.Client{}, err
	}

	firmwareRepositoryHostname := URL.Hostname()
	sshRepositoryHostname := firmwareRepositoryHostname + ":" + firmwareRepositorySSHPort

	fmt.Printf("Establishing connection to hostname %s.\n", sshRepositoryHostname)

	// Create a new SCP client.
	scpClient := scp.NewClient(sshRepositoryHostname, &clientConfig)

	// Connect to the remote server.
	sshClient, err := ssh.Dial("tcp", scpClient.Host, scpClient.ClientConfig)
	if err != nil {
		return scp.Client{}, &ssh.Client{}, err
	}

	if scpClient.Session != nil {
		return scpClient, sshClient, nil
	}

	scpClient.Conn = sshClient.Conn
	scpClient.Session, err = sshClient.NewSession()
	if err != nil {
		return scp.Client{}, &ssh.Client{}, nil
	}

	return scpClient, sshClient, nil
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func fileHash(filePath string, hashingAlgorithm string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	var fileHash hash.Hash
	switch hashingAlgorithm {
	case HashingAlgorithmMD5:
		fileHash = md5.New()
	case HashingAlgorithmSHA1:
		fileHash = sha1.New()
	default:
		validHashingAlgorithms := []string{HashingAlgorithmMD5, HashingAlgorithmSHA1}
		return "", fmt.Errorf("invalid hashing algorithm %s. Supported algorithms are %v", hashingAlgorithm, validHashingAlgorithms)
	}

	if _, err := io.Copy(fileHash, file); err != nil {
		return "", err
	}

	hashInBytes := fileHash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}
