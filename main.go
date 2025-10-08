package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mdl struct {
	rows     [][]string
	selected int
}

func (m mdl) Init() tea.Cmd { return nil }

func (m mdl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(m.rows)-1 {
				m.selected++
			}
		case "enter":
			fmt.Println("Kamu memilih:", m.rows[m.selected][0])
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m mdl) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("213")).
		Render("📋 Daftar Lamaran")

	var out string
	for i, r := range m.rows {
		row := fmt.Sprintf("%-10s | %-15s | %-10s", r[0], r[1], r[2])
		if i == m.selected {
			// highlight row yang sedang dipilih
			row = lipgloss.NewStyle().
				Background(lipgloss.Color("63")).
				Foreground(lipgloss.Color("230")).
				Render(row)
		}
		out += row + "\n"
	}

	help := lipgloss.NewStyle().Faint(true).
		Render("↑/↓ untuk pindah, Enter pilih, q keluar")

	return lipgloss.JoinVertical(lipgloss.Left, header, "", out, help)
}

func main() {
	data := [][]string{
		{"1", "Bulbasaur", "Grass", "Poison", "フシギダネ", "Fushigidane"},
		{"2", "Ivysaur", "Grass", "Poison", "フシギソウ", "Fushigisou"},
		{"3", "Venusaur", "Grass", "Poison", "フシギバナ", "Fushigibana"},
		{"4", "Charmander", "Fire", "", "ヒトカゲ", "Hitokage"},
		{"5", "Charmeleon", "Fire", "", "リザード", "Lizardo"},
		{"6", "Charizard", "Fire", "Flying", "リザードン", "Lizardon"},
		{"7", "Squirtle", "Water", "", "ゼニガメ", "Zenigame"},
		{"8", "Wartortle", "Water", "", "カメール", "Kameil"},
		{"9", "Blastoise", "Water", "", "カメックス", "Kamex"},
		{"10", "Caterpie", "Bug", "", "キャタピー", "Caterpie"},
		{"11", "Metapod", "Bug", "", "トランセル", "Trancell"},
		{"12", "Butterfree", "Bug", "Flying", "バタフリー", "Butterfree"},
		{"13", "Weedle", "Bug", "Poison", "ビードル", "Beedle"},
		{"14", "Kakuna", "Bug", "Poison", "コクーン", "Cocoon"},
		{"15", "Beedrill", "Bug", "Poison", "スピアー", "Spear"},
		{"16", "Pidgey", "Normal", "Flying", "ポッポ", "Poppo"},
		{"17", "Pidgeotto", "Normal", "Flying", "ピジョン", "Pigeon"},
		{"18", "Pidgeot", "Normal", "Flying", "ピジョット", "Pigeot"},
		{"19", "Rattata", "Normal", "", "コラッタ", "Koratta"},
		{"20", "Raticate", "Normal", "", "ラッタ", "Ratta"},
		{"21", "Spearow", "Normal", "Flying", "オニスズメ", "Onisuzume"},
		{"22", "Fearow", "Normal", "Flying", "オニドリル", "Onidrill"},
		{"23", "Ekans", "Poison", "", "アーボ", "Arbo"},
		{"24", "Arbok", "Poison", "", "アーボック", "Arbok"},
		{"25", "Pikachu", "Electric", "", "ピカチュウ", "Pikachu"},
		{"26", "Raichu", "Electric", "", "ライチュウ", "Raichu"},
		{"27", "Sandshrew", "Ground", "", "サンド", "Sand"},
		{"28", "Sandslash", "Ground", "", "サンドパン", "Sandpan"},
	}

	c := companies{
		items: data,
	}

	tea.NewProgram(c, tea.WithAltScreen()).Run()
}
