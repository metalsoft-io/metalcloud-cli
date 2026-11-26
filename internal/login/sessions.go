package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Session is the sign-in kept for one endpoint.
type Session struct {
	Token            string `json:"token"`
	ExpiresTimestamp string `json:"expiresTimestamp"`
	UserEmail        string `json:"userEmail"`
	CreatedTimestamp string `json:"createdTimestamp"`
	// ServiceAccount marks a session signed in from a service account's token file.
	ServiceAccount bool `json:"serviceAccount,omitempty"`
}

// SessionsPath returns the session file's path. It is a variable so tests can point it elsewhere.
var SessionsPath = func() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find the home folder: %w", err)
	}
	return filepath.Join(home, ".metalcloud", "sessions.json"), nil
}

func sessionKey(endpoint string) string {
	return strings.TrimRight(endpoint, "/")
}

// serviceAccountKey is the entry of a service-account session. It is kept apart from the person's
// entry for the same endpoint, and from the entries of other service accounts, by the token file's
// absolute path: each token file is one account.
func serviceAccountKey(endpoint, tokenFile string) string {
	path, err := filepath.Abs(tokenFile)
	if err != nil {
		path = filepath.Clean(tokenFile)
	}
	return sessionKey(endpoint) + "#service-account:" + path
}

func readSessions(path string) (map[string]Session, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]Session{}, nil
	}
	if err != nil {
		return nil, err
	}

	sessions := map[string]Session{}
	if len(strings.TrimSpace(string(data))) == 0 {
		return sessions, nil
	}
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return sessions, nil
}

// selfUID is the current user's id. It is a variable so tests can make their own files look like
// another user's.
var selfUID = os.Getuid()

// otherOwner tells whether a file of this owner can be changed by someone other than you and root,
// who can change any file anyway.
func otherOwner(uid int) bool {
	return uid != selfUID && uid != 0
}

// warnIfOpen warns on stderr when users other than you and root can read or change the session file,
// or replace it through its folder: by owning them, or through their modes. The file is used anyway;
// the next write (login, logout, an ended session) replaces it with one of yours with mode 0600, and
// makes the folder 0700 if it is yours. A missing file holds nothing, so it is not looked at.
func warnIfOpen(path string) {
	perm, uid, ok := unixStat(path)
	if !ok {
		return
	}
	switch {
	case otherOwner(uid):
		fmt.Fprintf(warnOut, "Warning: %s holds session tokens but belongs to another user (uid %d); remove it and sign in again.\n", path, uid)
	case perm&0o077 != 0:
		fmt.Fprintf(warnOut, "Warning: %s holds session tokens but other users can read or change it (mode %04o); run `chmod 600 %s`.\n", path, perm, path)
	}

	dir := filepath.Dir(path)
	dirPerm, dirUID, ok := unixStat(dir)
	switch {
	case !ok:
	case otherOwner(dirUID):
		fmt.Fprintf(warnOut, "Warning: %s belongs to another user (uid %d), who can replace the session file in it; run `sudo chown $USER %s`.\n", dir, dirUID, dir)
	case dirPerm&0o022 != 0:
		fmt.Fprintf(warnOut, "Warning: other users can write to %s (mode %04o) and replace the session file in it; run `chmod 700 %s`.\n", dir, dirPerm, dir)
	}
}

// writeSessions replaces the file through a 0600 temporary file in the same folder, so the result is
// 0600 whatever the old file's mode was.
func writeSessions(path string, sessions map[string]Session) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	// MkdirAll leaves the mode of an existing folder alone, and ~/.metalcloud often exists as 0755.
	if err := os.Chmod(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".sessions-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// LoadSession returns the saved session for the endpoint, or nil when there is none.
func LoadSession(endpoint string) (*Session, error) {
	path, err := SessionsPath()
	if err != nil {
		return nil, err
	}
	sessions, err := readSessions(path)
	if err != nil {
		return nil, err
	}
	warnIfOpen(path)
	s, ok := sessions[sessionKey(endpoint)]
	if !ok {
		return nil, nil
	}
	return &s, nil
}

// SaveSession stores the session for the endpoint and keeps the entries of other endpoints.
func SaveSession(endpoint string, s Session) error {
	path, err := SessionsPath()
	if err != nil {
		return err
	}
	sessions, err := readSessions(path)
	if err != nil {
		return err
	}
	sessions[sessionKey(endpoint)] = s
	return writeSessions(path, sessions)
}

// DeleteSession removes the entry for the endpoint. A missing entry is not an error.
func DeleteSession(endpoint string) error {
	path, err := SessionsPath()
	if err != nil {
		return err
	}
	sessions, err := readSessions(path)
	if err != nil {
		return err
	}
	key := sessionKey(endpoint)
	if _, ok := sessions[key]; !ok {
		return nil
	}
	delete(sessions, key)
	return writeSessions(path, sessions)
}

// LoadServiceAccountSession returns the saved service-account session for the endpoint and token
// file, or nil when there is none. The person's session and the sessions of other token files are not
// looked at.
func LoadServiceAccountSession(endpoint, tokenFile string) (*Session, error) {
	path, err := SessionsPath()
	if err != nil {
		return nil, err
	}
	sessions, err := readSessions(path)
	if err != nil {
		return nil, err
	}
	warnIfOpen(path)
	s, ok := sessions[serviceAccountKey(endpoint, tokenFile)]
	if !ok {
		return nil, nil
	}
	return &s, nil
}

// SaveServiceAccountSession stores the session as the service-account session of the endpoint and
// token file, marked as such, and keeps every other entry, the person's session for the endpoint and
// other token files' sessions included.
func SaveServiceAccountSession(endpoint, tokenFile string, s Session) error {
	path, err := SessionsPath()
	if err != nil {
		return err
	}
	sessions, err := readSessions(path)
	if err != nil {
		return err
	}
	s.ServiceAccount = true
	sessions[serviceAccountKey(endpoint, tokenFile)] = s
	return writeSessions(path, sessions)
}

// DeleteServiceAccountSession removes the service-account session of the endpoint and token file and
// keeps every other entry. A missing entry is not an error.
func DeleteServiceAccountSession(endpoint, tokenFile string) error {
	path, err := SessionsPath()
	if err != nil {
		return err
	}
	sessions, err := readSessions(path)
	if err != nil {
		return err
	}
	key := serviceAccountKey(endpoint, tokenFile)
	if _, ok := sessions[key]; !ok {
		return nil
	}
	delete(sessions, key)
	return writeSessions(path, sessions)
}
