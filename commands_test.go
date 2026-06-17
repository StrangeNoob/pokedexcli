package main

import (
	"math/rand"
	"testing"
)

func TestParseCatchArgs(t *testing.T) {
	cases := []struct {
		args     []string
		wild     string
		wantName string
		wantBall string
		wantErr  bool
	}{
		{nil, "pikachu", "pikachu", "pokeball", false},
		{[]string{"greatball"}, "pikachu", "pikachu", "greatball", false},
		{[]string{"bulbasaur"}, "", "bulbasaur", "pokeball", false},
		{[]string{"bulbasaur", "ultraball"}, "", "bulbasaur", "ultraball", false},
		{nil, "", "", "", true},                               // no target
		{[]string{"greatball"}, "", "", "", true},             // ball, no target
		{[]string{"pikachu", "masterball"}, "", "", "", true}, // bad ball
	}
	for i, c := range cases {
		name, ballName, err := parseCatchArgs(c.args, c.wild)
		if c.wantErr {
			if err == nil {
				t.Errorf("case %d: expected error, got (%q,%q)", i, name, ballName)
			}
			continue
		}
		if err != nil {
			t.Errorf("case %d: unexpected error: %v", i, err)
			continue
		}
		if name != c.wantName || ballName != c.wantBall {
			t.Errorf("case %d: got (%q,%q), want (%q,%q)", i, name, ballName, c.wantName, c.wantBall)
		}
	}
}

func TestRandomChoice(t *testing.T) {
	if got := randomChoice(nil, rand.New(rand.NewSource(1))); got != "" {
		t.Fatalf("empty slice should return \"\", got %q", got)
	}
	names := []string{"a", "b", "c"}
	r1 := randomChoice(names, rand.New(rand.NewSource(42)))
	r2 := randomChoice(names, rand.New(rand.NewSource(42)))
	if r1 != r2 {
		t.Fatalf("same seed must be deterministic: %q vs %q", r1, r2)
	}
	in := false
	for _, n := range names {
		if n == r1 {
			in = true
		}
	}
	if !in {
		t.Fatalf("result %q not from input", r1)
	}
}
