package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	label string
	desc  string
	to    screen
}

type menuModel struct {
	items  []menuItem
	cursor int
}

func newMenuModel() menuModel {
	return menuModel{
		items: []menuItem{
			{"Pokédex", "browse caught Pokémon", pokedexScreen},
			{"Explore", "catch wild Pokémon", exploreScreen},
			{"Battle", "test your team", battleScreen},
			{"Bag", "your Pokéballs", bagScreen},
			{"Quit", "exit", menuScreen},
		},
	}
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (menuModel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}
	case "enter":
		item := m.items[m.cursor]
		if item.label == "Quit" {
			return m, quitCmd
		}
		return m, switchTo(item.to)
	}
	return m, nil
}

func (m menuModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Pokédex CLI") + "\n\n")
	for i, item := range m.items {
		cursor := "  "
		line := fmt.Sprintf("%-10s %s", item.label, dimStyle.Render(item.desc))
		if i == m.cursor {
			cursor = "▸ "
			line = selectedStyle.Render(fmt.Sprintf("%-10s", item.label)) + " " + dimStyle.Render(item.desc)
		}
		b.WriteString(cursor + line + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ move · enter select · q quit"))
	return b.String()
}
