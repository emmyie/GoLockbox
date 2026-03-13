package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"GoLockbox/internal/vault"
)

type BrowserModel struct {
	viewport viewport.Model
	cursor   int
	files    []vault.FileEntry
}

func NewBrowserModel() BrowserModel {
	vp := viewport.New(60, 20)
	return BrowserModel{
		viewport: vp,
	}
}

func (m *BrowserModel) SetFiles(files []vault.FileEntry) {
	m.files = files
	if m.cursor >= len(m.files) {
		m.cursor = len(m.files) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	m.syncViewport()
}

func (m *BrowserModel) syncViewport() {
	lines := make([]string, len(m.files))
	for i, f := range m.files {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		lines[i] = fmt.Sprintf("%s %s", cursor, f.Name)
	}
	if len(lines) == 0 {
		lines = []string{"(no files)"}
	}
	m.viewport.SetContent(strings.Join(lines, "\n"))
}

func (m BrowserModel) Init() tea.Cmd {
	return nil
}

func (m BrowserModel) Update(msg tea.Msg) (BrowserModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.syncViewport()
			}
		case "down":
			if m.cursor < len(m.files)-1 {
				m.cursor++
				m.syncViewport()
			}
		case "enter":
			idx := m.cursor
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenFileView,
					Payload:    OpenFileMsg{Index: idx},
				}
			}
		case "e":
			idx := m.cursor
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenFileEditor,
					Payload:    EditFileMsg{Index: idx},
				}
			}
		case "n":
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenNewFileName,
					Payload:    nil,
				}
			}
		case "d":
			idx := m.cursor
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenBrowser,
					Payload:    DeleteFileMsg{Index: idx},
				}
			}
		case "r":
			if len(m.files) == 0 {
				return m, cmd
			}
			idx := m.cursor
			name := m.files[idx].Name
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenRenameFile,
					Payload:    RenameRequestMsg{Index: idx, Name: name},
				}
			}
		case "a":
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenAddFile,
					Payload:    nil,
				}
			}
		case "b":
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenMenu,
					Payload:    nil,
				}
			}
		}
	}

	return m, cmd
}

func (m BrowserModel) View(v *vault.Vault) string {
	s := "\n"
	if v != nil {
		s += fmt.Sprintf("Vault: %s\n\n", v.Path())
	} else {
		s += "No vault open.\n\n"
	}
	s += m.viewport.View()
	s += "\n\n↑/↓ move, Enter open file, e edit file, n new file, d delete file, r rename file, a add file, b go to menu\n"
	return s
}
