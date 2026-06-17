package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

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

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
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
