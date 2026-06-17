package pokedex

// XPForLevel returns the cumulative XP required to be at the given level: level^3.
func XPForLevel(level int) int {
	return level * level * level
}

// AddXP adds xp and levels up while the next threshold is reached.
// Returns the number of levels gained.
func (cp *CaughtPokemon) AddXP(xp int) int {
	cp.XP += xp
	gained := 0
	for cp.XP >= XPForLevel(cp.Level+1) {
		cp.Level++
		gained++
	}
	return gained
}

func (cp *CaughtPokemon) baseStat(name string) int {
	for _, s := range cp.Base.Stats {
		if s.Stat.Name == name {
			return s.BaseStat
		}
	}
	return 0
}

// Level-scaled battle stats.
func (cp *CaughtPokemon) HP() int      { return cp.baseStat("hp") + 2*cp.Level + 10 }
func (cp *CaughtPokemon) Attack() int  { return cp.baseStat("attack") + cp.Level }
func (cp *CaughtPokemon) Defense() int { return cp.baseStat("defense") + cp.Level }
func (cp *CaughtPokemon) Speed() int   { return cp.baseStat("speed") + cp.Level }

func (cp *CaughtPokemon) TypeNames() []string {
	names := make([]string, 0, len(cp.Base.Types))
	for _, t := range cp.Base.Types {
		names = append(names, t.Type.Name)
	}
	return names
}
