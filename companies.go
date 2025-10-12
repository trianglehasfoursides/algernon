package main

import (
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"gorm.io/gorm"
)

// gorm model
type Company struct {
	gorm.Model
	Name     string
	Industry string
	Location string
	Status   string
	Website  string
	Email    string
	Phone    string
	Notes    string
}

type company struct {
	width     int
	height    int
	items     [][]string
	records   []Company
	current   int
	selected  []string
	state     string
	form      tea.Model
	page      int
	pagesize  int
	totalpage int
}

func (c *company) load() {
	var total int64

	db.Model(&Company{}).Count(&total)

	_ = db.
		Limit(c.pagesize).
		Offset((c.page - 1) * c.pagesize).
		Find(&c.records)

	c.items = [][]string{}
	for _, r := range c.records {
		c.items = append(c.items, []string{
			r.Name,
			r.Email,
		})
	}

	c.totalpage = (int(total) + c.pagesize) / c.pagesize

}

func (c *company) Init() tea.Cmd {
	return nil
}

func (c *company) Update(msg tea.Msg) (company, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.height = msg.Height
		c.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return *c, func() tea.Msg {
				return childmsg{}
			}
		case "up":
			if c.current > 0 {
				c.current--
			}
		case "down":
			if c.current < len(c.items)-1 {
				c.current++
			}
		case "right":
			if c.page < c.totalpage {
				c.page++
				c.current = 0
				c.load()
			}
		case "left":
			if c.page > 1 {
				c.page--
				c.current = 0
				c.load()
			}
		case "enter":
			if len(c.items) > 0 {
				c.selected = c.items[c.current]
			}
		case "c":
			c.selected = nil
		case "a":
			c.state = "add"
			c.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Company name"),
				),
			)
		}
	}
	return *c, nil
}

func (c *company) View() string {
	switch c.state {
	case "add":
		c.form = huh.NewForm(
			huh.NewGroup(
				huh.NewInput(),
			),
		)

		return c.form.View()
	}

	// default
	re := lipgloss.NewRenderer(os.Stdout)

	baseStyle := re.NewStyle().Padding(0, 1)

	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)

	selectedStyle := baseStyle.Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))

	headers := []string{
		"Name",
		"Industry",
		"Location",
		"Status",
		"Website",
		"Email",
		"Phone",
		"Notes",
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(80).
		Rows(c.items...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}

			// highlight baris yang dipilih
			if row == c.current {
				return selectedStyle
			}

			even := row%2 == 0

			if even {
				return baseStyle.Foreground(lipgloss.Color("245"))
			}
			return baseStyle.Foreground(lipgloss.Color("252"))
		})

	total := lipgloss.NewStyle().MarginTop(4).Render("total pages : ", strconv.Itoa(c.totalpage))
	v := lipgloss.JoinVertical(lipgloss.Center, t.Render(), total)
	return lipgloss.Place(
		c.width, c.height,
		lipgloss.Center,
		lipgloss.Center,
		v,
	)
}
