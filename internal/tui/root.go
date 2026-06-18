package tui

import tea "github.com/charmbracelet/bubbletea"

type screen int

const (
	menuScreen screen = iota
	pokedexScreen
	exploreScreen
	battleScreen
	bagScreen
)

type screenModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (screenModel, tea.Cmd)
	View() string
}

type switchMsg struct{ to screen }

func switchTo(s screen) tea.Cmd {
	return func() tea.Msg { return switchMsg{to: s} }
}

type quitMsg struct{}

func quitCmd() tea.Msg { return quitMsg{} }

type rootModel struct {
	deps   Deps
	screen screen
	menu   menuModel
	active screenModel
	width  int
	height int
}

func newRootModel(deps Deps) rootModel {
	return rootModel{deps: deps, screen: menuScreen, menu: newMenuModel()}
}

func (m rootModel) Init() tea.Cmd { return nil }

func (m rootModel) save() {
	_ = m.deps.Dex.Save(m.deps.SavePath)
}

// newScreen builds a screen model. Later tasks replace the placeholder cases.
func (m rootModel) newScreen(s screen) screenModel {
	switch s {
	case bagScreen:
		return newBagModel(m.deps)
	case pokedexScreen:
		return newPokedexModel(m.deps)
	case exploreScreen:
		return newExploreModel(m.deps)
	case battleScreen:
		return newBattleModel(m.deps)
	}
	return newPlaceholder("?")
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.save()
			return m, tea.Quit
		case "q":
			if m.screen == menuScreen {
				m.save()
				return m, tea.Quit
			}
		case "esc":
			if m.screen != menuScreen {
				m.screen = menuScreen
				m.active = nil
				return m, nil
			}
		}
	case switchMsg:
		m.screen = msg.to
		if msg.to == menuScreen {
			m.active = nil
			return m, nil
		}
		m.active = m.newScreen(msg.to)
		return m, m.active.Init()
	case quitMsg:
		m.save()
		return m, tea.Quit
	}

	if m.screen == menuScreen {
		var cmd tea.Cmd
		m.menu, cmd = m.menu.Update(msg)
		return m, cmd
	}
	if m.active != nil {
		var cmd tea.Cmd
		m.active, cmd = m.active.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m rootModel) View() string {
	if m.screen == menuScreen {
		return m.menu.View()
	}
	if m.active != nil {
		return m.active.View()
	}
	return ""
}

type placeholderModel struct{ name string }

func newPlaceholder(name string) placeholderModel { return placeholderModel{name: name} }

func (m placeholderModel) Init() tea.Cmd { return nil }

func (m placeholderModel) Update(tea.Msg) (screenModel, tea.Cmd) { return m, nil }

func (m placeholderModel) View() string {
	return titleStyle.Render(m.name) + "\n\n" + dimStyle.Render("(coming soon)") + "\n\n" + helpStyle.Render("esc back")
}
