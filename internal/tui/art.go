package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/strangenoob/pokedexcli/internal/sprite"
)

const artWidth = 32

type artLoadedMsg struct {
	name string
	art  string
	err  error
}

// ArtStore caches rendered sprite art by Pokémon name. All map access happens on
// the Bubble Tea update thread; only fetching/rendering runs in a command goroutine.
type ArtStore struct {
	rendered map[string]string
	pending  map[string]bool
}

func NewArtStore() *ArtStore {
	return &ArtStore{rendered: map[string]string{}, pending: map[string]bool{}}
}

func (s *ArtStore) get(name string) string { return s.rendered[name] }

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
	if msg.err == nil && msg.art != "" {
		s.rendered[msg.name] = msg.art
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
		return artLoadedMsg{name: name, art: art, err: err}
	}
}
