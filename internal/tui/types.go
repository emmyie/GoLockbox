package tui

import "GoLockbox/internal/vault"

type ScreenID int

const (
	ScreenMenu ScreenID = iota
	ScreenVaultPath
	ScreenPassword
	ScreenBrowser
	ScreenFileView
	ScreenFileEditor
	ScreenAddFile
	ScreenRenameFile
	ScreenNewFileName
)

type VaultMode int

const (
	VaultModeCreate VaultMode = iota
	VaultModeOpen
)

type ScreenMsg struct {
	NextScreen ScreenID
	Payload    any
}

type MenuChoice int

const (
	MenuNewVault MenuChoice = iota
	MenuOpenVault
)

type VaultPathResult struct {
	Mode VaultMode
	Path string
}

type PasswordResult struct {
	Mode     VaultMode
	Path     string
	Password string
}

type OpenFileMsg struct {
	Index int
}

type EditFileMsg struct {
	Index int
}

type DeleteFileMsg struct {
	Index int
}

type RenameRequestMsg struct {
	Index int
	Name  string
}

type RenameResultMsg struct {
	Index   int
	NewName string
}

type SaveFileMsg struct {
	Index int
	Data  []byte
	Name  string
}

type AddFileMsg struct {
	SrcPath string
}

type BackToBrowserMsg struct{}

type VaultLoadedMsg struct {
	Vault *vault.Vault
	Files []vault.FileEntry
}

type NewFileNameResult struct {
	Name string
}
