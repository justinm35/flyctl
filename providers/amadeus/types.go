package amadeus

type SearchFlightResp struct {
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
	Data []struct {
		Type                  string `json:"type"`
		ID                    string `json:"id"`
		Source                string `json:"source"`
		OneWay                bool   `json:"oneWay"`
		NumberOfBookableSeats int    `json:"numberOfBookableSeats"`
		Itineraries           []struct {
			Duration string `json:"duration"`
			Segments []struct {
				Departure struct {
					IataCode string `json:"iataCode"`
					Terminal string `json:"terminal"`
					At       string `json:"at"`
				} `json:"departure"`
				Arrival struct {
					IataCode string `json:"iataCode"`
					Terminal string `json:"terminal"`
					At       string `json:"at"`
				} `json:"arrival,omitempty"`
				CarrierCode string `json:"carrierCode"`
				Number      string `json:"number"`
				Operating   struct {
					CarrierCode string `json:"carrierCode"`
				} `json:"operating"`
				Duration      string `json:"duration"`
				NumberOfStops int    `json:"numberOfStops"`
			} `json:"segments"`
		} `json:"itineraries"`
		Price struct {
			Currency   string `json:"currency"`
			Total      string `json:"total"`
			Base       string `json:"base"`
			GrandTotal string `json:"grandTotal"`
		} `json:"price"`
	} `json:"data"`
}
