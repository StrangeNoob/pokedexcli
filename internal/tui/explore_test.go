package tui

import (
	"strings"
	"testing"

	"github.com/strangenoob/pokedexcli/internal/pokeapi"
)

func TestExploreCachesWildStats(t *testing.T) {
	d := testDeps()
	m := newExploreModel(d)
	_, _ = m.Update(artLoadedMsg{
		name:    "magikarp",
		art:     "ART",
		pokemon: pokeapi.Pokemon{Name: "magikarp", Height: 9},
	})
	p, ok := d.Art.poke("magikarp")
	if !ok || p.Height != 9 {
		t.Fatalf("wild stats not cached: %+v ok=%v", p, ok)
	}
}

func TestWildStatsViewShowsStats(t *testing.T) {
	p := pokeapi.Pokemon{Name: "magikarp", Height: 9, Weight: 100}
	p.Stats = []pokeapi.Stat{{BaseStat: 20}}
	p.Stats[0].Stat.Name = "hp"
	p.Types = []pokeapi.TypeSlot{{}}
	p.Types[0].Type.Name = "water"

	out := wildStatsView(p)
	for _, want := range []string{"magikarp", "water", "hp", "20", "Height 9"} {
		if !strings.Contains(out, want) {
			t.Errorf("stats view missing %q in %q", want, out)
		}
	}
}

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

func TestExploreStoresArt(t *testing.T) {
	d := testDeps()
	m := newExploreModel(d)
	_, _ = m.Update(artLoadedMsg{name: "magikarp", art: "ART"})
	if d.Art.get("magikarp") != "ART" {
		t.Fatal("explore should store art via the shared store")
	}
}

var errTest = stringError("boom")

type stringError string

func (e stringError) Error() string { return string(e) }
