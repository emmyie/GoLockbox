package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"GoLockbox/internal/crypto"
)

type Vault struct {
	path      string
	masterKey []byte
	entries   []FileEntry
}

// CreateVault creates a new vault at base with a random master key encrypted by the password.
func CreateVault(base, password string) (*Vault, error) {
	expanded := expandPath(base)
	if err := os.MkdirAll(expanded, 0700); err != nil {
		return nil, err
	}

	// Generate master key
	masterKey, err := crypto.GenerateRandomKey(crypto.KeySize)
	if err != nil {
		return nil, err
	}

	// Derive key from password
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, err
	}

	derived, err := crypto.DeriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	defer zeroBytes(derived)

	encMaster, err := crypto.EncryptBytes(masterKey, derived)
	if err != nil {
		return nil, err
	}

	// master.key format: [1 byte version][16 bytes salt][ciphertext...]
	buf := make([]byte, 1+len(salt)+len(encMaster))
	buf[0] = masterKeyVersion

	copy(buf[1:], salt)
	copy(buf[1+len(salt):], encMaster)

	if err := os.WriteFile(masterKeyPath(expanded), buf, 0600); err != nil {
		return nil, err
	}

	v := &Vault{
		path:      expanded,
		masterKey: masterKey,
		entries:   []FileEntry{},
	}

	if err := v.saveMetadata(); err != nil {
		return nil, err
	}

	return v, nil
}

// OpenVault opens an existing vault at base using the password.
func OpenVault(base, password string) (*Vault, error) {
	expanded := expandPath(base)

	data, err := os.ReadFile(masterKeyPath(expanded))
	if err != nil {
		return nil, err
	}

	if len(data) < 1+crypto.SaltSize {
		return nil, errors.New("master key file too short")
	}

	version := data[0]
	if version != masterKeyVersion {
		return nil, fmt.Errorf("unsupported master key version: %d", version)
	}

	salt := data[1 : 1+crypto.SaltSize]
	ciphertext := data[1+crypto.SaltSize:]
	derived, err := crypto.DeriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	defer zeroBytes(derived)

	masterKey, err := crypto.DecryptBytes(ciphertext, derived)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	if len(masterKey) != crypto.KeySize {
		return nil, errors.New("invalid master key size")
	}

	v := &Vault{
		path:      expanded,
		masterKey: masterKey,
		entries:   []FileEntry{},
	}

	if err := v.loadMetadata(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Vault) filePath(name string) string {
	return filepath.Join(v.path, "files", name)
}

func (v *Vault) saveMetadata() error {
	meta := metadata{
		Version: metadataVersion,
		Entries: v.entries,
	}

	raw, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	enc, err := crypto.EncryptBytes(raw, v.masterKey)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(v.path, 0700); err != nil {
		return err
	}

	return os.WriteFile(metadataPath(v.path), enc, 0600)
}

func (v *Vault) loadMetadata() error {
	path := metadataPath(v.path)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			v.entries = []FileEntry{}
			return nil
		}
		return err
	}

	raw, err := crypto.DecryptBytes(data, v.masterKey)
	if err != nil {
		return err
	}

	var meta metadata
	if err := json.Unmarshal(raw, &meta); err != nil {
		return err
	}
	if meta.Version != metadataVersion {
		return fmt.Errorf("unsupported metadata version: %d", meta.Version)
	}

	v.entries = meta.Entries

	return nil
}

func (v *Vault) ListFiles() []FileEntry {
	out := make([]FileEntry, len(v.entries))
	copy(out, v.entries)

	return out
}

// SaveVault ensures metadata and directory structure are consistent.
func (v *Vault) SaveVault() error {
	filesDir := filepath.Join(v.path, "files")
	if err := os.MkdirAll(filesDir, 0700); err != nil {
		return err
	}

	var kept []FileEntry
	for _, e := range v.entries {
		p := v.filePath(e.Name)
		if _, err := os.Stat(p); err == nil {
			kept = append(kept, e)
		}
	}

	v.entries = kept

	return v.saveMetadata()
}

func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// Close saves metadata and zeros the master key from memory.
func (v *Vault) Close() error {
	if v == nil {
		return nil
	}

	if err := v.SaveVault(); err != nil {
		return err
	}

	if v.masterKey != nil {
		zeroBytes(v.masterKey)
		v.masterKey = nil
	}

	return nil
}

// Path returns the vault filesystem path.
func (v *Vault) Path() string {
	if v == nil {
		return ""
	}
	return v.path
}
