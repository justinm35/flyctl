package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/justinm35/flyctl/domain"
	"github.com/justinm35/flyctl/styles"
	"github.com/justinm35/flyctl/utils"
)

type ResultsState struct {
	table  table.Model
	offers []domain.FlightOffer
	err    string
}

func (resultsState *ResultsState) buildTable(width int) {
	inner := width - 14 // leave some space for borders/margins
	if inner < 40 {     // guard for tiny terminals
		inner = width
	}
	routeW := int(0.20 * float64(inner))
	departureW := int(0.14 * float64(inner))
	arrivalW := int(0.14 * float64(inner))
	durationW := int(0.22 * float64(inner))
	priceW := int(0.10 * float64(inner))
	carrierW := int(0.20 * float64(inner))

	// Make seats whatever is left so the total fits exactly
	columns := []table.Column{
		{Title: "Route", Width: routeW},
		{Title: "Departure Time", Width: departureW},
		{Title: "Arrival Time", Width: arrivalW},
		{Title: "Duration", Width: durationW},
		{Title: "Price", Width: priceW},
		{Title: "Carrier", Width: carrierW},
	}
	rows := utils.FormatResponseData(resultsState.offers)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)
	resultsState.table = t
}

func newResultsState() ResultsState {
	t := table.New(
		table.WithColumns([]table.Column{}),
	)
	return ResultsState{table: t}
}

func updateResults(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Select and open the details
		case "enter":
			return m, getFlightDetailsCmd(m)
		}
	}

	m.screenResults.table, cmd = m.screenResults.table.Update(msg)

	return m, cmd
}

func viewResults(m Model) string {
	s := ""
	s += lipgloss.NewStyle().Foreground(styles.NeonPurple).Bold(true).Width(30).Render("[Results]")
	s += "\n"
	s += m.screenResults.table.View()
	return s
}

func getFlightDetailsCmd(model Model) tea.Cmd {
	return func() tea.Msg {
		idx := model.screenResults.table.Cursor()
		if idx < 0 || idx >= len(model.screenResults.offers) {
			return nil
		}
		offer := model.screenResults.offers[idx]
		return flightDetailsSelectedMsg{offer: offer}
	}
}
