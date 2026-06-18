package tui

import (
	"testing"

	"github.com/strangenoob/pokedexcli/internal/pokeapi"
)

func TestPokedexPartyToggle(t *testing.T) {
	d := testDeps()
	d.Dex.Add(pokeapi.Pokemon{Name: "pikachu"})
	m := newPokedexModel(d)

	if d.Dex.InParty("pikachu") {
		t.Fatal("should not start in party")
	}
	updated, _ := m.Update(runeKey('p'))
	m = updated.(pokedexModel)
	if !d.Dex.InParty("pikachu") {
		t.Fatal("pressing p should add to party")
	}
	updated, _ = m.Update(runeKey('p'))
	if d.Dex.InParty("pikachu") {
		t.Fatal("pressing p again should remove from party")
	}
}
