package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/strangenoob/pokedexcli/internal/pokecache"
)

func TestFetchImageAndSpriteField(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pokemon/pikachu" {
			w.Write([]byte(`{"name":"pikachu","sprites":{"front_default":"http://example/img.png"}}`))
			return
		}
		w.Write([]byte("PNGBYTES"))
	}))
	defer srv.Close()

	c := &Client{baseURL: srv.URL, httpClient: srv.Client(), cache: pokecache.NewCache(5 * time.Second)}

	p, err := c.FetchPokemon("pikachu")
	if err != nil {
		t.Fatal(err)
	}
	if p.Sprites.FrontDefault != "http://example/img.png" {
		t.Fatalf("sprite url = %q", p.Sprites.FrontDefault)
	}

	data, err := c.FetchImage(srv.URL + "/img.png")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "PNGBYTES" {
		t.Fatalf("image bytes = %q", data)
	}
}
