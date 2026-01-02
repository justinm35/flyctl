package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	rapidgoogleflights "github.com/justinm35/flyctl/providers/rapid_google_flights"
)

type SearchState struct {
	inputs []textinput.Model
	focus  int
	err    string
}

func newSearchState() SearchState {
	makeInput := func(placeholder string) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.Prompt = "> "
		ti.CharLimit = 64
		ti.Width = 30
		return ti
	}

	inputs := []textinput.Model{
		makeInput("Source IATA (e.g. CPH)"),
		makeInput("Destination IATA (e.g. YYZ)"),
		makeInput("Departure date (YYYY-MM-DD)"),
	}

	inputs[0].Focus()

	return SearchState{
		inputs: inputs,
		focus:  0,
	}

}

func (s SearchState) initCmd() tea.Cmd { return textinput.Blink }

func updateSearch(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			if msg.String() == "shift+tab" || msg.String() == "up" {
				m.screenSearch.focus--
			} else {
				m.screenSearch.focus++
			}

			if m.screenSearch.focus < 0 {
				m.screenSearch.focus = len(m.screenSearch.inputs) - 1
			} else if m.screenSearch.focus >= len(m.screenSearch.inputs) {
				m.screenSearch.focus = 0
			}

			for i := range m.screenSearch.inputs {
				if i == m.screenSearch.focus {
					m.screenSearch.inputs[i].Focus()
					m.screenSearch.inputs[i].PromptStyle = m.screenSearch.inputs[i].PromptStyle.Bold(true)
					m.screenSearch.inputs[i].TextStyle = m.screenSearch.inputs[i].TextStyle.Bold(true)
				} else {
					m.screenSearch.inputs[i].Blur()
					m.screenSearch.inputs[i].PromptStyle = m.screenSearch.inputs[i].PromptStyle.Bold(false)
					m.screenSearch.inputs[i].TextStyle = m.screenSearch.inputs[i].TextStyle.Bold(false)
				}
			}

			return m, nil
		case "enter":
			// TODO: Validate input
			return m, getSearchResultsCmd(m)
		}
	}
	// Let the focused input handle the message
	var cmd tea.Cmd
	m.screenSearch.inputs[m.screenSearch.focus], cmd = m.screenSearch.inputs[m.screenSearch.focus].Update(msg)

	return m, cmd
}

func viewSeach(m Model) string {
	s := "Flight search\n\n"
	labels := []string{"From", "To", "Depart", "Return"}

	for i := range m.screenSearch.inputs {
		s += labels[i] + ":\n" + m.screenSearch.inputs[i].View() + "\n\n"
	}

	s += "(tab to switch fields, enter next/submit, esc to quit)\n"
	return s
}

func getSearchResultsCmd(model Model) tea.Cmd {
	return func() tea.Msg {
		input := rapidgoogleflights.GetSearchResultsInput{
			SourceIata:      model.screenSearch.inputs[0].Value(),
			DestinationIata: model.screenSearch.inputs[1].Value(),
			DepartureDate:   model.screenSearch.inputs[2].Value(),
			Adults:          1,
		}
		offers, err := rapidgoogleflights.SearchFlights(input)
		if err != nil {
			log.Printf("SearchFlights error: %v", err)
			log.Fatal(err)
		}

		return searchResultsMsg{offers: offers}
	}

}
