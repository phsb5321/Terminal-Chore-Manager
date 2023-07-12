package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Status int

type KeyMap struct{ New, Edit, Delete, Up, Down, Right, Left, Enter, Help, Quit, Back key.Binding }

type Task struct {
	status             Status
	title, description string
}

func NewTask(status Status, title, description string) Task {
	return Task{status: status, title: title, description: description}
}

type Column struct {
	focus  bool
	status Status
	list   list.Model
	height int
	width  int
}

func NewColumn(status Status) Column {
	var focus bool
	if status == todo {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return Column{focus: focus, status: status, list: defaultList}
}

type Board struct {
	help     help.Model
	loaded   bool
	focused  Status
	cols     []Column
	quitting bool
}

func NewBoard() *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: todo}
}

type Form struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	col         Column
	index       int
}

func NewForm(title, description string) *Form {
	form := Form{
		help:        help.New(),
		title:       textinput.New(),
		description: textarea.New(),
	}
	form.title.Placeholder = title
	form.description.Placeholder = description
	form.title.Focus()
	return &form
}

type MoveMsg struct {
	Task
}

func (s Status) getNext() Status {
	if s == done {
		return todo
	}
	return s + 1
}

func (s Status) getPrev() Status {
	if s == todo {
		return done
	}
	return s - 1
}

const margin = 4

var board *Board

const (
	todo Status = iota
	inProgress
	done
)

// Provides the mock data to fill the kanban board

func (b *Board) initLists() {
	b.cols = []Column{
		NewColumn(todo),
		NewColumn(inProgress),
		NewColumn(done),
	}
	// Init To Do
	b.cols[todo].list.Title = "To Do"
	b.cols[todo].list.SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
	})
	// Init in progress
	b.cols[inProgress].list.Title = "In Progress"
	b.cols[inProgress].list.SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init done
	b.cols[done].list.Title = "Done"
	b.cols[done].list.SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(Column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case Form:
		return m, m.cols[m.focused].Set(msg.index, msg.CreateTask())
	case MoveMsg:
		return m, m.cols[m.focused.getNext()].Set(APPEND, msg.Task)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		}
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(Column); ok {
		m.cols[m.focused] = res.(Column)
	} else {
		return res, cmd
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
		m.cols[todo].View(),
		m.cols[inProgress].View(),
		m.cols[done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

var keys = KeyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/l", "move left"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

func newDefaultForm() *Form {
	return NewForm("task name", "")
}

func (f Form) CreateTask() Task {
	return Task{f.col.status, f.title.Value(), f.description.Value()}
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case Column:
		f.col = msg
		f.col.list.Index()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return f, tea.Quit

		case key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}
			// Return the completed form as a message.
			return board.Update(f)
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

const APPEND = -1

func (c *Column) Focus() {
	c.focus = true
}

func (c *Column) Blur() {
	c.focus = false
}

func (c *Column) Focused() bool {
	return c.focus
}

// Init does initial setup for the column.
func (c Column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height/2)

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Edit):
			if len(c.list.VisibleItems()) != 0 {
				task := c.list.SelectedItem().(Task)
				f := NewForm(task.title, task.description)
				f.index = c.list.Index()
				f.col = c
				return f.Update(nil)
			}

		case key.Matches(msg, keys.New):
			f := newDefaultForm()
			f.index = APPEND
			f.col = c
			return f.Update(nil)

		case key.Matches(msg, keys.Delete):
			return c, c.DeleteCurrent()

		case key.Matches(msg, keys.Enter):
			return c, c.MoveToNext()

		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c Column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *Column) DeleteCurrent() tea.Cmd {
	if len(c.list.VisibleItems()) > 0 {
		c.list.RemoveItem(c.list.Index())
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)
	return cmd
}

func (c *Column) Set(i int, t Task) tea.Cmd {
	if i != APPEND {
		return c.list.SetItem(i, t)
	}
	return c.list.InsertItem(APPEND, t)
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
	var task Task
	var ok bool
	// If nothing is selected, the SelectedItem will return Nil.
	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}
	// move item
	c.list.RemoveItem(c.list.Index())
	task.status = c.status.getNext()

	// refresh list
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return MoveMsg{task} })
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// FilterValue implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	board = NewBoard()
	board.initLists()
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
