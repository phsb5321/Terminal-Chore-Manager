package form

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"kancli/pkg/column"
	"kancli/pkg/task"
)

type Form struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	Col         column.Column
	Index       int
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

func (f Form) CreateTask() task.Task {
	return task.NewTask(f.Col.Status, f.title.Value(), f.description.Value())
}

func (f *Form) Init() tea.Cmd {
	return nil
}

func NewDefaultForm() *Form {
	return NewForm("task name", "")
}
