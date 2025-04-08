package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSymlinksFile(t *testing.T) {
	// create a temporary directory for testing
	tempDir := t.TempDir()

	// create a symlink to a file in the temporary directory
	targetPath := filepath.Join(tempDir, "target.txt")
	symlinkPath := filepath.Join(tempDir, "symlink.txt")
	if err := os.WriteFile(targetPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(targetPath, symlinkPath); err != nil {
		t.Fatal(err)
	}

	// test the ReadSymlinksFile function
	info, err := ReadSymlinksFile(tempDir, "symlink.txt")
	if err != nil {
		t.Fatal(err)
	}
	if info.Name() != "target.txt" {
		t.Errorf("expected name to be 'target.txt', got '%s'", info.Name())
	}
	if info.Size() != 4 {
		t.Errorf("expected size to be 4, got %d", info.Size())
	}
}

func TestReadSymlinksFileRemove(t *testing.T) {
	// create a temporary directory for testing
	tempDir := t.TempDir()

	// create a symlink to a file in the temporary directory
	targetPath := filepath.Join(tempDir, "target.txt")
	symlinkPath := filepath.Join(tempDir, "symlink.txt")
	if err := os.WriteFile(targetPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(targetPath, symlinkPath); err != nil {
		t.Fatal(err)
	}

	err2 := os.Remove(targetPath)
	if err2 != nil {
		t.Fatal(err2)
	}
	// test the ReadSymlinksFile function
	_, err := ReadSymlinksFile(tempDir, "symlink.txt")
	if err != nil {
		if os.IsNotExist(err) {
			t.Logf("target.txt does not exist")
			return
		}
		t.Fatal(err)
	}

}
