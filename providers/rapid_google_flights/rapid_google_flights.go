package rapidgoogleflights

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/justinm35/flyctl/types"
	"github.com/spf13/viper"
)

const providerName = "rapidgoogleflights"

type GetSearchResultsInput struct {
	SourceIata      string
	DestinationIata string
	DepartureDate   string
	Adults          int
	Currency        string
}

func SearchFlights(input GetSearchResultsInput) ([]types.FlightOffer, error) {
	client := &http.Client{}

	u, _ := url.Parse("https://google-flights2.p.rapidapi.com/api/v1/searchFlights")

	q := u.Query()

	// command defined flags
	q.Set("departure_id", input.SourceIata)
	q.Set("arrival_id", input.DestinationIata)
	q.Set("outbound_date", input.DepartureDate)
	q.Set("adults", strconv.Itoa(input.Adults))

	log.Printf("SearchFlights query: \n sourceIata: %s \n arrivalIata: %s \n deparureDate: %s \n", input.SourceIata, input.DestinationIata, input.DepartureDate)

	currency := input.Currency
	if currency == "" {
		currency = "CAD"
	}
	q.Set("currency", currency)

	// defaults
	q.Set("travel_class", "ECONOMY")
	q.Set("show_hidden", "1")
	q.Set("language_code", "en-US")
	q.Set("country_code", "CA")
	q.Set("search_type", "best")

	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)

	rapid_api_key := viper.GetString("rapid_google_api_key")
	req.Header.Add("x-rapidapi-key", rapid_api_key)
	req.Header.Set("x-rapidapi-host", "google-flights2.p.rapidapi.com")

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("SearchFlights Error: %s \n", string(err.Error()))
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var result SearchFlightResp

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&result)
	b, _ := json.MarshalIndent(result, "", "  ")
	log.Printf("SearchFlights Responses: %s \n", string(b))

	if !result.Status {
		return nil, fmt.Errorf("rapidapi google flights error: %s", flattenMessages(result.Message))
	}

	adaptedRespone, _ := adaptSearchFlightResponse(result)

	return adaptedRespone, nil
}

func adaptSearchFlightResponse(data SearchFlightResp) ([]types.FlightOffer, error) {
	all := make([]FlightOption, 0, len(data.Data.Itineraries.TopFlights)+len(data.Data.Itineraries.OtherFlights))
	all = append(all, data.Data.Itineraries.TopFlights...)
	all = append(all, data.Data.Itineraries.OtherFlights...)

	offers := make([]types.FlightOffer, 0, len(all))

	for i, opt := range all {
		segs := make([]types.Segment, 0, len(opt.Flights))

		for _, leg := range opt.Flights {
			departAt, err := time.Parse("2006-1-2 15:04", strings.TrimSpace(leg.DepartureAirport.Time))
			if err != nil {
				return nil, fmt.Errorf("parse departure time %q: %w", leg.DepartureAirport.Time, err)
			}

			arriveAt, err := time.Parse("2006-1-2 15:04", strings.TrimSpace(leg.ArrivalAirport.Time))
			if err != nil {
				return nil, fmt.Errorf("parse arrival time %q: %w", leg.ArrivalAirport.Time, err)
			}

			segs = append(segs, types.Segment{
				From:     strings.TrimSpace(leg.DepartureAirport.AirportCode),
				To:       strings.TrimSpace(leg.ArrivalAirport.AirportCode),
				DepartAt: departAt,
				ArriveAt: arriveAt,
				Carrier:  strings.TrimSpace(leg.Airline),
				FlightNo: strings.TrimSpace(leg.FlightNumber),
				Cabin:    "", // not provided in this payload
			})
		}

		offerID := strings.TrimSpace(opt.NextToken)
		if offerID == "" {
			offerID = fmt.Sprintf("offer-%d-%d", data.Timestamp, i)
		}

		offers = append(offers, types.FlightOffer{
			Provider: providerName,
			OfferID:  offerID,
			TotalPrice: types.Money{
				Amount:   int64(opt.Price) * 100, // simple: major -> minor
				Currency: "USD",                  // payload doesn’t include currency; set default or plumb it in
			},
			Segments: segs,
		})
	}

	return offers, nil
}

func flattenMessages(msg []map[string]string) string {
	if len(msg) == 0 {
		return "unknown error"
	}
	var parts []string
	for _, m := range msg {
		for k, v := range m {
			if k != "" && v != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", k, v))
			} else if v != "" {
				parts = append(parts, v)
			}
		}
	}
	if len(parts) == 0 {
		return "unknown error"
	}
	return strings.Join(parts, "; ")
}
