package ball

var multipliers = map[string]float64{
	"pokeball":  1.0,
	"greatball": 1.5,
	"ultraball": 2.0,
}

var names = []string{"pokeball", "greatball", "ultraball"}

// Multiplier returns the catch-rate multiplier for a ball type, or 0 if unknown.
func Multiplier(name string) float64 {
	return multipliers[name]
}

func IsValid(name string) bool {
	_, ok := multipliers[name]
	return ok
}

// Names returns the ball types in stable display order (a fresh copy).
func Names() []string {
	out := make([]string, len(names))
	copy(out, names)
	return out
}
