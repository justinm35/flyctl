package amadeus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/justinm35/flyctl/types"
	"github.com/spf13/viper"
)

const providerName = "amadeus"

func SearchFlights(ctx context.Context, searchQuery types.SearchRequest) ([]types.FlightOffer, error) {
	client := &http.Client{}

	u, _ := url.Parse("https://test.api.amadeus.com/v2/shopping/flight-offers")

	q := u.Query()
	q.Set("originLocationCode", searchQuery.Origin)
	q.Set("destinationLocationCode", searchQuery.Destination)
	q.Set("departureDate", searchQuery.DepartDate.Format("2006-01-02"))
	q.Set("adults", strconv.Itoa(searchQuery.Adults))
	q.Set("max", strconv.Itoa(searchQuery.MaxResults))

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)

	req.Header.Add("Authorization", getAmadeusBearer())
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var result SearchFlightResp
	if err != nil {
		log.Fatal(err)
	} else {
		dec := json.NewDecoder(resp.Body)
		dec.Decode(&result)
	}

	adaptedRespone, _ := adaptSearchFlightResponse(result)
	return adaptedRespone, nil
}

func adaptSearchFlightResponse(data SearchFlightResp) ([]types.FlightOffer, error) {
	offers := make([]types.FlightOffer, 0, len(data.Data))

	for _, d := range data.Data {
		money, err := parseMoneyMinorUnits(d.Price.Total, d.Price.Currency, 2) // assumes 2dp
		if err != nil {
			return nil, fmt.Errorf("parse price for offer %s: %w", d.ID, err)
		}

		var segs []types.Segment
		for _, itin := range d.Itineraries {
			for _, s := range itin.Segments {
				departAt, err := parseTimeFlexible(s.Departure.At)
				if err != nil {
					return nil, fmt.Errorf("parse departure time for offer %s (%s->%s): %w",
						d.ID, s.Departure.IataCode, s.Arrival.IataCode, err)
				}

				arriveAt, err := parseTimeFlexible(s.Arrival.At)
				if err != nil {
					return nil, fmt.Errorf("parse arrival time for offer %s (%s->%s): %w",
						d.ID, s.Departure.IataCode, s.Arrival.IataCode, err)
				}

				carrier := strings.TrimSpace(s.CarrierCode)
				if op := strings.TrimSpace(s.Operating.CarrierCode); op != "" {
					carrier = op
				}

				flightNo := strings.TrimSpace(s.Number)
				if carrier != "" && flightNo != "" && !strings.HasPrefix(flightNo, carrier) {
					flightNo = carrier + flightNo
				}

				segs = append(segs, types.Segment{
					From:     s.Departure.IataCode,
					To:       s.Arrival.IataCode,
					DepartAt: departAt,
					ArriveAt: arriveAt,
					Carrier:  carrier,
					FlightNo: flightNo,
					Cabin:    "",
				})
			}
		}

		offers = append(offers, types.FlightOffer{
			Provider: providerName,
			OfferID:  d.ID,
			TotalPrice: types.Money{
				Amount:   money.Amount,
				Currency: money.Currency,
			},
			Segments: segs,
		})
	}

	return offers, nil
}

func getAmadeusBearer() string {
	baseURL := "https://test.api.amadeus.com/v1/security/oauth2/token"
	apiKey := viper.GetString("amadeus_api_key")
	apiSecret := viper.GetString("amadeus_api_secret")

	client := &http.Client{}

	reqBody := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", apiKey, apiSecret)
	reqBodyReader := strings.NewReader(reqBody)
	req, err := http.NewRequest("POST", baseURL, reqBodyReader)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	respBodyMap := responseToMap(resp.Body)

	accessToken := respBodyMap["access_token"]
	bearerPrefix := "Bearer"
	fullBearerToken := fmt.Sprintf("%s %s", bearerPrefix, accessToken)
	fmt.Println("FULL  BEARER", fullBearerToken)

	return fullBearerToken
}

func responseToMap(responseBody io.ReadCloser) map[string]any {
	bodyBytes, err := io.ReadAll(responseBody)
	if err != nil {
		log.Fatal(err)
	}

	var jsonMap map[string]any
	json.Unmarshal([]byte(bodyBytes), &jsonMap)

	return jsonMap
}

func parseTimeFlexible(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty time")
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %q", s)
}

// parseMoneyMinorUnits converts a decimal string (e.g. "123.45") into minor units (e.g. 12345).
// decimals=2 is typical, but some currencies differ (JPY=0, etc.).
func parseMoneyMinorUnits(amountStr, currency string, decimals int) (types.Money, error) {
	amountStr = strings.TrimSpace(amountStr)
	if amountStr == "" {
		return types.Money{}, fmt.Errorf("empty amount")
	}

	r := new(big.Rat)
	if _, ok := r.SetString(amountStr); !ok {
		return types.Money{}, fmt.Errorf("invalid decimal %q", amountStr)
	}

	scale := new(big.Rat).SetInt(big.NewInt(int64Pow10(decimals)))
	r.Mul(r, scale) // amount * 10^decimals

	// Round half away from zero to nearest int64
	f, _ := r.Float64()
	rounded := int64(math.Round(f))

	return types.Money{
		Amount:   rounded,
		Currency: strings.TrimSpace(currency),
	}, nil
}

func int64Pow10(n int) int64 {
	if n <= 0 {
		return 1
	}
	x := int64(1)
	for i := 0; i < n; i++ {
		x *= 10
	}
	return x
}
