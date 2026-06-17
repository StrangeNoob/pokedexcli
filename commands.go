package main

import (
	"fmt"
	"math/rand"

	"github.com/strangenoob/pokedexcli/internal/ball"
	"github.com/strangenoob/pokedexcli/internal/battle"
	"github.com/strangenoob/pokedexcli/internal/pokeapi"
	"github.com/strangenoob/pokedexcli/internal/pokedex"
)

type config struct {
	client      *pokeapi.Client
	dex         *pokedex.Pokedex
	savePath    string
	rng         *rand.Rand
	nextLocURL  *string
	prevLocURL  *string
	areaPokemon []string
	wildTarget  string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit":      {"exit", "Exit the Pokedex", commandExit},
		"help":      {"help", "Display a help message", commandHelp},
		"map":       {"map", "Display the next 20 location areas", commandMap},
		"mapb":      {"mapb", "Display the previous 20 location areas", commandMapb},
		"explore":   {"explore", "Explore a location area and list the Pokemon there", commandExplore},
		"catch":     {"catch", "Catch a Pokemon: catch [name] [ball]", commandCatch},
		"inspect":   {"inspect", "Inspect a Pokemon you have caught", commandInspect},
		"pokedex":   {"pokedex", "List all caught Pokemon", commandPokedex},
		"party":     {"party", "Show party, or 'party add|remove <name>'", commandParty},
		"battle":    {"battle", "Battle two caught Pokemon: battle <a> <b>", commandBattle},
		"encounter": {"encounter", "Find a random wild Pokemon in the explored area", commandEncounter},
		"bag":       {"bag", "Show your Pokeball inventory", commandBag},
		"save":      {"save", "Save your progress to disk", commandSave},
	}
}

func autoSave(cfg *config) {
	if err := cfg.dex.Save(cfg.savePath); err != nil {
		fmt.Println("Warning: could not save progress:", err)
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	autoSave(cfg)
	osExit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	areas, err := cfg.client.FetchLocationAreas(cfg.nextLocURL)
	if err != nil {
		return err
	}
	cfg.nextLocURL = areas.Next
	cfg.prevLocURL = areas.Previous
	for _, a := range areas.Results {
		fmt.Println(a.Name)
	}
	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.prevLocURL == nil {
		return fmt.Errorf("you're on the first page")
	}
	areas, err := cfg.client.FetchLocationAreas(cfg.prevLocURL)
	if err != nil {
		return err
	}
	cfg.nextLocURL = areas.Next
	cfg.prevLocURL = areas.Previous
	for _, a := range areas.Results {
		fmt.Println(a.Name)
	}
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location area name")
	}
	area, err := cfg.client.FetchLocationArea(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	names := make([]string, 0, len(area.PokemonEncounters))
	for _, e := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", e.Pokemon.Name)
		names = append(names, e.Pokemon.Name)
	}
	cfg.areaPokemon = names
	return nil
}

func commandCatch(cfg *config, args []string) error {
	name, ballName, err := parseCatchArgs(args, cfg.wildTarget)
	if err != nil {
		return err
	}
	if cfg.dex.BallCount(ballName) <= 0 {
		return fmt.Errorf("you have no %s left", ballName)
	}

	pokemon, err := cfg.client.FetchPokemon(name)
	if err != nil {
		return err
	}

	if err := cfg.dex.UseBall(ballName); err != nil {
		return err
	}
	fmt.Printf("Throwing a %s at %s...\n", ballName, name)

	base := 100 - (pokemon.BaseExperience / 4)
	if base < 5 {
		base = 5
	}
	chance := int(float64(base) * ball.Multiplier(ballName))
	if chance > 100 {
		chance = 100
	}

	if cfg.rng.Intn(100) < chance {
		cfg.dex.Add(pokemon)
		fmt.Printf("%s was caught!\n", pokemon.Name)
		if name == cfg.wildTarget {
			cfg.wildTarget = ""
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	autoSave(cfg)
	return nil
}

func commandEncounter(cfg *config, args []string) error {
	if len(cfg.areaPokemon) == 0 {
		return fmt.Errorf("explore an area first")
	}
	cfg.wildTarget = randomChoice(cfg.areaPokemon, cfg.rng)
	fmt.Printf("A wild %s appeared!\n", cfg.wildTarget)
	fmt.Println("Use 'catch' to throw a ball at it.")
	return nil
}

func commandBag(cfg *config, args []string) error {
	fmt.Println("Your bag:")
	for _, name := range ball.Names() {
		fmt.Printf(" - %s: %d\n", name, cfg.dex.BallCount(name))
	}
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	cp, ok := cfg.dex.Get(args[0])
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", cp.Base.Name)
	fmt.Printf("Level: %d  XP: %d\n", cp.Level, cp.XP)
	fmt.Printf("Height: %d\n", cp.Base.Height)
	fmt.Printf("Weight: %d\n", cp.Base.Weight)
	fmt.Println("Stats:")
	for _, s := range cp.Base.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, tn := range cp.TypeNames() {
		fmt.Printf("  - %s\n", tn)
	}
	return nil
}

func commandPokedex(cfg *config, args []string) error {
	fmt.Println("Your Pokedex:")
	for _, cp := range cfg.dex.Caught {
		fmt.Printf(" - %s (Lvl %d)\n", cp.Base.Name, cp.Level)
	}
	return nil
}

func commandParty(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Your party:")
		for _, name := range cfg.dex.Party {
			cp, _ := cfg.dex.Get(name)
			fmt.Printf(" - %s (Lvl %d)\n", name, cp.Level)
		}
		return nil
	}
	if len(args) < 2 {
		return fmt.Errorf("usage: party add|remove <name>")
	}
	switch args[0] {
	case "add":
		if err := cfg.dex.AddToParty(args[1]); err != nil {
			return err
		}
		fmt.Printf("%s joined your party\n", args[1])
	case "remove":
		if err := cfg.dex.RemoveFromParty(args[1]); err != nil {
			return err
		}
		fmt.Printf("%s left your party\n", args[1])
	default:
		return fmt.Errorf("unknown party action %q (use add or remove)", args[0])
	}
	autoSave(cfg)
	return nil
}

func toCombatant(cp *pokedex.CaughtPokemon) battle.Combatant {
	return battle.Combatant{
		Name:    cp.Base.Name,
		HP:      cp.HP(),
		Attack:  cp.Attack(),
		Defense: cp.Defense(),
		Speed:   cp.Speed(),
		Types:   cp.TypeNames(),
	}
}

func commandBattle(cfg *config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: battle <yourPokemon> <opponentPokemon>")
	}
	a, ok := cfg.dex.Get(args[0])
	if !ok {
		return fmt.Errorf("you have not caught %s", args[0])
	}
	b, ok := cfg.dex.Get(args[1])
	if !ok {
		return fmt.Errorf("you have not caught %s", args[1])
	}

	result := battle.Simulate(toCombatant(a), toCombatant(b), cfg.rng)
	for _, line := range result.Log {
		fmt.Println(line)
	}

	winner, _ := cfg.dex.Get(result.Winner)
	loser, _ := cfg.dex.Get(result.Loser)
	xp := loser.Level * 10
	if levels := winner.AddXP(xp); levels > 0 {
		fmt.Printf("%s gained %d XP and reached level %d!\n", winner.Base.Name, xp, winner.Level)
	} else {
		fmt.Printf("%s gained %d XP.\n", winner.Base.Name, xp)
	}
	autoSave(cfg)
	return nil
}

func commandSave(cfg *config, args []string) error {
	if err := cfg.dex.Save(cfg.savePath); err != nil {
		return err
	}
	fmt.Println("Progress saved.")
	return nil
}
