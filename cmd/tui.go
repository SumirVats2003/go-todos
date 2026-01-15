package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SumirVats2003/go-todo/internal"
	"github.com/SumirVats2003/go-todo/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4"))
	completedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Strikethrough(true)
)

type todoModel struct {
	width        int
	height       int
	repo         internal.Repository
	todos        []model.Todo
	cursor       int
	selected     int
	mode         string // "list", "edit", "new"
	titleInput   textinput.Model
	contentInput textinput.Model
	message      string
}

func initialTodoModel(repo internal.Repository) todoModel {
	ti := textinput.New()
	ti.Placeholder = "Todo title"
	ti.CharLimit = 50
	ti.Width = 40

	ci := textinput.New()
	ci.Placeholder = "Todo content"
	ci.CharLimit = 200
	ci.Width = 40

	return todoModel{
		repo:         repo,
		todos:        repo.GetAllTodos(),
		cursor:       0,
		selected:     -1,
		mode:         "list",
		titleInput:   ti,
		contentInput: ci,
		message:      "",
	}
}

func (m todoModel) Init() tea.Cmd {
	return nil
}

func (m todoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case "list":
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.todos)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.todos) > 0 {
					m.selected = m.cursor
					m.mode = "edit"
					m.titleInput.SetValue(m.todos[m.cursor].Title)
					m.contentInput.SetValue(m.todos[m.cursor].Content)
					m.titleInput.Focus()
				}
			case "n":
				m.mode = "new"
				m.selected = -1
				m.titleInput.SetValue("")
				m.contentInput.SetValue("")
				m.titleInput.Focus()
			case " ":
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					todo.Completed = !todo.Completed
					m.repo.UpdateTodo(todo.Id, todo)
					m.todos = m.repo.GetAllTodos()
				}
			case "d":
				if len(m.todos) > 0 {
					m.repo.DeleteTodo(m.todos[m.cursor].Id)
					m.todos = m.repo.GetAllTodos()
					if m.cursor >= len(m.todos) && m.cursor > 0 {
						m.cursor--
					}
				}
			}
		case "edit", "new":
			switch msg.String() {
			case "esc":
				m.mode = "list"
				m.titleInput.Blur()
				m.contentInput.Blur()
				m.message = ""
			case "tab":
				if m.titleInput.Focused() {
					m.titleInput.Blur()
					m.contentInput.Focus()
				} else {
					m.contentInput.Blur()
					m.titleInput.Focus()
				}
			case "enter":
				if m.titleInput.Value() != "" {
					if m.mode == "new" {
						newTodo := model.Todo{
							Title:     m.titleInput.Value(),
							Content:   m.contentInput.Value(),
							Completed: false,
						}
						m.repo.CreateTodo(newTodo)
						m.todos = m.repo.GetAllTodos()
						m.message = "Todo created!"
					} else {
						todo := m.todos[m.selected]
						todo.Title = m.titleInput.Value()
						todo.Content = m.contentInput.Value()
						m.repo.UpdateTodo(todo.Id, todo)
						m.todos = m.repo.GetAllTodos()
						m.message = "Todo updated!"
					}
					m.mode = "list"
					m.titleInput.Blur()
					m.contentInput.Blur()
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width - 7
		m.height = msg.Height - 5
	}

	var cmd tea.Cmd
	if m.mode == "edit" || m.mode == "new" {
		m.titleInput, cmd = m.titleInput.Update(msg)
		m.contentInput, cmd = m.contentInput.Update(msg)
	}

	return m, cmd
}

func (m todoModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	footerHeight := 3
	mainHeight := m.height - footerHeight

	leftWidth := int(float64(m.width) * 0.4)
	rightWidth := m.width - leftWidth

	// Left panel - Todo list
	var todoList strings.Builder
	for i, todo := range m.todos {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}

		title := todo.Title
		if todo.Completed {
			title = completedStyle.Render(title)
		}

		if i == m.cursor {
			todoList.WriteString(selectedStyle.Render(fmt.Sprintf("%s %s %s", cursor, status, title)))
		} else {
			todoList.WriteString(fmt.Sprintf("%s %s %s", cursor, status, title))
		}
		todoList.WriteString("\n")
	}

	leftPanel := borderStyle.
		Width(leftWidth).
		Height(mainHeight).
		Render("Todos\n\n" + todoList.String())

	// Right panel - Details/Edit
	var rightPanel string
	if m.mode == "edit" || m.mode == "new" {
		title := "Edit Todo"
		if m.mode == "new" {
			title = "New Todo"
		}

		rightPanel = borderStyle.
			Width(rightWidth).
			Height(mainHeight).
			Render(fmt.Sprintf("%s\n\nTitle:\n%s\n\nContent:\n%s\n\n%s",
				title,
				m.titleInput.View(),
				m.contentInput.View(),
				m.message))
	} else if m.selected >= 0 && m.selected < len(m.todos) {
		todo := m.todos[m.selected]
		details := fmt.Sprintf("Title: %s\n\nContent: %s\n\nStatus: %s\n\nCreated: %s",
			todo.Title,
			todo.Content,
			map[bool]string{true: "Completed", false: "Active"}[todo.Completed],
			time.Unix(int64(todo.CreatedAt), 0).Format("2006-01-02 15:04:05"))

		rightPanel = borderStyle.
			Width(rightWidth).
			Height(mainHeight).
			Render("Todo Details\n\n" + details)
	} else {
		rightPanel = borderStyle.
			Width(rightWidth).
			Height(mainHeight).
			Render("Select a todo to view details")
	}

	// Footer
	footer := borderStyle.
		Width(m.width).
		Height(footerHeight).
		Render("↑↓: navigate | enter: edit | n: new | space: toggle | d: delete | q: quit")

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftPanel,
		rightPanel,
	) + "\n" + footer
}

func InitTeaApp(repo internal.Repository) {
	p := tea.NewProgram(
		initialTodoModel(repo),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
