package main

import (
	tea "github.com/charmbracelet/bubbletea"
	say "github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var err error
	if db, err = gorm.Open(sqlite.Open("algernon.db"), &gorm.Config{}); err != nil {
		say.Fatal(err.Error())
	}

	// migration
	db.AutoMigrate(&Company{})

	// start!!!
	m := &model{
		company: company{
			page:     1,
			pagesize: 10,
		},
	}

	tea.NewProgram(m, tea.WithAltScreen()).Run()
}
