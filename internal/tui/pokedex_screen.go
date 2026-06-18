package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/strangenoob/pokedexcli/internal/pokedex"
)

type pokedexModel struct {
	deps   Deps
	names  []string
	cursor int
	status string
}

func newPokedexModel(deps Deps) pokedexModel {
	return pokedexModel{deps: deps, names: sortedCaught(deps.Dex)}
}

func sortedCaught(dex *pokedex.Pokedex) []string {
	names := make([]string, 0, len(dex.Caught))
	for n := range dex.Caught {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func (m pokedexModel) Init() tea.Cmd { return nil }

func (m pokedexModel) Update(msg tea.Msg) (screenModel, tea.Cmd) {
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
		if m.cursor < len(m.names)-1 {
			m.cursor++
		}
	case "p":
		if len(m.names) == 0 {
			return m, nil
		}
		name := m.names[m.cursor]
		if m.deps.Dex.InParty(name) {
			_ = m.deps.Dex.RemoveFromParty(name)
			m.status = name + " left the party"
		} else if err := m.deps.Dex.AddToParty(name); err != nil {
			m.status = err.Error()
		} else {
			m.status = name + " joined the party"
		}
		_ = m.deps.Dex.Save(m.deps.SavePath)
	}
	return m, nil
}

func (m pokedexModel) View() string {
	if len(m.names) == 0 {
		return titleStyle.Render("Pokédex") + "\n\n" +
			dimStyle.Render("No Pokémon caught yet.") + "\n\n" +
			helpStyle.Render("esc back")
	}

	var list strings.Builder
	for i, name := range m.names {
		marker := "  "
		star := " "
		if m.deps.Dex.InParty(name) {
			star = "★"
		}
		cp, _ := m.deps.Dex.Get(name)
		line := fmt.Sprintf("%-12s Lv%d %s", name, cp.Level, star)
		if i == m.cursor {
			marker = "▸ "
			line = selectedStyle.Render(line)
		}
		list.WriteString(marker + line + "\n")
	}

	cp, _ := m.deps.Dex.Get(m.names[m.cursor])
	body := lipgloss.JoinHorizontal(lipgloss.Top,
		boxStyle.Render(list.String()),
		boxStyle.Render(detailView(cp, m.deps.Dex.InParty(cp.Base.Name))),
	)

	out := titleStyle.Render("Pokédex") + "\n\n" + body + "\n"
	if m.status != "" {
		out += statusStyle.Render(m.status) + "\n"
	}
	return out + helpStyle.Render("↑/↓ move · p party · esc back")
}

func detailView(cp *pokedex.CaughtPokemon, inParty bool) string {
	var b strings.Builder
	b.WriteString(cp.Base.Name + "\n")
	b.WriteString(fmt.Sprintf("Level %d   XP %d\n", cp.Level, cp.XP))
	b.WriteString(fmt.Sprintf("HP  %3d   ATK %3d\n", cp.HP(), cp.Attack()))
	b.WriteString(fmt.Sprintf("DEF %3d   SPD %3d\n", cp.Defense(), cp.Speed()))
	b.WriteString("Types: " + strings.Join(cp.TypeNames(), ", ") + "\n")
	if inParty {
		b.WriteString("★ in party")
	} else {
		b.WriteString("  not in party")
	}
	return b.String()
}
