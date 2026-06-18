package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/strangenoob/pokedexcli/internal/pokeapi"
	"github.com/strangenoob/pokedexcli/internal/pokedex"
)

func mkStat(name string, v int) pokeapi.Stat {
	var s pokeapi.Stat
	s.BaseStat = v
	s.Stat.Name = name
	return s
}

func statPokemon(name string, hp, atk, def, spd int) pokeapi.Pokemon {
	p := pokeapi.Pokemon{Name: name}
	p.Stats = []pokeapi.Stat{
		mkStat("hp", hp), mkStat("attack", atk), mkStat("defense", def), mkStat("speed", spd),
	}
	return p
}

func TestBattleSkipAwardsXP(t *testing.T) {
	d := testDeps()
	d.Dex.Add(statPokemon("alpha", 100, 90, 30, 60))
	d.Dex.Add(statPokemon("beta", 40, 20, 30, 10))
	m := newBattleModel(d)

	// pick first (alpha at cursor 0 after sort)
	sm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = sm.(battleModel)
	if m.step != pickSecondStep {
		t.Fatalf("step after first pick = %d", m.step)
	}
	// pick opponent and start
	sm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = sm.(battleModel)
	if m.step != animateStep || cmd == nil {
		t.Fatalf("expected animateStep with tick cmd, step=%d cmd=%v", m.step, cmd)
	}
	// fast-forward
	sm, _ = m.Update(runeKey(' '))
	m = sm.(battleModel)
	if m.step != doneStep {
		t.Fatalf("step after skip = %d, want doneStep", m.step)
	}
	w, _ := d.Dex.Get(m.result.Winner)
	if w.XP <= pokedex.XPForLevel(pokedex.StartingLevel) {
		t.Fatalf("winner XP not increased: %d", w.XP)
	}
}

func TestBattleTickAdvances(t *testing.T) {
	d := testDeps()
	d.Dex.Add(statPokemon("alpha", 200, 90, 30, 60))
	d.Dex.Add(statPokemon("beta", 200, 90, 30, 10))
	m := newBattleModel(d)
	sm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = sm.(battleModel)
	sm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = sm.(battleModel)
	start := m.turnIdx
	sm, _ = m.Update(battleTickMsg{})
	m = sm.(battleModel)
	if m.turnIdx <= start && m.step != doneStep {
		t.Fatalf("tick did not advance: turnIdx=%d", m.turnIdx)
	}
}
