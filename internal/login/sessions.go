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
