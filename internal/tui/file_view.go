package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type FileViewModel struct {
	viewport viewport.Model
}

func NewFileViewModel() FileViewModel {
	vp := viewport.New(80, 20)
	return FileViewModel{viewport: vp}
}

func (m *FileViewModel) SetContent(data []byte) {
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("%4d | %s", i+1, lines[i])
	}
	m.viewport.SetContent(strings.Join(lines, "\n"))
}

func (m FileViewModel) Init() tea.Cmd {
	return nil
}

func (m FileViewModel) Update(msg tea.Msg) (FileViewModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "b", "esc":
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenBrowser,
					Payload:    BackToBrowserMsg{},
				}
			}
		}
	}

	return m, cmd
}

func (m FileViewModel) View() string {
	s := "\nFile view – b/esc to go back\n\n"
	s += m.viewport.View()
	s += "\n"
	return s
}
