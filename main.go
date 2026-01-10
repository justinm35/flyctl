package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/justinm35/flyctl/domain"
	"github.com/justinm35/flyctl/styles"
)

type screen int

const (
	screenSearch screen = iota
	screenResults
	screenFlightDetails
	screenCount
)

var allScreens = []screen{
	screenSearch,
	screenResults,
	screenFlightDetails,
}

func main() {
	InitConfig()

	m := NewModel()
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Model: App State
type Model struct {
	focusedPane         int
	screen              screen
	screenSearch        SearchState
	screenResults       ResultsState
	screenFlightDetails FlightDetailsState
	width               int
	height              int
}

type searchResultsMsg struct{ offers []domain.FlightOffer }
type flightDetailsSelectedMsg struct{ offer domain.FlightOffer }

type errMsg struct{ err error }

// NewModel: Initial model
func NewModel() Model {
	return Model{
		focusedPane:         0,
		screen:              screenSearch,
		screenSearch:        newSearchState(),
		screenResults:       newResultsState(),
		screenFlightDetails: newFlightDetailsState(),
	}
}

// Init: Kick off the event loop
func (m Model) Init() tea.Cmd {
	return m.screenSearch.initCmd()
}

// Update: handle Msgs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		if km.String() == "tab" {
			if m.focusedPane == int(screenCount)-1 {
				m.focusedPane = 0
			} else {
				m.focusedPane += 1
				return m, nil
			}
		}
		if km.String() == "shift+tab" {
			if m.focusedPane == 0 {
				m.focusedPane = int(screenCount) - 1
			} else {
				m.focusedPane -= 1
				return m, nil
			}

		}
		if km.String() == "ctrl+z" {
			return m, tea.Suspend
		}
		if km.String() == "ctrl+c" || km.String() == "esc" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case searchResultsMsg:
		m.screenSearch.loading = false
		m.focusedPane = 1
		m.screenResults.offers = msg.offers
		m.screenResults.buildTable(m.width)
		m.screen = screenResults
		return m, nil
	case flightDetailsSelectedMsg:
		m.screenFlightDetails.initFlightDetails(msg.offer)
		m.screen = screenFlightDetails
	}

	switch allScreens[m.focusedPane] {
	case screenSearch:
		return updateSearch(m, msg)
	case screenResults:
		return updateResults(m, msg)
	case screenFlightDetails:
		return updateFlightDetails(m, msg)
	default:
		return m, nil
	}
}

// View: Return a string based on the state of our model
func (m Model) View() string {
	totalHeight := m.height - 2
	totalWidth := m.width - 4
	halfWidth := totalWidth / 2
	halfHeight := m.height / 2

	paneStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.White)

	// Top pane (results)
	var topPane string
	if allScreens[m.focusedPane] == screenResults {
		topPane = paneStyle.
			Height(halfHeight - 30).
			Width(totalWidth).
			BorderForeground(styles.NeonPurple).
			Render(viewResults(m))
	} else {
		topPane = paneStyle.
			Height(halfHeight - 30).
			Width(totalWidth).
			Render(viewResults(m))
	}

	// Bottom left (search)
	var bottomLeftPane string
	if allScreens[m.focusedPane] == screenSearch {
		bottomLeftPane = paneStyle.
			Height(totalHeight - halfHeight).
			Width(halfWidth).
			BorderForeground(styles.NeonPurple).
			Render(viewSeach(m))
	} else {
		bottomLeftPane = paneStyle.
			Height(totalHeight - halfHeight).
			Width(halfWidth).
			Render(viewSeach(m))
	}

	// Bottom right (details)
	var bottomRightPane string
	if allScreens[m.focusedPane] == screenFlightDetails {
		bottomRightPane = paneStyle.
			Height(totalHeight - halfHeight).
			Width(halfWidth).
			BorderForeground(styles.NeonPurple).
			Render(viewFlightDetails(m))
	} else {
		bottomRightPane = paneStyle.
			Height(totalHeight - halfHeight).
			Width(halfWidth).
			Render(viewFlightDetails(m))
	}

	bottomBar := lipgloss.NewStyle().
		Foreground(styles.MutedGray).
		Width(totalWidth).
		Render("(esc to quit)")

	bottomHalf := lipgloss.JoinHorizontal(lipgloss.Bottom, bottomLeftPane, bottomRightPane)
	fullView := lipgloss.JoinVertical(lipgloss.Top, topPane, bottomHalf, bottomBar)

	return fullView
}
