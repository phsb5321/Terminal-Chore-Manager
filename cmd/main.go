package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"kancli/pkg/board"
)

func main() {
	logFile, err := setupLogging()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logging: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	kanbanBoard := board.NewBoard()

	err = runProgram(kanbanBoard)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run program: %v\n", err)
		os.Exit(1)
	}
}

func setupLogging() (*os.File, error) {
	return tea.LogToFile("debug.log", "debug")
}

func runProgram(model tea.Model) error {
	program := tea.NewProgram(model)
	_, err := program.Run()
	return err
}
