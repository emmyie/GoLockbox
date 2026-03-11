package vault_test

import (
	"os"
	"testing"

	"GoLockbox/internal/vault"
)

func FuzzOpenVault(f *testing.F) {
	// Seed with a valid vault
	dir, _ := os.MkdirTemp("", "vault-fuzz-*")
	v, _ := vault.CreateVault(dir, "pw123")
	_ = v

	// Read master.key to use as seed
	mk, _ := os.ReadFile(dir + "/master.key")
	f.Add(mk)

	f.Fuzz(func(t *testing.T, data []byte) {
		// Write fuzzed master.key
		_ = os.WriteFile(dir+"/master.key", data, 0600)

		// Try opening vault — should never panic
		_, err := vault.OpenVault(dir, "pw123")
		if err != nil && !os.IsNotExist(err) {
			// Expected: decryption may fail, but should not panic
			t.Logf("OpenVault error (expected): %v", err)
		}
	})
}

func FuzzDecryptMetadata(f *testing.F) {
	dir, _ := os.MkdirTemp("", "vault-meta-fuzz-*")
	_, err := vault.CreateVault(dir, "pw123")
	if err != nil {
		f.Fatalf("CreateVault failed: %v", err)
	}

	// Seed with real metadata
	meta, _ := os.ReadFile(dir + "/entries.enc")
	f.Add(meta)

	f.Fuzz(func(t *testing.T, data []byte) {
		// Replace metadata with fuzzed data
		_ = os.WriteFile(dir+"/entries.enc", data, 0600)

		// Should never panic
		_, err = vault.OpenVault(dir, "pw123")
		if err != nil && !os.IsNotExist(err) {
			// Expected: decryption may fail, but should not panic
			t.Logf("OpenVault error (expected): %v", err)
		}
	})
}
