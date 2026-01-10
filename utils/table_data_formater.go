package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/justinm35/flyctl/domain"
)

func FormatResponseData(offers []domain.FlightOffer) []table.Row {
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
		// seatsRemaining := "-" // you removed it from domain.FlightOffer

		allRows = append(allRows, table.Row{
			routeString,
			departureTime,
			arrivalTime,
			totalDurationString,
			totalPrice,
			carrierString,
			// seatsRemaining,
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
