package tui

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"GoLockbox/internal/vault"
)

type RootModel struct {
	screen ScreenID

	menu        MenuModel
	vaultPath   VaultPathModel
	password    PasswordModel
	browser     BrowserModel
	fileView    FileViewModel
	editor      EditorModel
	addFile     AddFileModel
	rename      RenameModel
	newFileName NewFileNameModel

	vault    *vault.Vault
	files    []vault.FileEntry
	mode     VaultMode
	vaultErr string
	path     string
}

func NewRootModel() RootModel {
	return RootModel{
		screen:      ScreenMenu,
		menu:        NewMenuModel(),
		vaultPath:   NewVaultPathModel(),
		password:    NewPasswordModel(),
		browser:     NewBrowserModel(),
		fileView:    NewFileViewModel(),
		editor:      NewEditorModel(),
		addFile:     NewAddFileModel(),
		rename:      NewRenameModel(),
		newFileName: NewNewFileNameModel(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.menu.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			if m.vault != nil {
				_ = m.vault.Close()
			}
			return m, tea.Quit
		}
	case ScreenMsg:
		return m.handleScreenMsg(msg)
	}

	switch m.screen {
	case ScreenMenu:
		var cmd tea.Cmd
		m.menu, cmd = m.menu.Update(msg)
		return m, cmd

	case ScreenVaultPath:
		var cmd tea.Cmd
		m.vaultPath, cmd = m.vaultPath.Update(msg, m.mode)
		return m, cmd

	case ScreenPassword:
		var cmd tea.Cmd
		m.password, cmd = m.password.Update(msg, m.mode, m.path)
		return m, cmd

	case ScreenBrowser:
		var cmd tea.Cmd
		m.browser, cmd = m.browser.Update(msg)
		return m, cmd

	case ScreenFileView:
		var cmd tea.Cmd
		m.fileView, cmd = m.fileView.Update(msg)
		return m, cmd

	case ScreenFileEditor:
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		return m, cmd

	case ScreenAddFile:
		var cmd tea.Cmd
		m.addFile, cmd = m.addFile.Update(msg)
		return m, cmd

	case ScreenRenameFile:
		var cmd tea.Cmd
		m.rename, cmd = m.rename.Update(msg)
		return m, cmd

	case ScreenNewFileName:
		var cmd tea.Cmd
		m.newFileName, cmd = m.newFileName.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m RootModel) handleScreenMsg(sm ScreenMsg) (tea.Model, tea.Cmd) {
	switch payload := sm.Payload.(type) {

	case MenuChoice:
		if payload == MenuNewVault {
			m.mode = VaultModeCreate
		} else {
			m.mode = VaultModeOpen
		}
		m.vaultPath.Reset()
		m.screen = ScreenVaultPath
		return m, m.vaultPath.Init()

	case VaultPathResult:
		m.mode = payload.Mode
		m.path = payload.Path
		m.password.Reset()
		m.screen = ScreenPassword
		return m, m.password.Init()

	case PasswordResult:
		m.mode = payload.Mode
		m.path = payload.Path
		pass := payload.Password

		var v *vault.Vault
		var err error
		if m.mode == VaultModeCreate {
			v, err = vault.CreateVault(m.path, pass)
		} else {
			v, err = vault.OpenVault(m.path, pass)
		}
		if err != nil {
			m.vaultErr = err.Error()
			m.menu.Reset()
			m.screen = ScreenMenu
			return m, nil
		}
		m.vault = v
		m.files = v.ListFiles()
		m.browser.SetFiles(m.files)
		m.vaultErr = ""
		m.screen = ScreenBrowser
		return m, nil

	case OpenFileMsg:
		if m.vault == nil || payload.Index < 0 || payload.Index >= len(m.files) {
			return m, nil
		}
		data, err := m.vault.OpenFile(payload.Index)
		if err != nil {
			m.vaultErr = err.Error()
			return m, nil
		}
		m.fileView.SetContent(data)
		m.screen = ScreenFileView
		return m, nil

	case EditFileMsg:
		if m.vault == nil || payload.Index < 0 || payload.Index >= len(m.files) {
			return m, nil
		}
		data, err := m.vault.OpenFile(payload.Index)
		if err != nil {
			m.vaultErr = err.Error()
			return m, nil
		}
		m.editor.SetFile(payload.Index, m.files[payload.Index].Name, data)
		m.screen = ScreenFileEditor
		return m, nil

	case DeleteFileMsg:
		if m.vault != nil && payload.Index >= 0 && payload.Index < len(m.files) {
			_ = m.vault.DeleteFile(payload.Index)
			m.files = m.vault.ListFiles()
			m.browser.SetFiles(m.files)
		}
		m.screen = ScreenBrowser
		return m, nil

	case RenameRequestMsg:
		m.rename.SetTarget(payload.Index, payload.Name)
		m.screen = ScreenRenameFile
		return m, m.rename.Init()

	case RenameResultMsg:
		if m.vault != nil && payload.Index >= 0 && payload.Index < len(m.files) {
			_ = m.vault.RenameFile(payload.Index, payload.NewName)
			m.files = m.vault.ListFiles()
			m.browser.SetFiles(m.files)
		}
		m.screen = ScreenBrowser
		return m, nil

	case SaveFileMsg:
		if payload.Index == -1 {
			_ = m.vault.CreateFile(payload.Name, payload.Data)
		} else {
			_ = m.vault.SaveFile(payload.Index, payload.Data)
		}

		m.files = m.vault.ListFiles()
		m.browser.SetFiles(m.files)
		m.screen = ScreenBrowser
		return m, nil

	case AddFileMsg:
		if m.vault != nil {
			name := filepath.Base(payload.SrcPath)
			if name == "" {
				name = payload.SrcPath
			}
			_ = m.vault.AddFile(payload.SrcPath, name)
			m.files = m.vault.ListFiles()
			m.browser.SetFiles(m.files)
		}
		m.screen = ScreenBrowser
		return m, nil

	case BackToBrowserMsg:
		// Used by file view/editor/add/rename to go back to browser
		m.screen = ScreenBrowser
		return m, nil

	case NewFileNameResult:
		// User chose a filename → open editor with empty content
		m.editor.SetFile(-1, payload.Name, nil)
		m.screen = ScreenFileEditor
		return m, nil
	}

	// Fallback: just switch screen
	m.screen = sm.NextScreen
	switch m.screen {
	case ScreenMenu:
		m.menu.Reset()
	case ScreenVaultPath:
		m.vaultPath.Reset()
		return m, m.vaultPath.Init()
	case ScreenPassword:
		m.password.Reset()
		return m, m.password.Init()
	case ScreenAddFile:
		return m, m.addFile.Init()
	case ScreenRenameFile:
		return m, m.rename.Init()
	case ScreenNewFileName:
		return m, m.newFileName.Init()
	}

	return m, nil
}

func (m RootModel) View() string {
	switch m.screen {
	case ScreenMenu:
		return m.menu.View(m.vaultErr)
	case ScreenVaultPath:
		return m.vaultPath.View(m.mode)
	case ScreenPassword:
		return m.password.View(m.mode, m.path)
	case ScreenBrowser:
		return m.browser.View(m.vault)
	case ScreenFileView:
		return m.fileView.View()
	case ScreenFileEditor:
		return m.editor.View()
	case ScreenAddFile:
		return m.addFile.View()
	case ScreenRenameFile:
		return m.rename.View()
	case ScreenNewFileName:
		return m.newFileName.View()
	default:
		return "\nUnknown screen\n"
	}
}

func Run() error {
	p := tea.NewProgram(NewRootModel())
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}
	return nil
}
