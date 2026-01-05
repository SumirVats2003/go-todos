package cmd

import (
	"log"

	"github.com/SumirVats2003/go-todo/internal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
)

type model struct {
	width  int
	height int
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width - 7
		m.height = msg.Height - 5
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	footerHeight := 1
	mainHeight := m.height - footerHeight

	leftWidth := int(float64(m.width) * 0.2)
	middleWidth := int(float64(m.width) * 0.4)
	rightWidth := m.width - leftWidth - middleWidth

	collections := borderStyle.
		Width(leftWidth).
		Height(m.height).
		Render("Collections\n- api/auth\n- api/login")

	request := borderStyle.
		Width(middleWidth).
		Height(mainHeight).
		Render("POST https://url.com\n\nHeaders | Body\n\nHeaders:\nKey: Value\n\nBody:\n{\n  \"key\": \"value\"\n}")

	response := borderStyle.
		Width(rightWidth).
		Height(mainHeight).
		Render("Response\n200 OK\n\n{\n  \"field\": \"value\"\n}")

	mainBody := lipgloss.JoinHorizontal(
		lipgloss.Top,
		request,
		response,
	)

	footer := borderStyle.
		Width(int(float64(m.width) * 0.8)).
		Height(footerHeight).
		Render("q: quit  |  tab: switch pane  |  enter: send request")

	rightView := lipgloss.JoinVertical(
		lipgloss.Left,
		mainBody,
		footer,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		collections,
		rightView,
	)
}

func InitTeaApp(repo internal.Repository) {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
