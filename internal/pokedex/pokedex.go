package pokedex

import (
	"fmt"

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
}

func New() *Pokedex {
	return &Pokedex{Caught: make(map[string]*CaughtPokemon)}
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
