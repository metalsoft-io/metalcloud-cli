package login

import (
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
