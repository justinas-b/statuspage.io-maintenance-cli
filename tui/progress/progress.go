package progress

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/justinas-b/statuspage.io-maintenance-cli/client"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Model struct {
	items                  []*client.Page
	index                  int
	width                  int
	height                 int
	spinner                spinner.Model
	progress               progress.Model
	done                   bool
	maintenanceTitle       string
	maintenanceDescription string
	maintenanceStartDate   string
	maintenanceDuration    string
}

type processedItemMsg struct {
	page  *client.Page
	error error
}

var (
	currentItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle        = lipgloss.NewStyle().Margin(1, 2)
	checkMark        = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	failureMark      = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).SetString("✗")
)

func InitialModel(pages []*client.Page, maintenanceTitle string, maintenanceDescription string, maintenanceStartDate string, maintenanceDuration string) Model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return Model{
		items:                  pages,
		spinner:                s,
		progress:               p,
		maintenanceTitle:       maintenanceTitle,
		maintenanceDescription: maintenanceDescription,
		maintenanceStartDate:   maintenanceStartDate,
		maintenanceDuration:    maintenanceDuration,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		processItem(
			m.items[m.index],
			m.maintenanceTitle,
			m.maintenanceDescription,
			m.maintenanceStartDate,
			m.maintenanceDuration,
		),
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case processedItemMsg:
		// Update progress bar
		progressCommand := m.progress.SetPercent(float64(m.index+1) / float64(len(m.items)))

		m.index++

		// Set the mark depending if request has failed or succeeded
		var (
			mark         lipgloss.Style
			nextCommand  tea.Cmd
			errorMessage string
		)

		if msg.error != nil {
			mark = failureMark
			errorMessage = fmt.Sprintf("(%s)", msg.error)
		} else {
			mark = checkMark
		}

		if m.index >= len(m.items) {
			// Everything's been installed. We're done!
			cmd := tea.Quit
			nextCommand = cmd
			m.done = true
			//return m, tea.Sequence(
			//	progressCommand,
			//	tea.Printf("%s %s", mark, m.items[m.index-1]),
			//	tea.Quit,
			//)
		} else {
			cmd := processItem(
				m.items[m.index],
				m.maintenanceTitle,
				m.maintenanceDescription,
				m.maintenanceStartDate,
				m.maintenanceDuration,
			)
			nextCommand = cmd
			//return m, tea.Sequence(
			//	progressCommand,
			//	tea.Printf("%s %s", mark, m.items[m.index-1]), // print success message above our program
			//	processItem(m.items[m.index]),                 // download the next package
			//)
		}
		return m, tea.Sequence(
			progressCommand,
			tea.Printf("%s %s %s", mark, m.items[m.index-1], errorMessage),
			nextCommand,
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	n := len(m.items)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! Processed %d items.\n", n))
	}

	itemCount := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)
	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+itemCount))
	itemName := currentItemStyle.Render(m.items[m.index].String())
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Processing " + itemName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+itemCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + itemCount
}

func processItem(item *client.Page, maintenanceTitle string, maintenanceDescription string, maintenanceStartDate string, maintenanceDuration string) tea.Cmd {
	// This is where you'd do i/o stuff to download and install items. In
	// our case we're just pausing for a moment to simulate the process.
	d := time.Millisecond * time.Duration(rand.Intn(100)) //nolint:gosec
	duration, _ := strconv.Atoi(maintenanceDuration)
	startTime, _ := time.Parse(time.DateTime, maintenanceStartDate)
	endTime := startTime.Add(time.Hour * time.Duration(duration))

	err := item.SetMaintenance(
		maintenanceTitle,
		maintenanceDescription,
		startTime,
		endTime,
	)
	
	return tea.Tick(d, func(t time.Time) tea.Msg {
		time.Sleep(2 * time.Second)
		return processedItemMsg{
			page:  item,
			error: err,
		}
	})

	//return tea.Cmd(update(item))
}

func Run(pages []*client.Page, maintenanceTitle string, maintenanceDescription string, maintenanceStartDate string, maintenanceDuration string) {
	model := InitialModel(pages, maintenanceTitle, maintenanceDescription, maintenanceStartDate, maintenanceDuration)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
