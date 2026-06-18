package tui

import (
	"strings"
	"testing"

	"github.com/strangenoob/pokedexcli/internal/pokeapi"
)

func TestExploreAreasLoaded(t *testing.T) {
	m := newExploreModel(testDeps())
	updated, _ := m.Update(areasLoadedMsg{areas: []string{"area-a", "area-b"}})
	em := updated.(exploreModel)
	if len(em.areas) != 2 || em.loading {
		t.Fatalf("areas=%v loading=%v", em.areas, em.loading)
	}
}

func TestExploreCatchAddsToDex(t *testing.T) {
	d := testDeps()
	m := newExploreModel(d)
	before := d.Dex.BallCount("pokeball")

	// BaseExperience 0 => base 100 => chance 100 => always caught.
	updated, _ := m.Update(pokemonFetchedMsg{
		pokemon: pokeapi.Pokemon{Name: "pikachu"},
		name:    "pikachu",
		ball:    "pokeball",
	})
	em := updated.(exploreModel)

	if _, ok := d.Dex.Get("pikachu"); !ok {
		t.Fatal("pikachu should be caught and added to the dex")
	}
	if d.Dex.BallCount("pokeball") != before-1 {
		t.Fatalf("ball not consumed: %d", d.Dex.BallCount("pokeball"))
	}
	if !strings.Contains(em.status, "caught") {
		t.Fatalf("status = %q", em.status)
	}
}

func TestExploreFetchErrorDoesNotSpendBall(t *testing.T) {
	d := testDeps()
	m := newExploreModel(d)
	before := d.Dex.BallCount("pokeball")
	updated, _ := m.Update(pokemonFetchedMsg{name: "x", ball: "pokeball", err: errTest})
	_ = updated
	if d.Dex.BallCount("pokeball") != before {
		t.Fatal("fetch error must not consume a ball")
	}
}

var errTest = stringError("boom")

type stringError string

func (e stringError) Error() string { return string(e) }
