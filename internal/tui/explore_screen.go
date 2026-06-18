package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/strangenoob/pokedexcli/internal/ball"
	"github.com/strangenoob/pokedexcli/internal/pokeapi"
)

type exploreStep int

const (
	areaListStep exploreStep = iota
	wildListStep
)

type areasLoadedMsg struct {
	areas []string
	next  *string
	prev  *string
	err   error
}

type wildPokemonMsg struct {
	areaName string
	names    []string
	err      error
}

type pokemonFetchedMsg struct {
	pokemon pokeapi.Pokemon
	name    string
	ball    string
	err     error
}

type exploreModel struct {
	deps    Deps
	step    exploreStep
	spinner spinner.Model
	loading bool

	areas   []string
	nextURL *string
	prevURL *string
	areaCur int

	areaName string
	wild     []string
	wildCur  int
	ballIdx  int

	status string
}

func newExploreModel(deps Deps) exploreModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	return exploreModel{deps: deps, step: areaListStep, spinner: sp, loading: true}
}

func (m exploreModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadAreasCmd(nil))
}

func (m exploreModel) loadAreasCmd(page *string) tea.Cmd {
	client := m.deps.Client
	return func() tea.Msg {
		res, err := client.FetchLocationAreas(page)
		if err != nil {
			return areasLoadedMsg{err: err}
		}
		names := make([]string, 0, len(res.Results))
		for _, r := range res.Results {
			names = append(names, r.Name)
		}
		return areasLoadedMsg{areas: names, next: res.Next, prev: res.Previous}
	}
}

func (m exploreModel) loadWildCmd(area string) tea.Cmd {
	client := m.deps.Client
	return func() tea.Msg {
		res, err := client.FetchLocationArea(area)
		if err != nil {
			return wildPokemonMsg{areaName: area, err: err}
		}
		names := make([]string, 0, len(res.PokemonEncounters))
		for _, e := range res.PokemonEncounters {
			names = append(names, e.Pokemon.Name)
		}
		return wildPokemonMsg{areaName: area, names: names}
	}
}

func (m exploreModel) catchCmd(name, ballName string) tea.Cmd {
	client := m.deps.Client
	return func() tea.Msg {
		pokemon, err := client.FetchPokemon(name)
		return pokemonFetchedMsg{pokemon: pokemon, name: name, ball: ballName, err: err}
	}
}

func (m exploreModel) Update(msg tea.Msg) (screenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case artLoadedMsg:
		m.deps.Art.handle(msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case areasLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			return m, nil
		}
		m.areas = msg.areas
		m.nextURL = msg.next
		m.prevURL = msg.prev
		if m.areaCur >= len(m.areas) {
			m.areaCur = 0
		}
		return m, nil
	case wildPokemonMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			return m, nil
		}
		m.areaName = msg.areaName
		m.wild = msg.names
		m.wildCur = 0
		m.step = wildListStep
		return m, m.requestWildArt()
	case pokemonFetchedMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			return m, nil
		}
		if err := m.deps.Dex.UseBall(msg.ball); err != nil {
			m.status = err.Error()
			return m, nil
		}
		base := 100 - msg.pokemon.BaseExperience/4
		if base < 5 {
			base = 5
		}
		chance := int(float64(base) * ball.Multiplier(msg.ball))
		if chance > 100 {
			chance = 100
		}
		if m.deps.RNG.Intn(100) < chance {
			m.deps.Dex.Add(msg.pokemon)
			m.status = fmt.Sprintf("%s was caught!", msg.name)
		} else {
			m.status = fmt.Sprintf("%s escaped!", msg.name)
		}
		_ = m.deps.Dex.Save(m.deps.SavePath)
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m exploreModel) handleKey(key tea.KeyMsg) (screenModel, tea.Cmd) {
	switch m.step {
	case areaListStep:
		switch key.String() {
		case "up", "k":
			if m.areaCur > 0 {
				m.areaCur--
			}
		case "down", "j":
			if m.areaCur < len(m.areas)-1 {
				m.areaCur++
			}
		case "n":
			if m.nextURL != nil {
				m.loading = true
				return m, tea.Batch(m.spinner.Tick, m.loadAreasCmd(m.nextURL))
			}
		case "b":
			if m.prevURL != nil {
				m.loading = true
				return m, tea.Batch(m.spinner.Tick, m.loadAreasCmd(m.prevURL))
			}
		case "enter":
			if len(m.areas) == 0 {
				return m, nil
			}
			m.loading = true
			m.status = ""
			return m, tea.Batch(m.spinner.Tick, m.loadWildCmd(m.areas[m.areaCur]))
		}
	case wildListStep:
		switch key.String() {
		case "up", "k":
			if m.wildCur > 0 {
				m.wildCur--
			}
			return m, m.requestWildArt()
		case "down", "j":
			if m.wildCur < len(m.wild)-1 {
				m.wildCur++
			}
			return m, m.requestWildArt()
		case "left", "h":
			if m.ballIdx > 0 {
				m.ballIdx--
			}
		case "right", "l":
			if m.ballIdx < len(ball.Names())-1 {
				m.ballIdx++
			}
		case "backspace":
			m.step = areaListStep
		case "enter":
			if len(m.wild) == 0 {
				return m, nil
			}
			ballName := ball.Names()[m.ballIdx]
			if m.deps.Dex.BallCount(ballName) <= 0 {
				m.status = "You have no " + ballName + " left"
				return m, nil
			}
			m.loading = true
			m.status = ""
			return m, tea.Batch(m.spinner.Tick, m.catchCmd(m.wild[m.wildCur], ballName))
		}
	}
	return m, nil
}

func (m exploreModel) requestWildArt() tea.Cmd {
	if m.step != wildListStep || len(m.wild) == 0 {
		return nil
	}
	return m.deps.Art.request(m.deps, m.wild[m.wildCur])
}

func (m exploreModel) View() string {
	var b strings.Builder
	switch m.step {
	case areaListStep:
		b.WriteString(titleStyle.Render("Explore") + "\n\n")
		if m.loading && len(m.areas) == 0 {
			b.WriteString(m.spinner.View() + " loading areas…\n")
		}
		for i, name := range m.areas {
			cursor := "  "
			line := name
			if i == m.areaCur {
				cursor = "▸ "
				line = selectedStyle.Render(line)
			}
			b.WriteString(cursor + line + "\n")
		}
		b.WriteString("\n" + helpStyle.Render("↑/↓ move · enter explore · n/b page · esc back"))
	case wildListStep:
		b.WriteString(titleStyle.Render("Explore: "+m.areaName) + "\n\n")

		var list strings.Builder
		for i, name := range m.wild {
			cursor := "  "
			line := name
			if i == m.wildCur {
				cursor = "▸ "
				line = selectedStyle.Render(line)
			}
			list.WriteString(cursor + line + "\n")
		}
		ballName := ball.Names()[m.ballIdx]
		fmt.Fprintf(&list, "\nBall: ‹ %s ›  (you have %d)\n", ballName, m.deps.Dex.BallCount(ballName))
		if m.loading {
			list.WriteString(m.spinner.View() + " throwing…\n")
		}

		if len(m.wild) > 0 {
			name := m.wild[m.wildCur]
			detail := []string{}
			if art := m.deps.Art.get(name); art != "" {
				detail = append(detail, art)
			}
			if p, ok := m.deps.Art.poke(name); ok {
				detail = append(detail, wildStatsView(p))
			} else {
				detail = append(detail, dimStyle.Render("loading stats…"))
			}
			b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
				boxStyle.Render(list.String()),
				"  ",
				boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, detail...))))
		} else {
			b.WriteString(boxStyle.Render(list.String()))
		}
		b.WriteString("\n" + helpStyle.Render("↑/↓ pick · ←/→ ball · enter throw · backspace areas · esc back"))
	}
	if m.status != "" {
		b.WriteString("\n" + statusStyle.Render(m.status))
	}
	return b.String()
}

// wildStatsView renders a wild Pokémon's base stats, types, and size.
func wildStatsView(p pokeapi.Pokemon) string {
	var b strings.Builder
	b.WriteString(p.Name + "\n")

	types := make([]string, 0, len(p.Types))
	for _, t := range p.Types {
		types = append(types, t.Type.Name)
	}
	b.WriteString("Types: " + strings.Join(types, ", ") + "\n")

	for _, s := range p.Stats {
		fmt.Fprintf(&b, "%-15s %3d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Fprintf(&b, "Height %d  Weight %d", p.Height, p.Weight)
	return b.String()
}
