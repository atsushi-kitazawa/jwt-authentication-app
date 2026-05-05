package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuItem struct {
	title       string
	description string
	status      string
}

type screen int

const (
	screenHome screen = iota
	screenPlayground
)

type model struct {
	currentScreen screen
	items         []menuItem
	cursor        int
	selected      int
	width         int
	height        int
}

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("62")).
				Bold(true).
				Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	statusReadyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Bold(true)

	statusMutedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214")).
				Bold(true)

	heroStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")).
			Padding(1, 2)
)

func initialModel() model {
	items := []menuItem{
		{
			title:       "Login",
			description: "JWT login screen image. Name and password fields will live here later.",
			status:      "Ready to prototype",
		},
		{
			title:       "Users",
			description: "A list view for users fetched from the API. Think table, filters, and pagination.",
			status:      "Placeholder data",
		},
		{
			title:       "User Detail",
			description: "Detail pane for one user with update and soft-delete actions.",
			status:      "UI sketch only",
		},
		{
			title:       "Health Check",
			description: "A small network status card that can ping the unauthenticated health endpoint.",
			status:      "Easy first integration",
		},
	}

	return model{
		currentScreen: screenHome,
		items:         items,
		selected:      0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			if m.currentScreen == screenHome {
				m.currentScreen = screenPlayground
				return m, nil
			}

			m.selected = m.cursor
		case "esc", "b":
			if m.currentScreen == screenPlayground {
				m.currentScreen = screenHome
				return m, nil
			}
		case "up", "k":
			if m.currentScreen == screenPlayground && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.currentScreen == screenPlayground && m.cursor < len(m.items)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	switch m.currentScreen {
	case screenHome:
		return m.renderHome()
	case screenPlayground:
		return m.renderPlayground()
	default:
		return "unknown screen"
	}
}

func (m model) renderHome() string {
	header := titleStyle.Render("JWT Client Playground")
	subtitle := subtleStyle.Render("A small Bubble Tea sample with a simple screen transition.")

	hero := heroStyle.Width(max(48, m.width-12)).Render(strings.Join([]string{
		"Top Screen",
		"",
		"This is the entry screen for the TUI client.",
		"Press enter to move to the current playground screen and feel the transition flow.",
	}, "\n"))

	info := panelStyle.Width(max(48, m.width-12)).Render(strings.Join([]string{
		"Planned next steps",
		"",
		"- login screen",
		"- users list",
		"- user detail",
		"- health check card",
	}, "\n"))

	footer := subtleStyle.Render("enter: open sample screen   q: quit")

	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		subtitle,
		"",
		hero,
		"",
		info,
		"",
		footer,
	))
}

func (m model) renderPlayground() string {

	header := titleStyle.Render("JWT Client Playground")
	subtitle := subtleStyle.Render("Bubble Tea + Lip Gloss sample. Move with j/k or arrows, select with enter, back with b or esc.")

	left := m.renderMenu()
	right := m.renderDetail()
	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	footer := subtleStyle.Render("This is a sample TUI only. API calls will be wired in after we shape the UI. Press b to return home.")

	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		subtitle,
		"",
		content,
		"",
		footer,
	))
}

func (m model) renderMenu() string {
	var rows []string

	rows = append(rows, lipgloss.NewStyle().Bold(true).Render("Screens"))
	rows = append(rows, subtleStyle.Render("What should the client feel like?"))
	rows = append(rows, "")

	for i, item := range m.items {
		prefix := "  "
		style := itemStyle
		if i == m.cursor {
			prefix = "> "
			style = selectedItemStyle
		}

		line := fmt.Sprintf("%s%s", prefix, item.title)
		rows = append(rows, style.Render(line))
	}

	menuWidth := max(28, m.width/3)
	return panelStyle.Width(menuWidth).Render(strings.Join(rows, "\n"))
}

func (m model) renderDetail() string {
	item := m.items[m.selected]
	statusStyle := statusMutedStyle
	if m.selected == 0 || m.selected == 3 {
		statusStyle = statusReadyStyle
	}

	body := []string{
		lipgloss.NewStyle().Bold(true).Render(item.title),
		"",
		item.description,
		"",
		"Status",
		statusStyle.Render(item.status),
		"",
		"Notes",
		m.renderNotes(item.title),
	}

	detailWidth := max(42, m.width-38)
	return panelStyle.Width(detailWidth).Render(strings.Join(body, "\n"))
}

func (m model) renderNotes(title string) string {
	switch title {
	case "Login":
		return "- JWT token acquisition flow\n- persistent auth state\n- failed login feedback"
	case "Users":
		return "- scrollable list\n- selected row actions\n- refresh behavior"
	case "User Detail":
		return "- editable fields\n- optimistic updates\n- soft delete confirmation"
	case "Health Check":
		return "- unauthenticated request\n- online/offline badge\n- latency display"
	default:
		return "- custom screen notes"
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func main() {
	program := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run tui: %v\n", err)
		os.Exit(1)
	}
}
