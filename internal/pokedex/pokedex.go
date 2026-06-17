package pokedex

import (
	"fmt"

	"github.com/strangenoob/pokedexcli/internal/ball"
	"github.com/strangenoob/pokedexcli/internal/pokeapi"
)

const (
	MaxPartySize  = 6
	StartingLevel = 5
)

type CaughtPokemon struct {
	Base  pokeapi.Pokemon `json:"base"`
	Level int             `json:"level"`
	XP    int             `json:"xp"`
}

type Pokedex struct {
	Caught map[string]*CaughtPokemon `json:"caught"`
	Party  []string                  `json:"party"`
	Balls  map[string]int            `json:"balls"`
}

// StartingBalls returns a fresh starting ball inventory.
func StartingBalls() map[string]int {
	return map[string]int{"pokeball": 20, "greatball": 10, "ultraball": 5}
}

func New() *Pokedex {
	return &Pokedex{
		Caught: make(map[string]*CaughtPokemon),
		Balls:  StartingBalls(),
	}
}

func (p *Pokedex) BallCount(name string) int {
	return p.Balls[name]
}

// UseBall consumes one ball of the given type.
func (p *Pokedex) UseBall(name string) error {
	if !ball.IsValid(name) {
		return fmt.Errorf("unknown ball type %q", name)
	}
	if p.Balls[name] <= 0 {
		return fmt.Errorf("you have no %s left", name)
	}
	p.Balls[name]--
	return nil
}

func (p *Pokedex) Add(base pokeapi.Pokemon) *CaughtPokemon {
	cp := &CaughtPokemon{
		Base:  base,
		Level: StartingLevel,
		XP:    XPForLevel(StartingLevel),
	}
	p.Caught[base.Name] = cp
	return cp
}

func (p *Pokedex) Get(name string) (*CaughtPokemon, bool) {
	cp, ok := p.Caught[name]
	return cp, ok
}

func (p *Pokedex) InParty(name string) bool {
	for _, n := range p.Party {
		if n == name {
			return true
		}
	}
	return false
}

func (p *Pokedex) AddToParty(name string) error {
	if _, ok := p.Caught[name]; !ok {
		return fmt.Errorf("you have not caught %s", name)
	}
	if p.InParty(name) {
		return fmt.Errorf("%s is already in your party", name)
	}
	if len(p.Party) >= MaxPartySize {
		return fmt.Errorf("your party is full (max %d)", MaxPartySize)
	}
	p.Party = append(p.Party, name)
	return nil
}

func (p *Pokedex) RemoveFromParty(name string) error {
	for i, n := range p.Party {
		if n == name {
			p.Party = append(p.Party[:i], p.Party[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%s is not in your party", name)
}
