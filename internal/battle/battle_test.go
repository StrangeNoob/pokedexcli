package battle

import (
	"math/rand"
	"testing"
)

func TestSimulateIsDeterministicForSameSeed(t *testing.T) {
	a := Combatant{Name: "pikachu", HP: 60, Attack: 55, Defense: 40, Speed: 90, Types: []string{"electric"}}
	b := Combatant{Name: "geodude", HP: 80, Attack: 80, Defense: 100, Speed: 20, Types: []string{"rock", "ground"}}

	r1 := Simulate(a, b, rand.New(rand.NewSource(42)))
	r2 := Simulate(a, b, rand.New(rand.NewSource(42)))

	if r1.Winner != r2.Winner || len(r1.Log) != len(r2.Log) {
		t.Fatalf("non-deterministic: %q/%d vs %q/%d", r1.Winner, len(r1.Log), r2.Winner, len(r2.Log))
	}
}

func TestStrongerPokemonWins(t *testing.T) {
	strong := Combatant{Name: "dragonite", HP: 150, Attack: 134, Defense: 95, Speed: 80, Types: []string{"dragon"}}
	weak := Combatant{Name: "magikarp", HP: 30, Attack: 10, Defense: 55, Speed: 80, Types: []string{"water"}}

	r := Simulate(strong, weak, rand.New(rand.NewSource(1)))
	if r.Winner != "dragonite" {
		t.Fatalf("expected dragonite to win, got %q", r.Winner)
	}
}

func TestFasterPokemonStrikesFirst(t *testing.T) {
	fast := Combatant{Name: "fast", HP: 50, Attack: 50, Defense: 10, Speed: 200, Types: []string{"normal"}}
	slow := Combatant{Name: "slow", HP: 50, Attack: 50, Defense: 10, Speed: 1, Types: []string{"normal"}}

	r := Simulate(slow, fast, rand.New(rand.NewSource(7)))
	// First log line after the intro names the first attacker.
	if r.Log[1][:4] != "fast" {
		t.Fatalf("expected 'fast' to strike first, got log line %q", r.Log[1])
	}
}
