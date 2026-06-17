package ball

import "testing"

func TestMultiplier(t *testing.T) {
	cases := map[string]float64{
		"pokeball": 1.0, "greatball": 1.5, "ultraball": 2.0, "masterball": 0,
	}
	for n, want := range cases {
		if got := Multiplier(n); got != want {
			t.Errorf("Multiplier(%q) = %v, want %v", n, got, want)
		}
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid("pokeball") || !IsValid("ultraball") {
		t.Fatal("expected pokeball/ultraball valid")
	}
	if IsValid("masterball") || IsValid("") {
		t.Fatal("expected invalid for unknown/empty")
	}
}

func TestNames(t *testing.T) {
	want := []string{"pokeball", "greatball", "ultraball"}
	got := Names()
	if len(got) != len(want) {
		t.Fatalf("Names() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Names()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
	got[0] = "mutated"
	if Names()[0] != "pokeball" {
		t.Fatal("Names() must return a defensive copy")
	}
}
