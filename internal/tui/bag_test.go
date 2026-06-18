package tui

import (
	"strings"
	"testing"

	"github.com/strangenoob/pokedexcli/internal/ball"
)

func TestBagViewShowsCounts(t *testing.T) {
	m := newBagModel(testDeps())
	out := m.View()
	for _, name := range ball.Names() {
		if !strings.Contains(out, name) {
			t.Errorf("bag view missing %q", name)
		}
	}
	if !strings.Contains(out, "20") {
		t.Error("bag view should show 20 pokeballs")
	}
}
