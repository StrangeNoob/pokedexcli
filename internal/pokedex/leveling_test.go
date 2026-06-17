package pokedex

import (
	"testing"

	"github.com/strangenoob/pokedexcli/internal/pokeapi"
)

func sampleBase() pokeapi.Pokemon {
	p := pokeapi.Pokemon{Name: "bulbasaur"}
	p.Stats = []pokeapi.Stat{
		{BaseStat: 45}, {BaseStat: 49}, {BaseStat: 49}, {BaseStat: 45},
	}
	p.Stats[0].Stat.Name = "hp"
	p.Stats[1].Stat.Name = "attack"
	p.Stats[2].Stat.Name = "defense"
	p.Stats[3].Stat.Name = "speed"
	p.Types = []pokeapi.TypeSlot{{Slot: 1}, {Slot: 2}}
	p.Types[0].Type.Name = "grass"
	p.Types[1].Type.Name = "poison"
	return p
}

func TestXPForLevel(t *testing.T) {
	if XPForLevel(5) != 125 {
		t.Fatalf("XPForLevel(5) = %d, want 125", XPForLevel(5))
	}
	if XPForLevel(6) != 216 {
		t.Fatalf("XPForLevel(6) = %d, want 216", XPForLevel(6))
	}
}

func TestAddXPLevelsUp(t *testing.T) {
	dex := New()
	cp := dex.Add(sampleBase()) // level 5, XP 125

	if got := cp.AddXP(90); got != 0 || cp.Level != 5 {
		t.Fatalf("AddXP(90): gained=%d level=%d, want 0 and 5", got, cp.Level)
	}
	if got := cp.AddXP(1); got != 1 || cp.Level != 6 { // XP now 216
		t.Fatalf("AddXP(1): gained=%d level=%d, want 1 and 6", got, cp.Level)
	}
}

func TestScaledStats(t *testing.T) {
	dex := New()
	cp := dex.Add(sampleBase()) // level 5
	if cp.HP() != 45+10+10 {
		t.Fatalf("HP() = %d, want 65", cp.HP())
	}
	if cp.Attack() != 54 {
		t.Fatalf("Attack() = %d, want 54", cp.Attack())
	}
	if got := cp.TypeNames(); len(got) != 2 || got[0] != "grass" || got[1] != "poison" {
		t.Fatalf("TypeNames() = %v", got)
	}
}
