package pokedex

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStartingBalls(t *testing.T) {
	p := New()
	if p.BallCount("pokeball") != 20 || p.BallCount("greatball") != 10 || p.BallCount("ultraball") != 5 {
		t.Fatalf("starting stash = %v", p.Balls)
	}
}

func TestUseBall(t *testing.T) {
	p := New()
	if err := p.UseBall("pokeball"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.BallCount("pokeball") != 19 {
		t.Fatalf("count after use = %d, want 19", p.BallCount("pokeball"))
	}
	if err := p.UseBall("masterball"); err == nil {
		t.Fatal("expected error for unknown ball type")
	}
	p.Balls["ultraball"] = 0
	if err := p.UseBall("ultraball"); err == nil {
		t.Fatal("expected error when out of balls")
	}
}

func TestLoadBackfillsBalls(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.json")
	if err := os.WriteFile(path, []byte(`{"caught":{},"party":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	dex, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if dex.BallCount("pokeball") != 20 {
		t.Fatalf("expected backfilled stash, got %d pokeballs", dex.BallCount("pokeball"))
	}
}

func TestSaveLoadPreservesBalls(t *testing.T) {
	p := New()
	if err := p.UseBall("pokeball"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "save.json")
	if err := p.Save(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.BallCount("pokeball") != 19 {
		t.Fatalf("pokeball after round-trip = %d, want 19", loaded.BallCount("pokeball"))
	}
}
