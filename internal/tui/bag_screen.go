package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/strangenoob/pokedexcli/internal/ball"
)

type bagModel struct {
	deps Deps
}

func newBagModel(deps Deps) bagModel { return bagModel{deps: deps} }

func (m bagModel) Init() tea.Cmd { return nil }

func (m bagModel) Update(tea.Msg) (screenModel, tea.Cmd) { return m, nil }

func (m bagModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Bag") + "\n\n")
	for _, name := range ball.Names() {
		b.WriteString(fmt.Sprintf("  %-10s × %2d\n", name, m.deps.Dex.BallCount(name)))
	}
	b.WriteString("\n" + helpStyle.Render("esc back"))
	return b.String()
}
