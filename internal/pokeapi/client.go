package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/strangenoob/pokedexcli/internal/pokecache"
)

const defaultBaseURL = "https://pokeapi.co/api/v2"

type Client struct {
	baseURL    string
	httpClient *http.Client
	cache      *pokecache.Cache
}

func NewClient(cache *pokecache.Cache) *Client {
	return &Client{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      cache,
	}
}

func (c *Client) fetch(url string) ([]byte, error) {
	if val, ok := c.cache.Get(url); ok {
		return val, nil
	}

	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("request to %s failed: %s", url, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	c.cache.Add(url, body)
	return body, nil
}

func (c *Client) FetchLocationAreas(pageURL *string) (LocationAreas, error) {
	url := c.baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	body, err := c.fetch(url)
	if err != nil {
		return LocationAreas{}, err
	}

	var out LocationAreas
	if err := json.Unmarshal(body, &out); err != nil {
		return LocationAreas{}, err
	}
	return out, nil
}

func (c *Client) FetchLocationArea(name string) (LocationArea, error) {
	url := fmt.Sprintf("%s/location-area/%s", c.baseURL, name)

	body, err := c.fetch(url)
	if err != nil {
		return LocationArea{}, err
	}

	var out LocationArea
	if err := json.Unmarshal(body, &out); err != nil {
		return LocationArea{}, err
	}
	return out, nil
}

func (c *Client) FetchPokemon(name string) (Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, name)

	body, err := c.fetch(url)
	if err != nil {
		return Pokemon{}, err
	}

	var out Pokemon
	if err := json.Unmarshal(body, &out); err != nil {
		return Pokemon{}, err
	}
	return out, nil
}
