package confirmation

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
)

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

type Model struct {
	cursor  int
	choice  string
	choices []string
	title   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(m.title)

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}

	return s.String()
}

func Run(title string, choices []string) string {
	model := Model{
		cursor:  0,
		choice:  "",
		choices: choices,
		title:   title,
	}
	program := tea.NewProgram(model)

	// Run returns the model as a tea.Model.
	programModel, err := program.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	var value string
	if m, ok := programModel.(Model); ok && m.choice != "" {
		value = m.choice
	}

	fmt.Printf("\nYou chose %s!\n\n", value)
	return value

}
