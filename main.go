package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	m := NewModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
}
