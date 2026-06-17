package pokedex

import "testing"

func TestPartyRules(t *testing.T) {
	dex := New()

	if err := dex.AddToParty("bulbasaur"); err == nil {
		t.Fatal("expected error adding uncaught pokemon")
	}

	dex.Add(sampleBase()) // caught bulbasaur
	if err := dex.AddToParty("bulbasaur"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dex.InParty("bulbasaur") {
		t.Fatal("expected bulbasaur in party")
	}
	if err := dex.AddToParty("bulbasaur"); err == nil {
		t.Fatal("expected error adding duplicate")
	}

	if err := dex.RemoveFromParty("bulbasaur"); err != nil {
		t.Fatalf("unexpected remove error: %v", err)
	}
	if dex.InParty("bulbasaur") {
		t.Fatal("expected bulbasaur removed")
	}
	if err := dex.RemoveFromParty("bulbasaur"); err == nil {
		t.Fatal("expected error removing absent pokemon")
	}
}

func TestPartyFull(t *testing.T) {
	dex := New()
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for _, n := range names {
		base := sampleBase()
		base.Name = n
		dex.Add(base)
	}
	for i := 0; i < MaxPartySize; i++ {
		if err := dex.AddToParty(names[i]); err != nil {
			t.Fatalf("unexpected error at %d: %v", i, err)
		}
	}
	if err := dex.AddToParty(names[MaxPartySize]); err == nil {
		t.Fatal("expected party-full error")
	}
}
