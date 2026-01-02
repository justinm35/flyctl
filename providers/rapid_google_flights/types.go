package rapidgoogleflights

type SearchFlightResp struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Data      struct {
		Itineraries struct {
			TopFlights   []FlightOption `json:"topFlights"`
			OtherFlights []FlightOption `json:"otherFlights"`
		} `json:"itineraries"`
	} `json:"data"`
}
type FlightOption struct {
	DepartureTime string          `json:"departure_time"`
	ArrivalTime   string          `json:"arrival_time"`
	Duration      DurationInfo    `json:"duration"`
	Flights       []FlightLeg     `json:"flights"`
	Layovers      []Layover       `json:"layovers"`
	Bags          BagsInfo        `json:"bags"`
	Carbon        CarbonEmissions `json:"carbon_emissions"`
	Price         int             `json:"price"`
	Stops         int             `json:"stops"`
	AirlineLogo   string          `json:"airline_logo"`
	NextToken     string          `json:"next_token"`
}

type DurationInfo struct {
	Raw  int    `json:"raw"`
	Text string `json:"text"`
}

type FlightLeg struct {
	DepartureAirport AirportTime `json:"departure_airport"`
	ArrivalAirport   AirportTime `json:"arrival_airport"`

	DurationLabel string   `json:"duration_label"`
	Duration      int      `json:"duration"`
	Airline       string   `json:"airline"`
	AirlineLogo   string   `json:"airline_logo"`
	FlightNumber  string   `json:"flight_number"`
	Aircraft      string   `json:"aircraft"`
	Seat          string   `json:"seat"`
	Legroom       string   `json:"legroom"`
	Extensions    []string `json:"extensions"`
}

type AirportTime struct {
	AirportName string `json:"airport_name"`
	AirportCode string `json:"airport_code"`
	Time        string `json:"time"`
}

type Layover struct {
	AirportCode   string `json:"airport_code"`
	AirportName   string `json:"airport_name"`
	DurationLabel string `json:"duration_label"`
	Duration      int    `json:"duration"`
	City          string `json:"city"`
}

type BagsInfo struct {
	CarryOn int  `json:"carry_on"`
	Checked *int `json:"checked"`
}

type CarbonEmissions struct {
	DifferencePercent int `json:"difference_percent"`
	CO2e              int `json:"CO2e"`
	TypicalForRoute   int `json:"typical_for_this_route"`
	Higher            int `json:"higher"`
}
