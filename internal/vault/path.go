package vault

import (
	"os"
	"path/filepath"
	"strings"
)

func expandPath(base string) string {
	if base == "~" {
		if h, err := os.UserHomeDir(); err == nil {
			return h
		}
		return base
	}

	if strings.HasPrefix(base, "~/") || strings.HasPrefix(base, "~\\") {
		if h, err := os.UserHomeDir(); err == nil {
			return filepath.Join(h, base[2:])
		}
	}

	return base
}

func masterKeyPath(base string) string {
	return filepath.Join(base, masterKeyFile)
}

func metadataPath(base string) string {
	return filepath.Join(base, metadataFile)
}

func MasterKeyExists(base string) (bool, error) {
	expanded := expandPath(base)
	_, err := os.Stat(masterKeyPath(expanded))
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func ExpandPath(base string) string {
	return expandPath(base)
}
