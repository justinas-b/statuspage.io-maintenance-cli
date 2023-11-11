package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/justinas-b/statuspage.io-maintenance-cli/client"
	"os"
	"sort"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Select   key.Binding
	Continue key.Binding
}

var keys = keyMap{
	Select: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "Select"),
	),
	Continue: key.NewBinding(
		key.WithKeys("↲"),
		key.WithHelp("↲", "Continue"),
	),
}

type Item struct {
	Page    *client.Page
	Checked bool
}

func (i Item) Title() string {

	var check string
	if i.Checked {
		check = "✓"
	} else {
		check = "✗"
	}

	return fmt.Sprintf("%s %s", check, i.Page.Name)
}

func (i Item) Description() string { return i.Page.Domain }
func (i Item) FilterValue() string { return i.Page.Name }

type Model struct {
	List list.Model
}

func InitialModel(items []list.Item) Model {
	m := Model{
		List: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.List.Title = "StatusPage.io pages:"
	m.List.SetShowPagination(true)
	m.List.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Select,
			keys.Continue,
		}
	}
	m.List.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Select,
			keys.Continue,
		}
	}
	return m
}
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) GetIndexOfSelectedItem(i Item) int {
	allItems := m.List.Items()
	for idx, item := range allItems {
		if item == i {
			return idx
		}
	}
	panic("Item not found")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.List.FilterState() == list.Filtering {
			break
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			os.Exit(1)
			//return m, tea.Quit
		case tea.KeySpace:
			if m.List.IsFiltered() {
				selectedItem := m.List.SelectedItem().(Item)
				idx := m.GetIndexOfSelectedItem(selectedItem)
				selectedItem.Checked = !selectedItem.Checked
				m.List.SetItem(idx, selectedItem)
				m.List.ResetFilter()

			} else {
				i := m.List.SelectedItem().(Item)
				i.Checked = !i.Checked
				m.List.SetItem(m.List.Index(), i)
			}

			// TODO: extract as method to sort list
			items := m.List.Items()
			sort.Slice(items, func(i, j int) bool {
				if items[i].(Item).Checked && !items[j].(Item).Checked {
					return true
				}
				return false
			})
			m.List.SetItems(items)

		case tea.KeyEnter:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	//_, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.List.View())
}

func pagesToItems(pages []*client.Page) []list.Item {
	var items []list.Item
	for idx := range pages {
		items = append(
			items,
			Item{
				Page:    pages[idx],
				Checked: false,
			},
		)
	}
	return items
}

func isSelected(i list.Item) bool {
	return i.(Item).Checked
}

func Run(pages []*client.Page) []*client.Page {
	items := pagesToItems(pages)
	model := InitialModel(items)
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run returns the listModel as a tea.Model.
	programModel, err := program.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Filter the initial list of items
	filteredItems := make([]*client.Page, 0)
	fmt.Printf("Following status pages will be affected:\n")
	if m, ok := programModel.(Model); ok {
		for _, i := range m.List.Items() {
			if isSelected(i) {
				//filteredItems = append(filteredItems, &(*i.(Item).Page))
				filteredItems = append(filteredItems, i.(Item).Page)
				x := i.(Item).Page
				_ = x
				fmt.Printf("🪓 %s!\n", i.(Item).Page.Name)
			}
		}
	}
	fmt.Printf("\n")

	// Exit program if none if the items are selected
	if len(filteredItems) == 0 {
		fmt.Println("No items selected")
		os.Exit(1)
	}
	return filteredItems
}
