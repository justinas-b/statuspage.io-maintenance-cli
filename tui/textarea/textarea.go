package textarea

// A simple program demonstrating the TextArea component from the Bubbles
// component library.

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type errMsg error

type Model struct {
	TextArea textarea.Model
	err      error
}

func InitialModel() Model {
	ta := textarea.New()
	ta.CharLimit = 0
	ta.Placeholder = "Once upon a time..."
	ta.Focus()

	return Model{
		TextArea: ta,
		err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	var command tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.TextArea.Focused() {
				m.TextArea.Blur()
			}

		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.TextArea.Focused() {
				command = m.TextArea.Focus()
				commands = append(commands, command)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.TextArea, command = m.TextArea.Update(msg)
	commands = append(commands, command)
	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Maintenance description: %s\n%s\n",
		"(ctrl+c to confirm)",
		m.TextArea.View(),
	) + "\n"
}

func Run() string {
	program := tea.NewProgram(InitialModel())
	model, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}

	var value string
	if m, ok := model.(Model); ok {
		value = m.TextArea.Value()
	}
	return value
}
