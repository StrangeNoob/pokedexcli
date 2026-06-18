package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/strangenoob/pokedexcli/internal/pokeapi"
	"github.com/strangenoob/pokedexcli/internal/sprite"
)

const artWidth = 32

type artLoadedMsg struct {
	name    string
	art     string
	pokemon pokeapi.Pokemon
	err     error
}

// ArtStore caches rendered sprite art and the fetched Pokémon (for stats) by name.
// All map access happens on the Bubble Tea update thread; only fetching/rendering
// runs in a command goroutine.
type ArtStore struct {
	rendered map[string]string
	pokes    map[string]pokeapi.Pokemon
	pending  map[string]bool
}

func NewArtStore() *ArtStore {
	return &ArtStore{
		rendered: map[string]string{},
		pokes:    map[string]pokeapi.Pokemon{},
		pending:  map[string]bool{},
	}
}

func (s *ArtStore) get(name string) string { return s.rendered[name] }

// poke returns the fetched Pokémon data for name, if it has been loaded.
func (s *ArtStore) poke(name string) (pokeapi.Pokemon, bool) {
	p, ok := s.pokes[name]
	return p, ok
}

// request returns a command to load art for name, or nil if it is empty, already
// rendered, or already in flight.
func (s *ArtStore) request(deps Deps, name string) tea.Cmd {
	if name == "" || s.rendered[name] != "" || s.pending[name] {
		return nil
	}
	s.pending[name] = true
	return artCmd(deps, name)
}

func (s *ArtStore) handle(msg artLoadedMsg) {
	s.pending[msg.name] = false
	if msg.err == nil {
		if msg.art != "" {
			s.rendered[msg.name] = msg.art
		}
		s.pokes[msg.name] = msg.pokemon
	}
}

func artCmd(deps Deps, name string) tea.Cmd {
	client := deps.Client
	return func() tea.Msg {
		p, err := client.FetchPokemon(name)
		if err != nil {
			return artLoadedMsg{name: name, err: err}
		}
		url := p.Sprites.FrontDefault
		if url == "" {
			return artLoadedMsg{name: name, err: fmt.Errorf("no sprite for %s", name)}
		}
		data, err := client.FetchImage(url)
		if err != nil {
			return artLoadedMsg{name: name, err: err}
		}
		art, err := sprite.Render(data, artWidth)
		return artLoadedMsg{name: name, art: art, pokemon: p, err: err}
	}
}
