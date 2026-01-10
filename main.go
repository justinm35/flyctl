package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/justinm35/flyctl/domain"
)

type screen int

const (
	screenSearch screen = iota
	screenResults
	screenFlightDetails
)

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
	screen              screen
	screenSearch        SearchState
	screenResults       ResultsState
	screenFlightDetails FlightDetailsState
}

type searchResultsMsg struct{ offers []domain.FlightOffer }
type flightDetailsSelectedMsg struct{ offer domain.FlightOffer }

type errMsg struct{ err error }

// NewModel: Initial model
func NewModel() Model {
	return Model{
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
		if km.String() == "ctrl+z" {
			return m, tea.Suspend
		}
		if km.String() == "ctrl+c" || km.String() == "esc" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case searchResultsMsg:
		m.screenSearch.loading = false
		m.screenResults.offers = msg.offers
		m.screenResults.buildTable()
		m.screen = screenResults
		return m, nil
	case flightDetailsSelectedMsg:
		m.screenFlightDetails.initFlightDetails(msg.offer)
		m.screen = screenFlightDetails
	}

	switch m.screen {
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
	switch m.screen {
	case screenSearch:
		return viewSeach(m)
	case screenResults:
		return viewResults(m)
	case screenFlightDetails:
		return viewFlightDetails(m)
	default:
		return ""
	}
}
