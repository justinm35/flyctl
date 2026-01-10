package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/justinm35/flyctl/styles"
	"github.com/justinm35/flyctl/types"
)

type FlightDetailsState struct {
	offer    types.FlightOffer
	viewport viewport.Model
	err      string
}

func (flightDetailsState *FlightDetailsState) initFlightDetails(selectedOffer types.FlightOffer) {
	vp := viewport.New(50, 50)

	flightDetailsState.viewport = vp

	flightDetailsState.offer = selectedOffer
}

func newFlightDetailsState() FlightDetailsState {
	return FlightDetailsState{}
}

func updateFlightDetails(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	headerH := 1
	// 	footerH := 1
	// 	verticalPadding := 0
	// 	m.screenFlightDetails.viewport.Width = msg.Width
	// 	m.screenFlightDetails.viewport.Height = msg.Height - headerH - footerH - verticalPadding
	// 	if m.screenFlightDetails.viewport.Height < 1 {
	// 		m.screenFlightDetails.viewport.Height = 1
	// 	}
	// 	return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, nil
		case "b":
			m.screen = screenResults
			return m, nil
		}
	}

	return m, cmd
}

func viewFlightDetails(m Model) string {
	const width = 78
	const glamourGutter = 2
	vp := m.screenFlightDetails.viewport
	halfHeight := (m.height / 2) - 2

	vp.Height = halfHeight
	vp.Width = m.width / 2

	lipGlossRender := lipGlossRender(m.screenFlightDetails.offer, m.width)

	vp.SetContent(lipGlossRender)

	return vp.View()
}

func lipGlossRender(offer types.FlightOffer, width int) string {

	const dateLayout = "Mon, 02 Jan 2006"
	const timeLayout = "15:04 MST"

	header := lipgloss.NewStyle().Foreground(styles.NeonPurple).Bold(true).Width(30).Render("[Selected Flight Details] \n")
	noResults := lipgloss.NewStyle().Foreground(styles.MutedGray).Align(lipgloss.Center).MarginTop(6).Width(width / 2).Render("Search & Select a flight...")
	if offer.Segments == nil {
		return fmt.Sprintf("%s \n\n\n %s", header, noResults)
	}

	departingFlihtLine := fmt.Sprintf("Departure Date: %s", offer.Segments[0].DepartAt.UTC().Format(dateLayout))
	totalPriceLine := fmt.Sprintf("Price (%s): %d", offer.TotalPrice.Currency, offer.TotalPrice.Amount)

	departingFlight := lipgloss.NewStyle().Render(departingFlihtLine)
	totalPrice := lipgloss.NewStyle().Render(totalPriceLine)

	spacer := lipgloss.NewStyle().Width(6).Render("") // 4-character gap

	departureAndPrice := lipgloss.JoinHorizontal(
		lipgloss.Left,
		departingFlight,
		spacer,
		totalPrice,
	)

	renderer, err := glamour.NewTermRenderer()
	if err != nil {
		return ""
	}
	var b strings.Builder

	b.WriteString("```text\n")
	for i, s := range offer.Segments {

		if i != 0 {
			fmt.Fprintf(&b, "│\n")
		}

		fmt.Fprintf(&b, "○ %s %s  \n", s.DepartAt.UTC().Format(timeLayout), s.From)
		fmt.Fprintf(&b, "│  \n")
		fmt.Fprintf(&b, "│ Travel Time: %s  \n", formatDuration(s.ArriveAt.Sub(s.DepartAt)))
		fmt.Fprintf(&b, "│  \n")
		fmt.Fprintf(&b, "○ %s\n", s.ArriveAt.UTC().Format(timeLayout))
		fmt.Fprintf(&b, "│ %s · %s · %s\n", emptyDash(s.Carrier), emptyDash(s.FlightNo), emptyDash(s.Cabin))
		fmt.Fprintf(&b, "│\n")
		if len(offer.Segments) > i+1 {
			next := offer.Segments[i+1]
			layover := next.DepartAt.Sub(s.ArriveAt)

			fmt.Fprintf(&b, "────────────────────────────────────────────────────────────────\n")
			fmt.Fprintf(&b, "%s layover • %s\n", formatDuration(layover), s.To)
			fmt.Fprintf(&b, "────────────────────────────────────────────────────────────────\n")
		}
	}
	b.WriteString("```\n\n")

	routeDetails, err := renderer.Render(b.String())
	if err != nil {
		return ""
	}

	fillView := lipgloss.JoinVertical(lipgloss.Left, header, departureAndPrice, routeDetails)
	return lipgloss.NewStyle().Render(fillView)
}

func offerMarkdown(offer types.FlightOffer) string {
	var b strings.Builder

	const dateLayout = "Mon, 02 Jan 2006"
	const timeLayout = "15:04 MST"

	// Summary
	fmt.Fprintf(&b, "### Selected Journey: %s\n\n", routeLine(offer.Segments))
	fmt.Fprintf(&b, "**Depature %s**\n\n", offer.Segments[0].DepartAt.UTC().Format(dateLayout))
	fmt.Fprintf(&b, "**Price (%s): %d**\n\n", offer.TotalPrice.Currency, offer.TotalPrice.Amount)

	// Segments
	if len(offer.Segments) == 0 {
		b.WriteString("> No segments available.\n\n")
		return b.String()
	}

	b.WriteString("```text\n")
	for i, s := range offer.Segments {

		if i != 0 {
			fmt.Fprintf(&b, "│\n")
		}

		fmt.Fprintf(&b, "○ %s %s  \n", s.DepartAt.UTC().Format(timeLayout), s.From)
		fmt.Fprintf(&b, "│  \n")
		fmt.Fprintf(&b, "│ Travel Time: %s  \n", formatDuration(s.ArriveAt.Sub(s.DepartAt)))
		fmt.Fprintf(&b, "│  \n")
		fmt.Fprintf(&b, "○ %s\n", s.ArriveAt.UTC().Format(timeLayout))
		fmt.Fprintf(&b, "│ %s · %s · %s\n", emptyDash(s.Carrier), emptyDash(s.FlightNo), emptyDash(s.Cabin))
		fmt.Fprintf(&b, "│\n")
		if len(offer.Segments) > i+1 {

			next := offer.Segments[i+1]
			layover := next.DepartAt.Sub(s.ArriveAt)

			fmt.Fprintf(&b, "│────────────────────────────────────────────────────────────────\n")
			fmt.Fprintf(&b, "│ Layover: %s\n", formatDuration(layover))
			fmt.Fprintf(&b, "│────────────────────────────────────────────────────────────────\n")
		}
	}

	b.WriteString("```\n\n")

	b.WriteString("\n\n*(Press **b** to go back)*\n")

	return b.String()
}

func routeLine(segs []types.Segment) string {
	if len(segs) == 0 {
		return "-"
	}
	var parts []string
	parts = append(parts, emptyDash(segs[0].From))
	for _, s := range segs {
		parts = append(parts, emptyDash(s.To))
	}
	return strings.Join(parts, " → ")
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
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

func emptyDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
