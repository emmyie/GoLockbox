package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"GoLockbox/internal/crypto"
)

func sanitize(name string) string {
	name = strings.TrimSpace(name)
	name = filepath.Base(name)

	return strings.ReplaceAll(name, string(os.PathSeparator), "-")
}

func resolveCollision(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	for i := 1; ; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s-%d%s", name, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

// AddFile encrypts src and stores it in the vault under the given name.
func (v *Vault) AddFile(src, name string) error {
	name = sanitize(name)
	if name == "" {
		return fmt.Errorf("invalid name")
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	enc, err := crypto.EncryptBytes(data, v.masterKey)
	if err != nil {
		return err
	}

	filesDir := filepath.Join(v.path, "files")
	if err := os.MkdirAll(filesDir, 0700); err != nil {
		return err
	}

	dst := resolveCollision(v.filePath(name))
	if err := os.WriteFile(dst, enc, 0600); err != nil {
		return err
	}

	v.entries = append(v.entries, FileEntry{Name: filepath.Base(dst)})

	return v.saveMetadata()
}

// OpenFile decrypts and returns the contents of the file at index.
func (v *Vault) OpenFile(index int) ([]byte, error) {
	if index < 0 || index >= len(v.entries) {
		return nil, fmt.Errorf("invalid index")
	}

	entry := v.entries[index]
	path := v.filePath(entry.Name)
	enc, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return crypto.DecryptBytes(enc, v.masterKey)
}

// SaveFile encrypts and overwrites the file at index with data.
func (v *Vault) SaveFile(index int, data []byte) error {
	if index < 0 || index >= len(v.entries) {
		return fmt.Errorf("invalid index")
	}

	entry := v.entries[index]
	path := v.filePath(entry.Name)
	enc, err := crypto.EncryptBytes(data, v.masterKey)
	if err != nil {
		return err
	}

	return os.WriteFile(path, enc, 0600)
}

// ExportFile decrypts the file at index and writes it to dst.
// If dst is a directory, the file name is appended.
func (v *Vault) ExportFile(index int, dst string) error {
	if index < 0 || index >= len(v.entries) {
		return fmt.Errorf("invalid index")
	}

	data, err := v.OpenFile(index)
	if err != nil {
		return err
	}

	entry := v.entries[index]
	info, err := os.Stat(dst)
	if err == nil && info.IsDir() {
		dst = filepath.Join(dst, entry.Name)
	}

	return os.WriteFile(dst, data, 0600)
}

// RenameFile renames the file at index to newName (with collision handling).
func (v *Vault) RenameFile(index int, newName string) error {
	if index < 0 || index >= len(v.entries) {
		return fmt.Errorf("invalid index")
	}

	newName = sanitize(newName)
	if newName == "" {
		return fmt.Errorf("invalid name")
	}

	oldEntry := v.entries[index]
	oldPath := v.filePath(oldEntry.Name)
	newPath := resolveCollision(v.filePath(newName))
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	v.entries[index].Name = filepath.Base(newPath)

	return v.saveMetadata()
}

// DeleteFile removes the file at index from disk and metadata.
func (v *Vault) DeleteFile(index int) error {
	if index < 0 || index >= len(v.entries) {
		return fmt.Errorf("invalid index")
	}

	entry := v.entries[index]
	path := v.filePath(entry.Name)

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	v.entries = append(v.entries[:index], v.entries[index+1:]...)

	return v.saveMetadata()
}

// CreateFile creates a new file in the vault with the provided data.
func (v *Vault) CreateFile(name string, data []byte) error {
	name = sanitize(name)
	if name == "" {
		return fmt.Errorf("invalid name")
	}

	enc, err := crypto.EncryptBytes(data, v.masterKey)
	if err != nil {
		return err
	}

	filesDir := filepath.Join(v.path, "files")
	if err := os.MkdirAll(filesDir, 0700); err != nil {
		return err
	}

	dst := resolveCollision(v.filePath(name))
	if err := os.WriteFile(dst, enc, 0600); err != nil {
		return err
	}

	v.entries = append(v.entries, FileEntry{Name: filepath.Base(dst)})

	return v.saveMetadata()
}
