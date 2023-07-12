package column

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kancli/pkg/keymap"
	"kancli/pkg/message"
	"kancli/pkg/task"
)

const (
	APPEND = -1
	margin = 4
)

type Task = task.Task

type MoveMsg struct {
	Task task.Task
}

func (c *Column) Focus() {
	c.focus = true
}

func (c *Column) Blur() {
	c.focus = false
}

func (c *Column) Focused() bool {
	return c.focus
}

func (c Column) Init() tea.Cmd {
	return nil
}

func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.List.SetSize(msg.Width/margin, msg.Height/2)

	case tea.KeyMsg:
		return c.handleKeyMsg(msg)

	default:
		c.List, cmd = c.List.Update(msg)
	}
	return c, cmd
}

func (c *Column) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {

	case key.Matches(msg, keymap.Keys.Edit):
	case key.Matches(msg, keymap.Keys.Delete):
		return c, c.DeleteCurrent()

	case key.Matches(msg, keymap.Keys.Enter):
		return c, c.MoveToNext()
	}
	return c, nil
}

func (c Column) View() string {
	return c.getStyle().Render(c.List.View())
}

func (c *Column) DeleteCurrent() tea.Cmd {
	if len(c.List.VisibleItems()) > 0 {
		c.List.RemoveItem(c.List.Index())
	}

	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)
	return cmd
}

func (c *Column) Set(i int, t task.Task) tea.Cmd {
	if i != APPEND {
		return c.List.SetItem(i, t)
	}
	return c.List.InsertItem(APPEND, t)
}

func (c *Column) setSize(width, height int) {
	c.width = width / margin
}

func (c *Column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}

	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

func (c *Column) MoveToNext() tea.Cmd {
	task, ok := c.List.SelectedItem().(Task)
	if !ok {
		return nil
	}
	c.List.RemoveItem(c.List.Index())
	task.Status = c.Status.GetNext()

	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg {
		return message.MoveMsg{Task: task}
	})
}
