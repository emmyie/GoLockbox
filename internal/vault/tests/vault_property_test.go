package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"GoLockbox/internal/vault"
)

func TestProperty_MetadataRoundTrip(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.CreateVault(dir, "pw123")

	// Add files
	for i := 0; i < 10; i++ {
		src := filepath.Join(dir, "f.txt")
		os.WriteFile(src, []byte("X"), 0600)
		v.AddFile(src, "file.txt")
	}

	// Reopen vault
	v2, err := vault.OpenVault(dir, "pw123")
	if err != nil {
		t.Fatalf("OpenVault failed: %v", err)
	}

	if len(v2.ListFiles()) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(v2.ListFiles()))
	}
}

func TestProperty_RenamePreservesContent(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.CreateVault(dir, "pw123")

	src := filepath.Join(dir, "a.txt")
	os.WriteFile(src, []byte("DATA"), 0600)
	v.AddFile(src, "a.txt")

	v.RenameFile(0, "b.txt")

	out, _ := v.OpenFile(0)
	if string(out) != "DATA" {
		t.Fatalf("content changed after rename")
	}
}
