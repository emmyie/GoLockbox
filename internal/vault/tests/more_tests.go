package vault_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"GoLockbox/internal/vault"
)

func TestOpenVault_WrongPassword(t *testing.T) {
	dir := t.TempDir()
	if _, err := vault.CreateVault(dir, "correct"); err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	_, err := vault.OpenVault(dir, "incorrect")
	if err == nil {
		t.Fatalf("expected error opening vault with wrong password")
	}
	if !errors.Is(err, vault.ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestSaveFile_PersistsAcrossOpen(t *testing.T) {
	dir := t.TempDir()
	v, err := vault.CreateVault(dir, "pw123")
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	src := filepath.Join(dir, "orig.txt")
	os.WriteFile(src, []byte("ORIG"), 0600)
	if err := v.AddFile(src, "data.txt"); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	if err := v.SaveFile(0, []byte("UPDATED")); err != nil {
		t.Fatalf("SaveFile failed: %v", err)
	}

	// reopen
	v2, err := vault.OpenVault(dir, "pw123")
	if err != nil {
		t.Fatalf("OpenVault failed: %v", err)
	}

	data, err := v2.OpenFile(0)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	if string(data) != "UPDATED" {
		t.Fatalf("expected UPDATED, got %s", string(data))
	}
}
