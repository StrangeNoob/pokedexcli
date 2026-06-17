package pokedex

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func DefaultSavePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pokedexcli", "save.json"), nil
}

func (p *Pokedex) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func Load(path string) (*Pokedex, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return New(), nil
	}
	if err != nil {
		return nil, err
	}

	var p Pokedex
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if p.Caught == nil {
		p.Caught = make(map[string]*CaughtPokemon)
	}
	if p.Balls == nil {
		p.Balls = StartingBalls()
	}
	return &p, nil
}
