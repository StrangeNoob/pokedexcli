package tui

import (
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/strangenoob/pokedexcli/internal/pokeapi"
	"github.com/strangenoob/pokedexcli/internal/pokedex"
)

type Deps struct {
	Dex      *pokedex.Pokedex
	Client   *pokeapi.Client
	RNG      *rand.Rand
	SavePath string
	Art      *ArtStore
}

func Run(deps Deps) error {
	p := tea.NewProgram(newRootModel(deps), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
