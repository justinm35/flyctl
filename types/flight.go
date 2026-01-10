package types

import "time"

type SearchRequest struct {
	Origin      string
	Destination string
	DepartDate  time.Time
	ReturnDate  *time.Time
	Adults      int
	MaxResults  int
	Currency    string
}

type FlightOffer struct {
	Provider   string
	OfferID    string
	TotalPrice Money
	Segments   []Segment
}

type Segment struct {
	From     string
	To       string
	DepartAt time.Time
	ArriveAt time.Time
	Carrier  string
	FlightNo string
	Cabin    string
}

type Money struct {
	Amount   int64
	Currency string
}
