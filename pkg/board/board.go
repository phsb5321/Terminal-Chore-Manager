package board

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kancli/pkg/column"
	"kancli/pkg/form"
	"kancli/pkg/keymap"
	"kancli/pkg/message"
	"kancli/pkg/task"
)

const margin = 4

const APPEND = -1

type Board struct {
	help     help.Model
	loaded   bool
	focused  task.Status
	cols     []column.Column
	quitting bool
}

func NewBoard() *Board {
	helpModel := help.New()
	helpModel.ShowAll = true
	newBoard := &Board{help: helpModel, focused: task.ToDo}
	newBoard.InitLists()
	return newBoard
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := 0; i < len(m.cols); i++ {
			res, cmd := m.cols[i].Update(msg)
			m.cols[i] = res.(column.Column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case form.Form:
		return m, m.cols[m.focused].Set(msg.Index, msg.CreateTask())
	case message.MoveMsg:
		return m, m.cols[m.focused.GetNext()].Set(APPEND, msg.Task)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keymap.Keys.Left):
			m.changeFocus(m.focused.GetPrev())
		case key.Matches(msg, keymap.Keys.Right):
			m.changeFocus(m.focused.GetNext())
		}
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(column.Column); ok {
		m.cols[m.focused] = res.(column.Column)
	}
	return m, cmd
}

func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[task.ToDo].View(),
		m.cols[task.InProgress].View(),
		m.cols[task.Done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keymap.Keys))
}

func (b *Board) InitLists() {
	b.cols = []column.Column{
		column.NewColumn(task.ToDo),
		column.NewColumn(task.InProgress),
		column.NewColumn(task.Done),
	}

	b.cols[task.ToDo].List.Title = "To Do"
	b.cols[task.ToDo].List.SetItems([]list.Item{
		task.NewTask(task.ToDo, "buy milk", "strawberry milk"),
	})

	b.cols[task.InProgress].List.Title = "In Progress"
	b.cols[task.InProgress].List.SetItems([]list.Item{
		task.NewTask(task.InProgress, "write code", "don't worry, it's Go"),
	})

	b.cols[task.Done].List.Title = "Done"
	b.cols[task.Done].List.SetItems([]list.Item{
		task.NewTask(task.Done, "stay cool", "as a cucumber"),
	})
}

func (m *Board) changeFocus(nextStatus task.Status) {
	m.cols[m.focused].Blur()
	m.focused = nextStatus
	m.cols[m.focused].Focus()
}
