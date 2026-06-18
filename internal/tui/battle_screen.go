package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

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

func (m battleModel) Init() tea.Cmd { return nil }

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
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
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
	return m, battleTick()
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
		case "down", "j":
			if m.cursor < len(m.names)-1 {
				m.cursor++
			}
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
		}
	case pickSecondStep:
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.secondList)-1 {
				m.cursor++
			}
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
		b.WriteString(renderChoiceList(m.names, m.cursor))
		b.WriteString("\n" + helpStyle.Render("↑/↓ move · enter pick · esc back"))
	case pickSecondStep:
		b.WriteString(fmt.Sprintf("First: %s. Choose the opponent:\n\n", m.firstName))
		b.WriteString(renderChoiceList(m.secondList, m.cursor))
		b.WriteString("\n" + helpStyle.Render("↑/↓ move · enter pick · esc back"))
	case animateStep, doneStep:
		b.WriteString(fmt.Sprintf("%-12s %s %d/%d\n", m.aName, hpBar(m.aHP, m.aMaxHP, 12), m.aHP, m.aMaxHP))
		b.WriteString(fmt.Sprintf("%-12s %s %d/%d\n", m.bName, hpBar(m.bHP, m.bMaxHP, 12), m.bHP, m.bMaxHP))
		b.WriteString("\n")
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
