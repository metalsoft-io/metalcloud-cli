package login

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func useTempSessionsPath(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".metalcloud", "sessions.json")
	old := SessionsPath
	SessionsPath = func() (string, error) { return path, nil }
	t.Cleanup(func() { SessionsPath = old })
	return path
}

func TestSaveSessionWritesOwnerOnlyFileAndKeepsOtherEndpoints(t *testing.T) {
	path := useTempSessionsPath(t)

	if err := SaveSession("https://a.example.com/", Session{Token: "ta", UserEmail: "a@x"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveSession("https://b.example.com", Session{Token: "tb"}); err != nil {
		t.Fatal(err)
	}

	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if dirInfo.Mode().Perm() != 0o700 {
		t.Errorf("folder mode = %o, want 700", dirInfo.Mode().Perm())
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Mode().Perm() != 0o600 {
		t.Errorf("file mode = %o, want 600", fileInfo.Mode().Perm())
	}

	// the trailing slash is dropped from the key, in both directions
	a, err := LoadSession("https://a.example.com")
	if err != nil || a == nil || a.Token != "ta" || a.UserEmail != "a@x" {
		t.Fatalf("a = %+v, err = %v", a, err)
	}
	b, err := LoadSession("https://b.example.com/")
	if err != nil || b == nil || b.Token != "tb" {
		t.Fatalf("b = %+v, err = %v", b, err)
	}
	other, err := LoadSession("https://c.example.com")
	if err != nil || other != nil {
		t.Fatalf("other = %+v, err = %v; want nil, nil", other, err)
	}

	if err := DeleteSession("https://a.example.com/"); err != nil {
		t.Fatal(err)
	}
	if a, _ := LoadSession("https://a.example.com"); a != nil {
		t.Errorf("a still present after delete: %+v", a)
	}
	if b, _ := LoadSession("https://b.example.com"); b == nil || b.Token != "tb" {
		t.Errorf("b lost by deleting a: %+v", b)
	}
	// deleting an entry that is not there is not an error
	if err := DeleteSession("https://nope.example.com"); err != nil {
		t.Errorf("delete of missing entry: %v", err)
	}
}

func TestSaveSessionKeeps0600WhenRewriting(t *testing.T) {
	path := useTempSessionsPath(t)

	if err := SaveSession("https://a.example.com", Session{Token: "t1"}); err != nil {
		t.Fatal(err)
	}
	// someone loosened the file in the meantime
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SaveSession("https://a.example.com", Session{Token: "t2"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("file mode after rewrite = %o, want 600", info.Mode().Perm())
	}
	if s, _ := LoadSession("https://a.example.com"); s == nil || s.Token != "t2" {
		t.Errorf("session = %+v, want token t2", s)
	}
}

func TestSaveSessionTightensExistingFolder(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix modes do not apply on Windows")
	}
	dir := filepath.Join(t.TempDir(), ".metalcloud")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "sessions.json")
	old := SessionsPath
	SessionsPath = func() (string, error) { return path, nil }
	t.Cleanup(func() { SessionsPath = old })

	if err := SaveSession("https://a.example.com", Session{Token: "ta"}); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("folder mode = %o, want 700", got)
	}
}

func TestServiceAccountSessionIsKeptApartFromThePersonsSession(t *testing.T) {
	path := useTempSessionsPath(t)
	const endpoint = "https://a.example.com/"
	tokenFile := filepath.Join(t.TempDir(), "token")

	if err := SaveSession(endpoint, Session{Token: "person", UserEmail: "p@x"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveServiceAccountSession(endpoint, tokenFile, Session{Token: "robot", ExpiresTimestamp: "2030-01-01T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}

	person, err := LoadSession("https://a.example.com")
	if err != nil || person == nil || person.Token != "person" || person.ServiceAccount {
		t.Fatalf("person = %+v, err = %v", person, err)
	}
	robot, err := LoadServiceAccountSession("https://a.example.com", tokenFile)
	if err != nil || robot == nil || robot.Token != "robot" || !robot.ServiceAccount {
		t.Fatalf("robot = %+v, err = %v", robot, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	key := "https://a.example.com#service-account:" + tokenFile
	if _, ok := raw[key]; !ok {
		t.Errorf("no %q entry in %s", key, data)
	}
	if _, has := raw["https://a.example.com"]["serviceAccount"]; has {
		t.Errorf("the person's entry carries a serviceAccount field: %v", raw["https://a.example.com"])
	}
	if raw[key]["serviceAccount"] != true {
		t.Errorf("the marked entry is not marked: %v", raw[key])
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && fileInfo.Mode().Perm() != 0o600 {
		t.Errorf("file mode = %o, want 600", fileInfo.Mode().Perm())
	}
}

func TestDeleteServiceAccountSessionKeepsThePersonsSession(t *testing.T) {
	useTempSessionsPath(t)
	const endpoint = "https://a.example.com"
	const tokenFile = "/run/secrets/robot-token"

	if err := SaveSession(endpoint, Session{Token: "person"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveServiceAccountSession(endpoint, tokenFile, Session{Token: "robot"}); err != nil {
		t.Fatal(err)
	}

	if err := DeleteServiceAccountSession(endpoint, tokenFile); err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadServiceAccountSession(endpoint, tokenFile); s != nil {
		t.Errorf("marked session still there: %+v", s)
	}
	if s, _ := LoadSession(endpoint); s == nil || s.Token != "person" {
		t.Errorf("the person's session was touched: %+v", s)
	}

	// deleting a missing entry is fine, and DeleteSession leaves a marked one alone
	if err := DeleteServiceAccountSession(endpoint, tokenFile); err != nil {
		t.Errorf("deleting a missing entry: %v", err)
	}
	if err := SaveServiceAccountSession(endpoint, tokenFile, Session{Token: "robot"}); err != nil {
		t.Fatal(err)
	}
	if err := DeleteSession(endpoint); err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadServiceAccountSession(endpoint, tokenFile); s == nil {
		t.Error("DeleteSession removed the marked session")
	}
}

func TestServiceAccountSessionsAreKeyedByTokenFile(t *testing.T) {
	useTempSessionsPath(t)
	const endpoint = "https://a.example.com"
	dir := t.TempDir()
	fileA := filepath.Join(dir, "a")
	fileB := filepath.Join(dir, "b")

	if err := SaveSession(endpoint, Session{Token: "person"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveServiceAccountSession(endpoint, fileA, Session{Token: "robot-a"}); err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadServiceAccountSession(endpoint, fileB); s != nil {
		t.Fatalf("another token file sees the session of a: %+v", s)
	}
	if err := SaveServiceAccountSession(endpoint, fileB, Session{Token: "robot-b"}); err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadServiceAccountSession(endpoint, fileA); s == nil || s.Token != "robot-a" {
		t.Errorf("signing in with b changed a: %+v", s)
	}
	if s, _ := LoadServiceAccountSession(endpoint, fileB); s == nil || s.Token != "robot-b" {
		t.Errorf("b = %+v", s)
	}

	if err := DeleteServiceAccountSession(endpoint, fileA); err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadServiceAccountSession(endpoint, fileA); s != nil {
		t.Errorf("a is still there: %+v", s)
	}
	if s, _ := LoadServiceAccountSession(endpoint, fileB); s == nil || s.Token != "robot-b" {
		t.Errorf("deleting a removed b: %+v", s)
	}
	if s, _ := LoadSession(endpoint); s == nil || s.Token != "person" {
		t.Errorf("the person's session was touched: %+v", s)
	}
}

func TestServiceAccountSessionKeyUsesTheAbsoluteCleanedPath(t *testing.T) {
	useTempSessionsPath(t)
	const endpoint = "https://a.example.com"
	dir := t.TempDir()
	t.Chdir(dir)

	if err := SaveServiceAccountSession(endpoint, "./sub/../token", Session{Token: "robot"}); err != nil {
		t.Fatal(err)
	}
	// the same file, spelled absolutely, finds the session
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if s, _ := LoadServiceAccountSession(endpoint, filepath.Join(cwd, "token")); s == nil || s.Token != "robot" {
		t.Errorf("absolute spelling: %+v", s)
	}
}
