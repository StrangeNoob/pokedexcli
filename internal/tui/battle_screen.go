package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/strangenoob/pokedexcli/internal/battle"
	"github.com/strangenoob/pokedexcli/internal/pokedex"
)

type battleStep int

const (
	pickFirstStep battleStep = iota
	pickSecondStep
	animateStep
	doneStep
)

type battleTickMsg struct{}

func battleTick() tea.Cmd {
	return tea.Tick(700*time.Millisecond, func(time.Time) tea.Msg { return battleTickMsg{} })
}

type battleModel struct {
	deps   Deps
	step   battleStep
	names  []string
	cursor int

	firstName  string
	secondList []string

	result  battle.Result
	turnIdx int

	aName, bName   string
	aMaxHP, bMaxHP int
	aHP, bHP       int

	status string
}

func newBattleModel(deps Deps) battleModel {
	return battleModel{deps: deps, step: pickFirstStep, names: sortedCaught(deps.Dex)}
}

func (m battleModel) Init() tea.Cmd { return m.requestSelectionArt() }

// requestSelectionArt loads art for the currently highlighted selection list entry.
func (m battleModel) requestSelectionArt() tea.Cmd {
	var items []string
	switch m.step {
	case pickFirstStep:
		items = m.names
	case pickSecondStep:
		items = m.secondList
	default:
		return nil
	}
	if len(items) == 0 || m.cursor >= len(items) {
		return nil
	}
	return m.deps.Art.request(m.deps, items[m.cursor])
}

func toCombatantTUI(cp *pokedex.CaughtPokemon) battle.Combatant {
	return battle.Combatant{
		Name:    cp.Base.Name,
		HP:      cp.HP(),
		Attack:  cp.Attack(),
		Defense: cp.Defense(),
		Speed:   cp.Speed(),
		Types:   cp.TypeNames(),
	}
}

func hpBar(cur, max, width int) string {
	if max <= 0 {
		max = 1
	}
	filled := cur * width / max
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	ratio := float64(cur) / float64(max)
	barColor := lipgloss.Color("42") // green
	if ratio < 0.5 {
		barColor = lipgloss.Color("214") // yellow
	}
	if ratio < 0.25 {
		barColor = lipgloss.Color("196") // red
	}
	return lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", filled)) +
		dimStyle.Render(strings.Repeat("░", width-filled))
}

func (m battleModel) applyTurn(tn battle.TurnEvent) battleModel {
	switch tn.Defender {
	case m.aName:
		m.aHP = tn.DefenderHP
	case m.bName:
		m.bHP = tn.DefenderHP
	}
	return m
}

func (m battleModel) startBattle(aName, bName string) (screenModel, tea.Cmd) {
	a, _ := m.deps.Dex.Get(aName)
	b, _ := m.deps.Dex.Get(bName)
	ca := toCombatantTUI(a)
	cb := toCombatantTUI(b)
	m.result = battle.Simulate(ca, cb, m.deps.RNG)
	m.aName, m.bName = aName, bName
	m.aMaxHP, m.aHP = ca.HP, ca.HP
	m.bMaxHP, m.bHP = cb.HP, cb.HP
	m.turnIdx = 0
	m.step = animateStep
	return m, tea.Batch(
		battleTick(),
		m.deps.Art.request(m.deps, aName),
		m.deps.Art.request(m.deps, bName),
	)
}

func (m battleModel) finish() (screenModel, tea.Cmd) {
	if m.step == doneStep {
		return m, nil
	}
	m.step = doneStep
	winner, wok := m.deps.Dex.Get(m.result.Winner)
	loser, lok := m.deps.Dex.Get(m.result.Loser)
	if wok && lok {
		xp := loser.Level * 10
		winner.AddXP(xp)
		m.status = fmt.Sprintf("%s wins! +%d XP (Lv%d)", winner.Base.Name, xp, winner.Level)
		_ = m.deps.Dex.Save(m.deps.SavePath)
	}
	return m, nil
}

func (m battleModel) Update(msg tea.Msg) (screenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case artLoadedMsg:
		m.deps.Art.handle(msg)
		return m, nil
	case battleTickMsg:
		if m.step != animateStep {
			return m, nil
		}
		if m.turnIdx >= len(m.result.Turns) {
			return m.finish()
		}
		m = m.applyTurn(m.result.Turns[m.turnIdx])
		m.turnIdx++
		if m.turnIdx >= len(m.result.Turns) {
			return m.finish()
		}
		return m, battleTick()
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m battleModel) handleKey(key tea.KeyMsg) (screenModel, tea.Cmd) {
	switch m.step {
	case pickFirstStep:
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, m.requestSelectionArt()
		case "down", "j":
			if m.cursor < len(m.names)-1 {
				m.cursor++
			}
			return m, m.requestSelectionArt()
		case "enter":
			if len(m.names) < 2 {
				m.status = "You need at least 2 caught Pokémon"
				return m, nil
			}
			m.firstName = m.names[m.cursor]
			m.secondList = nil
			for _, n := range m.names {
				if n != m.firstName {
					m.secondList = append(m.secondList, n)
				}
			}
			m.cursor = 0
			m.step = pickSecondStep
			return m, m.requestSelectionArt()
		}
	case pickSecondStep:
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, m.requestSelectionArt()
		case "down", "j":
			if m.cursor < len(m.secondList)-1 {
				m.cursor++
			}
			return m, m.requestSelectionArt()
		case "enter":
			return m.startBattle(m.firstName, m.secondList[m.cursor])
		}
	case animateStep:
		if key.String() == " " {
			for m.turnIdx < len(m.result.Turns) {
				m = m.applyTurn(m.result.Turns[m.turnIdx])
				m.turnIdx++
			}
			return m.finish()
		}
	}
	return m, nil
}

func (m battleModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Battle") + "\n\n")
	switch m.step {
	case pickFirstStep:
		b.WriteString("Choose your Pokémon:\n\n")
		b.WriteString(m.selectionView(m.names))
		b.WriteString("\n" + helpStyle.Render("↑/↓ move · enter pick · esc back"))
	case pickSecondStep:
		b.WriteString(fmt.Sprintf("First: %s. Choose the opponent:\n\n", m.firstName))
		b.WriteString(m.selectionView(m.secondList))
		b.WriteString("\n" + helpStyle.Render("↑/↓ move · enter pick · esc back"))
	case animateStep, doneStep:
		b.WriteString(m.battlefieldView() + "\n\n")
		shown := m.turnIdx + 1
		if m.step == doneStep {
			shown = len(m.result.Log)
		}
		if shown > len(m.result.Log) {
			shown = len(m.result.Log)
		}
		for _, line := range m.result.Log[:shown] {
			b.WriteString(line + "\n")
		}
		if m.step == doneStep {
			b.WriteString("\n" + statusStyle.Render(m.status) + "\n" + helpStyle.Render("esc back"))
		} else {
			b.WriteString("\n" + helpStyle.Render("space skip · esc back"))
		}
		return b.String()
	}
	if m.status != "" {
		b.WriteString("\n" + statusStyle.Render(m.status))
	}
	return b.String()
}

// battlefieldView renders the two combatants as equal-size side-by-side panels.
func (m battleModel) battlefieldView() string {
	h := lipgloss.Height(m.deps.Art.get(m.aName))
	if hb := lipgloss.Height(m.deps.Art.get(m.bName)); hb > h {
		h = hb
	}
	left := m.combatantPanel(m.aName, m.aHP, m.aMaxHP, h)
	right := m.combatantPanel(m.bName, m.bHP, m.bMaxHP, h)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, "   ", right)
}

// combatantPanel places a fixed-size sprite box beside the combatant's stats.
func (m battleModel) combatantPanel(name string, hp, maxHP, artH int) string {
	placed := lipgloss.Place(artWidth, artH, lipgloss.Center, lipgloss.Center, m.deps.Art.get(name))
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(placed)
	return lipgloss.JoinHorizontal(lipgloss.Top, box, "  ", m.combatantStats(name, hp, maxHP))
}

// combatantStats renders the combatant's name, level, live HP bar, and battle stats.
func (m battleModel) combatantStats(name string, hp, maxHP int) string {
	lvl, atk, def, spd, types := 0, 0, 0, 0, ""
	if cp, ok := m.deps.Dex.Get(name); ok {
		lvl, atk, def, spd = cp.Level, cp.Attack(), cp.Defense(), cp.Speed()
		types = strings.Join(cp.TypeNames(), ", ")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s  Lv%d\n", name, lvl)
	fmt.Fprintf(&b, "%s  %d/%d\n", hpBar(hp, maxHP, 14), hp, maxHP)
	fmt.Fprintf(&b, "ATK %3d  DEF %3d\n", atk, def)
	fmt.Fprintf(&b, "SPD %3d\n", spd)
	b.WriteString("Types: " + types)
	return b.String()
}

// selectionView renders the choice list beside a preview of the highlighted entry.
func (m battleModel) selectionView(items []string) string {
	left := boxStyle.Render(renderChoiceList(items, m.cursor))
	if len(items) == 0 || m.cursor >= len(items) {
		return left
	}
	preview := m.selectionPreview(items[m.cursor])
	if preview == "" {
		return left
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", boxStyle.Render(preview))
}

// selectionPreview shows a caught Pokémon's stats, with its sprite beside them
// once the art has loaded.
func (m battleModel) selectionPreview(name string) string {
	stats := m.statBlock(name)
	art := m.deps.Art.get(name)
	if art == "" {
		return stats
	}
	placed := lipgloss.Place(artWidth, lipgloss.Height(art), lipgloss.Center, lipgloss.Center, art)
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(placed)
	return lipgloss.JoinHorizontal(lipgloss.Top, box, "  ", stats)
}

// statBlock renders a caught Pokémon's level, stats, and types.
func (m battleModel) statBlock(name string) string {
	cp, ok := m.deps.Dex.Get(name)
	if !ok {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s  Lv%d\n", name, cp.Level)
	fmt.Fprintf(&b, "HP  %3d  ATK %3d\n", cp.HP(), cp.Attack())
	fmt.Fprintf(&b, "DEF %3d  SPD %3d\n", cp.Defense(), cp.Speed())
	b.WriteString("Types: " + strings.Join(cp.TypeNames(), ", "))
	return b.String()
}

func renderChoiceList(items []string, cursor int) string {
	var b strings.Builder
	for i, name := range items {
		c := "  "
		line := name
		if i == cursor {
			c = "▸ "
			line = selectedStyle.Render(line)
		}
		b.WriteString(c + line + "\n")
	}
	return b.String()
}
