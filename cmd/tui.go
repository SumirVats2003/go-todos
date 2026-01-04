package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/SumirVats2003/go-todo/internal"
	"github.com/SumirVats2003/go-todo/internal/model"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type todoApp struct {
	width      int
	height     int
	todos      []model.Todo
	repo       internal.Repository
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initialModel() todoApp {
	a := todoApp{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range a.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Content"
			t.CharLimit = 64
		}

		a.inputs[i] = t
	}

	return a
}

func (t todoApp) Init() tea.Cmd {
	return textinput.Blink
}

func (t todoApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height

	case []model.Todo:
		t.todos = msg
		return t, tea.Quit

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return t, tea.Quit
		}
	}

	return t, nil
}

func (t todoApp) View() string {
	var b strings.Builder

	for i := range t.inputs {
		b.WriteString(t.inputs[i].View())
		if i < len(t.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	return b.String()
}

func InitTeaApp(repo internal.Repository) {
	todoApp := todoApp{
		repo: repo,
	}

	p := tea.NewProgram(todoApp, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
