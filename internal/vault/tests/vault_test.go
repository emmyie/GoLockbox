package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"GoLockbox/internal/vault"
)

// helper: create a temp vault directory
func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "vault-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return dir
}

// helper: write a plaintext file
func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestCreateAndOpenVault(t *testing.T) {
	dir := tempDir(t)
	password := "test-password"

	// create vault
	_, err := vault.CreateVault(dir, password)
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	// ensure master.key exists
	if _, err := os.Stat(filepath.Join(dir, "master.key")); err != nil {
		t.Fatalf("master.key missing: %v", err)
	}

	// reopen vault
	v2, err := vault.OpenVault(dir, password)
	if err != nil {
		t.Fatalf("OpenVault failed: %v", err)
	}

	if len(v2.ListFiles()) != 0 {
		t.Fatalf("expected empty vault, got %d entries", len(v2.ListFiles()))
	}
}

func TestAddAndOpenFile(t *testing.T) {
	dir := tempDir(t)
	password := "pw123"

	v, err := vault.CreateVault(dir, password)
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	// create plaintext file
	src := writeTempFile(t, dir, "hello.txt", "Hello World")

	// add file to vault
	if err := v.AddFile(src, "hello.txt"); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// verify metadata
	files := v.ListFiles()
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	// open file (decrypt)
	data, err := v.OpenFile(0)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	if string(data) != "Hello World" {
		t.Fatalf("unexpected file contents: %s", string(data))
	}
}

func TestRenameFile(t *testing.T) {
	dir := tempDir(t)
	password := "pw123"

	v, _ := vault.CreateVault(dir, password)
	src := writeTempFile(t, dir, "a.txt", "AAA")

	if err := v.AddFile(src, "a.txt"); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	if err := v.RenameFile(0, "b.txt"); err != nil {
		t.Fatalf("RenameFile failed: %v", err)
	}

	files := v.ListFiles()
	if files[0].Name != "b.txt" {
		t.Fatalf("expected renamed file to be b.txt, got %s", files[0].Name)
	}

	// ensure file still decrypts correctly
	data, err := v.OpenFile(0)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	if string(data) != "AAA" {
		t.Fatalf("unexpected file contents after rename: %s", string(data))
	}
}

func TestDeleteFile(t *testing.T) {
	dir := tempDir(t)
	password := "pw123"

	v, _ := vault.CreateVault(dir, password)
	src := writeTempFile(t, dir, "x.txt", "XXX")

	if err := v.AddFile(src, "x.txt"); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	if err := v.DeleteFile(0); err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}

	if len(v.ListFiles()) != 0 {
		t.Fatalf("expected 0 files after delete")
	}
}

func TestExportFile(t *testing.T) {
	dir := tempDir(t)
	password := "pw123"

	v, _ := vault.CreateVault(dir, password)
	src := writeTempFile(t, dir, "exp.txt", "EXPORT ME")

	if err := v.AddFile(src, "exp.txt"); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	exportDir := tempDir(t)
	dst := filepath.Join(exportDir, "out.txt")

	if err := v.ExportFile(0, dst); err != nil {
		t.Fatalf("ExportFile failed: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read exported file: %v", err)
	}

	if string(data) != "EXPORT ME" {
		t.Fatalf("unexpected exported contents: %s", string(data))
	}
}

func TestReopenVaultWithFiles(t *testing.T) {
	dir := tempDir(t)
	password := "pw123"

	// create vault and add file
	v, _ := vault.CreateVault(dir, password)
	src := writeTempFile(t, dir, "z.txt", "ZZZ")
	v.AddFile(src, "z.txt")

	// reopen vault
	v2, err := vault.OpenVault(dir, password)
	if err != nil {
		t.Fatalf("OpenVault failed: %v", err)
	}

	files := v2.ListFiles()
	if len(files) != 1 {
		t.Fatalf("expected 1 file after reopen, got %d", len(files))
	}

	data, err := v2.OpenFile(0)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	if string(data) != "ZZZ" {
		t.Fatalf("unexpected file contents after reopen: %s", string(data))
	}
}
