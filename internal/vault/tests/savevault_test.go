package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"GoLockbox/internal/vault"
)

// Test that SaveVault removes entries whose files are deleted from disk.
func TestSaveVault_RemovesMissingFiles(t *testing.T) {
	dir := t.TempDir()
	v, err := vault.CreateVault(dir, "pw123")
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	// add two files
	src1 := filepath.Join(dir, "a.txt")
	os.WriteFile(src1, []byte("A"), 0600)
	src2 := filepath.Join(dir, "b.txt")
	os.WriteFile(src2, []byte("B"), 0600)

	if err := v.AddFile(src1, "a.txt"); err != nil {
		t.Fatalf("AddFile a failed: %v", err)
	}
	if err := v.AddFile(src2, "b.txt"); err != nil {
		t.Fatalf("AddFile b failed: %v", err)
	}

	files := v.ListFiles()
	if len(files) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(files))
	}

	// remove the first file from disk
	os.Remove(filepath.Join(dir, "files", files[0].Name))

	// call SaveVault to reconcile
	if err := v.SaveVault(); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	files = v.ListFiles()
	if len(files) != 1 {
		t.Fatalf("expected 1 entry after SaveVault, got %d", len(files))
	}

	if files[0].Name != filepath.Base(files[0].Name) {
		t.Fatalf("unexpected remaining entry name: %s", files[0].Name)
	}

	// reopen and verify metadata persisted
	v2, err := vault.OpenVault(dir, "pw123")
	if err != nil {
		t.Fatalf("OpenVault failed: %v", err)
	}
	if len(v2.ListFiles()) != 1 {
		t.Fatalf("expected 1 entry after reopen, got %d", len(v2.ListFiles()))
	}
}

// Test that SaveVault preserves multiple colliding filenames (AddFile handles collision suffixes).
func TestSaveVault_PreservesCollisions(t *testing.T) {
	dir := t.TempDir()
	v, err := vault.CreateVault(dir, "pw123")
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	// create same-named source file and add twice
	src := filepath.Join(dir, "dup.txt")
	os.WriteFile(src, []byte("X"), 0600)
	if err := v.AddFile(src, "dup.txt"); err != nil {
		t.Fatalf("AddFile first failed: %v", err)
	}
	// write again and add second time
	os.WriteFile(src, []byte("Y"), 0600)
	if err := v.AddFile(src, "dup.txt"); err != nil {
		t.Fatalf("AddFile second failed: %v", err)
	}

	files := v.ListFiles()
	if len(files) != 2 {
		t.Fatalf("expected 2 entries for collisions, got %d", len(files))
	}

	// Ensure both underlying files exist
	for _, e := range files {
		if _, err := os.Stat(filepath.Join(dir, "files", e.Name)); err != nil {
			t.Fatalf("expected file for entry %s to exist: %v", e.Name, err)
		}
	}

	// SaveVault should be a no-op but must not drop entries
	if err := v.SaveVault(); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	if len(v.ListFiles()) != 2 {
		t.Fatalf("expected 2 entries after SaveVault, got %d", len(v.ListFiles()))
	}
}
