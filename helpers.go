package main

import (
	"fmt"
	"math/rand"

	"github.com/strangenoob/pokedexcli/internal/ball"
)

// parseCatchArgs resolves the target Pokémon name and ball type from catch args,
// falling back to the current wild target and a default pokeball.
func parseCatchArgs(args []string, wildTarget string) (name, ballName string, err error) {
	name = wildTarget
	ballName = "pokeball"

	switch len(args) {
	case 0:
		// keep wild target + default ball
	case 1:
		if ball.IsValid(args[0]) {
			ballName = args[0]
		} else {
			name = args[0]
		}
	default:
		name = args[0]
		ballName = args[1]
	}

	if name == "" {
		return "", "", fmt.Errorf("nothing to catch — try 'encounter' or 'catch <name>'")
	}
	if !ball.IsValid(ballName) {
		return "", "", fmt.Errorf("unknown ball type %q (try: pokeball, greatball, ultraball)", ballName)
	}
	return name, ballName, nil
}

// randomChoice returns a random element of names, or "" if names is empty.
func randomChoice(names []string, rng *rand.Rand) string {
	if len(names) == 0 {
		return ""
	}
	return names[rng.Intn(len(names))]
}
