package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kancli/pkg/column"
	"kancli/pkg/keymap"
)

var keys = keymap.NewKeyMap()

func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case column.Column:
		f.Col = msg
		f.Col.List.Index()
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Quit):
			return f, tea.Quit

		case key.Matches(msg, keys.Back):
			return f, nil

		case key.Matches(msg, keys.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}

			return f, nil
		}
	}

	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	}
	f.description, cmd = f.description.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create a new task",
		f.title.View(),
		f.description.View(),
		f.help.View(keys))
}
