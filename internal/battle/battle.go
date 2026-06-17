package battle

import (
	"fmt"
	"math/rand"
)

type Combatant struct {
	Name    string
	HP      int
	Attack  int
	Defense int
	Speed   int
	Types   []string
}

type Result struct {
	Winner string
	Loser  string
	Log    []string
}

func damage(attacker, defender Combatant, rng *rand.Rand) int {
	base := attacker.Attack - defender.Defense/2
	if base < 1 {
		base = 1
	}
	mult := 1.0
	if len(attacker.Types) > 0 {
		mult = Multiplier(attacker.Types[0], defender.Types)
	}
	roll := 85 + rng.Intn(16) // 85..100
	dmg := int(float64(base) * mult * float64(roll) / 100.0)
	if dmg < 1 && mult > 0 {
		dmg = 1
	}
	return dmg
}

// Simulate runs an alternating-turn battle until one combatant faints.
func Simulate(a, b Combatant, rng *rand.Rand) Result {
	first, second := a, b
	if b.Speed > a.Speed || (b.Speed == a.Speed && b.Name < a.Name) {
		first, second = b, a
	}

	log := []string{fmt.Sprintf("%s (HP %d) vs %s (HP %d)!", first.Name, first.HP, second.Name, second.HP)}

	attacker, defender := &first, &second
	for first.HP > 0 && second.HP > 0 {
		dmg := damage(*attacker, *defender, rng)
		defender.HP -= dmg
		if defender.HP < 0 {
			defender.HP = 0
		}
		log = append(log, fmt.Sprintf("%s hits %s for %d (HP %d left)",
			attacker.Name, defender.Name, dmg, defender.HP))
		attacker, defender = defender, attacker
	}

	winner, loser := first, second
	if first.HP <= 0 {
		winner, loser = second, first
	}
	log = append(log, fmt.Sprintf("%s wins!", winner.Name))
	return Result{Winner: winner.Name, Loser: loser.Name, Log: log}
}
