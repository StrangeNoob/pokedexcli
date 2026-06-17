package battle

// matchups lists, per attacking type, the defending types it is strong/weak/
// useless against. Anything not listed is 1x.
type matchups struct {
	double []string // 2x
	half   []string // 0.5x
	zero   []string // 0x
}

var chart = map[string]matchups{
	"normal":   {half: []string{"rock", "steel"}, zero: []string{"ghost"}},
	"fire":     {double: []string{"grass", "ice", "bug", "steel"}, half: []string{"fire", "water", "rock", "dragon"}},
	"water":    {double: []string{"fire", "ground", "rock"}, half: []string{"water", "grass", "dragon"}},
	"electric": {double: []string{"water", "flying"}, half: []string{"electric", "grass", "dragon"}, zero: []string{"ground"}},
	"grass":    {double: []string{"water", "ground", "rock"}, half: []string{"fire", "grass", "poison", "flying", "bug", "dragon", "steel"}},
	"ice":      {double: []string{"grass", "ground", "flying", "dragon"}, half: []string{"fire", "water", "ice", "steel"}},
	"fighting": {double: []string{"normal", "ice", "rock", "dark", "steel"}, half: []string{"poison", "flying", "psychic", "bug", "fairy"}, zero: []string{"ghost"}},
	"poison":   {double: []string{"grass", "fairy"}, half: []string{"poison", "ground", "rock", "ghost"}, zero: []string{"steel"}},
	"ground":   {double: []string{"fire", "electric", "poison", "rock", "steel"}, half: []string{"grass", "bug"}, zero: []string{"flying"}},
	"flying":   {double: []string{"grass", "fighting", "bug"}, half: []string{"electric", "rock", "steel"}},
	"psychic":  {double: []string{"fighting", "poison"}, half: []string{"psychic", "steel"}, zero: []string{"dark"}},
	"bug":      {double: []string{"grass", "psychic", "dark"}, half: []string{"fire", "fighting", "poison", "flying", "ghost", "steel", "fairy"}},
	"rock":     {double: []string{"fire", "ice", "flying", "bug"}, half: []string{"fighting", "ground", "steel"}},
	"ghost":    {double: []string{"psychic", "ghost"}, half: []string{"dark"}, zero: []string{"normal"}},
	"dragon":   {double: []string{"dragon"}, half: []string{"steel"}, zero: []string{"fairy"}},
	"dark":     {double: []string{"psychic", "ghost"}, half: []string{"fighting", "dark", "fairy"}},
	"steel":    {double: []string{"ice", "rock", "fairy"}, half: []string{"fire", "water", "electric", "steel"}},
	"fairy":    {double: []string{"fighting", "dragon", "dark"}, half: []string{"fire", "poison", "steel"}},
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// Multiplier returns the combined effectiveness of attackingType against a
// (possibly dual-typed) defender. Unknown types contribute 1x.
func Multiplier(attackingType string, defenderTypes []string) float64 {
	m, ok := chart[attackingType]
	if !ok {
		return 1.0
	}
	total := 1.0
	for _, d := range defenderTypes {
		switch {
		case contains(m.zero, d):
			total *= 0.0
		case contains(m.double, d):
			total *= 2.0
		case contains(m.half, d):
			total *= 0.5
		}
	}
	return total
}
