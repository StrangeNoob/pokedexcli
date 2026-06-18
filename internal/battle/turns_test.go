package battle

import (
	"math/rand"
	"testing"
)

func TestSimulatePopulatesTurns(t *testing.T) {
	a := Combatant{Name: "pikachu", HP: 60, Attack: 55, Defense: 40, Speed: 90, Types: []string{"electric"}}
	b := Combatant{Name: "geodude", HP: 80, Attack: 80, Defense: 100, Speed: 20, Types: []string{"rock", "ground"}}

	r := Simulate(a, b, rand.New(rand.NewSource(42)))

	if len(r.Turns) == 0 {
		t.Fatal("expected at least one turn")
	}
	// Log = intro + one line per turn + winner line.
	if len(r.Turns) != len(r.Log)-2 {
		t.Fatalf("turns=%d, log=%d", len(r.Turns), len(r.Log))
	}
	last := r.Turns[len(r.Turns)-1]
	if last.DefenderHP != 0 {
		t.Fatalf("last DefenderHP = %d, want 0", last.DefenderHP)
	}
	if last.Defender != r.Loser {
		t.Fatalf("last turn defender %q != loser %q", last.Defender, r.Loser)
	}
	for i, tn := range r.Turns {
		if tn.DefenderHP < 0 {
			t.Fatalf("turn %d DefenderHP negative: %d", i, tn.DefenderHP)
		}
	}
}
