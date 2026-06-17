package pokedex

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dex := New()
	cp := dex.Add(sampleBase())
	cp.AddXP(200)
	if err := dex.AddToParty("bulbasaur"); err != nil {
		t.Fatalf("party add: %v", err)
	}

	path := filepath.Join(t.TempDir(), "save.json")
	if err := dex.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	got, ok := loaded.Get("bulbasaur")
	if !ok {
		t.Fatal("bulbasaur missing after load")
	}
	if got.Level != cp.Level || got.XP != cp.XP {
		t.Fatalf("level/xp mismatch: got %d/%d want %d/%d", got.Level, got.XP, cp.Level, cp.XP)
	}
	if !loaded.InParty("bulbasaur") {
		t.Fatal("party not persisted")
	}
}

func TestLoadMissingFileIsFreshStart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")
	dex, err := Load(path)
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if dex == nil || dex.Caught == nil {
		t.Fatal("expected initialized empty pokedex")
	}
}
