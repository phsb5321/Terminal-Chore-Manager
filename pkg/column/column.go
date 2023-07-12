package column

import (
	"github.com/charmbracelet/bubbles/list"

	"kancli/pkg/task"
)

type Column struct {
	focus  bool
	Status task.Status
	List   list.Model
	height int
	width  int
}

func NewColumn(status task.Status) Column {
	var focus bool
	if status == task.ToDo {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return Column{focus: focus, Status: status, List: defaultList}
}
