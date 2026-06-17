package battle

import "testing"

func TestMultiplier(t *testing.T) {
	cases := []struct {
		atk  string
		def  []string
		want float64
	}{
		{"water", []string{"fire"}, 2.0},
		{"fire", []string{"water"}, 0.5},
		{"electric", []string{"ground"}, 0.0},
		{"normal", []string{"ghost"}, 0.0},
		{"grass", []string{"water", "ground"}, 4.0}, // 2x * 2x
		{"fire", []string{"grass", "dragon"}, 1.0},  // 2x * 0.5x
		{"water", []string{"normal"}, 1.0},          // neutral
		{"mystery", []string{"fire"}, 1.0},          // unknown attacker
	}
	for _, c := range cases {
		if got := Multiplier(c.atk, c.def); got != c.want {
			t.Errorf("Multiplier(%q, %v) = %v, want %v", c.atk, c.def, got, c.want)
		}
	}
}
