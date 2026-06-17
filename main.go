package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/chzyer/readline"
	"github.com/strangenoob/pokedexcli/internal/pokeapi"
	"github.com/strangenoob/pokedexcli/internal/pokecache"
	"github.com/strangenoob/pokedexcli/internal/pokedex"
)

// osExit is indirected so commandExit stays testable.
var osExit = os.Exit

func main() {
	savePath, err := pokedex.DefaultSavePath()
	if err != nil {
		fmt.Println("Warning: could not determine save path:", err)
		savePath = "pokedex_save.json"
	}

	dex, err := pokedex.Load(savePath)
	if err != nil {
		fmt.Println("Warning: could not load save, starting fresh:", err)
		dex = pokedex.New()
	}

	cfg := &config{
		client:   pokeapi.NewClient(pokecache.NewCache(5 * time.Second)),
		dex:      dex,
		savePath: savePath,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	rl, err := readline.New("Pokedex > ")
	if err != nil {
		fmt.Println("Error: could not start readline:", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue // Ctrl+C clears the line; keep going
		}
		if err != nil { // io.EOF (Ctrl+D) or other read error
			fmt.Println("Closing the Pokedex... Goodbye!")
			autoSave(cfg)
			return
		}

		words := cleanInput(line)
		if len(words) == 0 {
			continue
		}
		command, exists := getCommands()[words[0]]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}
		if err := command.callback(cfg, words[1:]); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
