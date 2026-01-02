package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/justinm35/flyctl/domain"
	"github.com/justinm35/flyctl/utils"
)

type ResultsState struct {
	table  table.Model
	offers []domain.FlightOffer
	err    string
}

func (resultsState *ResultsState) buildTable() {
	columns := []table.Column{
		{Title: "Route", Width: 30},
		{Title: "Departure Time", Width: 20},
		{Title: "Arrival Time", Width: 20},
		{Title: "Duration", Width: 40},
		{Title: "Price", Width: 15},
		{Title: "Carrier", Width: 20},
		{Title: "Seats", Width: 4},
	}
	rows := utils.FormatResponseData(resultsState.offers)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	resultsState.table = t
}

func newResultsState() ResultsState {
	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Route", Width: 12},
			{Title: "Price", Width: 8},
		}),
	)
	t.SetHeight(12)
	return ResultsState{table: t}
}

// func (s ResultsState) initCmd() tea.Cmd { return textinput.Blink }

func updateResults(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.screenResults.table, cmd = m.screenResults.table.Update(msg)

	return m, cmd
}

func viewResults(m Model) string {
	return "Results\n\n" + m.screenResults.table.View() + "\n\n(esc to quit)\n"
}
