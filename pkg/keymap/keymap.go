package keymap

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	New, Edit, Delete, Up, Down, Right, Left, Enter, Help, Quit, Back key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

func NewKeyMap() KeyMap {
	return KeyMap{
		New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
		Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
		Right:  key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "move right")),
		Left:   key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "move left")),
		Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "enter")),
		Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
		Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

var Keys = NewKeyMap()
