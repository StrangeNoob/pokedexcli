package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/strangenoob/pokedexcli/internal/pokecache"
)

func TestFetchPokemonUsesCacheAndParsesJSON(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = w.Write([]byte(`{"name":"pikachu","base_experience":112,"height":4,"weight":60,
			"stats":[{"base_stat":35,"stat":{"name":"hp"}}],
			"types":[{"slot":1,"type":{"name":"electric"}}]}`))
	}))
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		httpClient: srv.Client(),
		cache:      pokecache.NewCache(5 * time.Second),
	}

	got, err := c.FetchPokemon("pikachu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "pikachu" || got.BaseExperience != 112 {
		t.Fatalf("bad parse: %+v", got)
	}
	if got.Types[0].Type.Name != "electric" {
		t.Fatalf("bad types parse: %+v", got.Types)
	}

	// Second call must hit the cache, not the server.
	if _, err := c.FetchPokemon("pikachu"); err != nil {
		t.Fatalf("unexpected error on cached call: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 server call, got %d", calls)
	}
}
