package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
)

func Wish(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := new(model)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	width     int
	height    int
	company   company
	component string
	view      string
}

type childmsg struct{}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		// Update child dan simpan hasilnya
		m.company, cmd = m.company.Update(msg)
		return m, cmd

	case childmsg:
		m.view = ""
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "c":
			m.view = "c"
			m.company.load()
			return m, nil
		}
	}

	if m.view == "c" {
		m.company, cmd = m.company.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.view {
	case "c":
		return m.company.View()
	default:
		return "hshshshh"
	}
}
