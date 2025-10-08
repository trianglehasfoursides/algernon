package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/ssh"
)

func Wish(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := model{}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	component string
	view      string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c":
			m.view = "companies"
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.view {
	case "companies":
		table.New().Render()
	}

	return m.component
}
