package tui

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/strangenoob/pokedexcli/internal/pokedex"
)

// Shared test helpers for the tui package.
func testDeps() Deps {
	return Deps{
		Dex:      pokedex.New(),
		Client:   nil,
		RNG:      rand.New(rand.NewSource(1)),
		SavePath: filepath.Join(os.TempDir(), "pokedexcli-tui-test.json"),
	}
}

func runeKey(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func TestMenuNavigation(t *testing.T) {
	m := newMenuModel()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Fatalf("cursor after down = %d, want 1", m.cursor)
	}
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Fatalf("cursor after up = %d, want 0", m.cursor)
	}
}

func TestMenuSelectEmitsSwitch(t *testing.T) {
	m := newMenuModel()
	m.cursor = 3 // Bag
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	sw, ok := cmd().(switchMsg)
	if !ok || sw.to != bagScreen {
		t.Fatalf("expected switch to bagScreen, got %#v", cmd())
	}
}

func TestRootSwitchAndBack(t *testing.T) {
	m := newRootModel(testDeps())
	updated, _ := m.Update(switchMsg{to: bagScreen})
	rm := updated.(rootModel)
	if rm.screen != bagScreen {
		t.Fatalf("screen = %d, want bagScreen", rm.screen)
	}
	updated, _ = rm.Update(tea.KeyMsg{Type: tea.KeyEsc})
	rm = updated.(rootModel)
	if rm.screen != menuScreen {
		t.Fatalf("screen after esc = %d, want menuScreen", rm.screen)
	}
}
