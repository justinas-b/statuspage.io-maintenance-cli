package textinput

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

type (
	errMsg error
)

type Model struct {
	TextInput        textinput.Model
	err              error
	errorPlaceholder string
	validationFunc   func(string) error
	title            string
}

func InitialModel(title string, placeholder string, errorPlaceholder string, validation func(string) error) Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return Model{
		TextInput:        ti,
		err:              nil,
		errorPlaceholder: errorPlaceholder,
		validationFunc:   validation,
		title:            title,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			if err := m.validationFunc(m.TextInput.Value()); err != nil {
				m.TextInput.Placeholder = m.errorPlaceholder
				m.TextInput.SetValue("")
			} else {
				return m, tea.Quit
			}
		case tea.KeyCtrlC:
			os.Exit(1)
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s (enter to confirm)\n%s\n",
		m.title,
		m.TextInput.View(),
	) + "\n"
}

func Run(title string, placeholderInitial string, placeholderError string, validationFunc func(string) error) string {
	model := InitialModel(title, placeholderInitial, placeholderError, validationFunc)
	program := tea.NewProgram(model)
	programModel, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}

	var value string
	if m, ok := programModel.(Model); ok {
		value = m.TextInput.Value()
	}
	return value
}
