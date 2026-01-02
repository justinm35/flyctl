package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/justinm35/flyctl/domain"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table    table.Model
	viewport viewport.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[0]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func RenderFlightSearchTable(f []domain.FlightOffer) {
	columns := []table.Column{
		{Title: "Route", Width: 30},
		{Title: "Departure Time", Width: 20},
		{Title: "Arrival Time", Width: 20},
		{Title: "Duration", Width: 40},
		{Title: "Price", Width: 15},
		{Title: "Carrier", Width: 20},
		{Title: "Seats", Width: 4},
	}

	rows := formatResponseData(f)

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

	m := model{
		table:    t,
		viewport: FlightPreview("STUFF"),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func formatResponseData(offers []domain.FlightOffer) []table.Row {
	var allRows []table.Row

	for _, o := range offers {
		if len(o.Segments) == 0 {
			continue
		}

		// -------- Route string (A → B → C)
		var routeString string
		for i, seg := range o.Segments {
			if i == 0 {
				routeString = seg.From
			}
			routeString = fmt.Sprintf("%s → %s", routeString, seg.To)
		}

		// -------- Time string (dep → arr per segment, newline between segments)
		const outLayout = "Mon, Jan 2, 3:04 PM"
		var timeString string
		for i, seg := range o.Segments {
			part := fmt.Sprintf("%s → %s", seg.DepartAt.UTC().Format(outLayout), seg.ArriveAt.UTC().Format(outLayout))
			if i == 0 {
				timeString = part
			} else {
				timeString = fmt.Sprintf("%s \n%s", timeString, part)
			}
		}

		// Departure Time
		departureTime := o.Segments[0].DepartAt.UTC().Format(outLayout)
		arrivalTime := o.Segments[len(o.Segments)-1].ArriveAt.UTC().Format(outLayout)

		// -------- Duration string (per segment + total)
		var totalDur time.Duration
		var durationParts []string
		for _, seg := range o.Segments {
			d := seg.ArriveAt.Sub(seg.DepartAt)
			if d < 0 {
				// guard for weird timezone/provider issues
				d = 0
			}
			totalDur += d
			durationParts = append(durationParts, formatDuration(d))
		}
		// mimic your "a | b | c" style; append total at end
		totalDurationString := strings.Join(durationParts, " | ")
		if totalDurationString == "" {
			totalDurationString = formatDuration(totalDur)
		} else {
			totalDurationString = fmt.Sprintf("%s | total %s", totalDurationString, formatDuration(totalDur))
		}

		// -------- Price (Money is minor units)
		totalPrice := formatMoney(o.TotalPrice)

		// -------- Carrier (choose unique carriers encountered)
		carrierString := joinUniqueCarriers(o.Segments)

		// -------- Seats remaining (not in domain yet)
		seatsRemaining := "-" // you removed it from domain.FlightOffer

		allRows = append(allRows, table.Row{
			routeString,
			departureTime,
			arrivalTime,
			totalDurationString,
			totalPrice,
			carrierString,
			seatsRemaining,
		})
	}

	return allRows
}
func formatDuration(d time.Duration) string {
	// "2h15m" -> "2h 15m"
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}

func formatMoney(m domain.Money) string {
	// assumes 2dp; matches your adapter parseMoneyMinorUnits(..., 2)
	abs := m.Amount
	sign := ""
	if abs < 0 {
		sign = "-"
		abs = -abs
	}
	major := abs / 100
	minor := abs % 100
	return fmt.Sprintf("%s%s %d.%02d", sign, m.Currency, major, minor)
}

func joinUniqueCarriers(segs []domain.Segment) string {
	seen := map[string]struct{}{}
	var carriers []string
	for _, s := range segs {
		c := strings.TrimSpace(s.Carrier)
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		carriers = append(carriers, c)
	}
	if len(carriers) == 0 {
		return "-"
	}
	return strings.Join(carriers, ", ")
}
