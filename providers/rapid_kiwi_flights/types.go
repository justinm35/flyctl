package rapidkiwiflights

type ItinerariesResponse struct {
	Typename    string      `json:"__typename"`
	Metadata    Metadata    `json:"metadata"`
	Itineraries []Itinerary `json:"itineraries"`
}

type Metadata struct {
	TopFiveResultsBaggageEligibility TopFiveResultsBaggageEligibility `json:"topFiveResultsBaggageEligibility"`
	Carriers                         []Carrier                        `json:"carriers"`
	StopoverCountries                []string                         `json:"stopoverCountries"`
	InboundDays                      []string                         `json:"inboundDays"`
	OutboundDays                     []string                         `json:"outboundDays"`
	TravelTips                       []interface{}                    `json:"travelTips"`
	TopResults                       TopResults                       `json:"topResults"`

	PriceAlertExists   *bool                 `json:"priceAlertExists"`
	ExistingPriceAlert interface{}           `json:"existingPriceAlert"`
	SearchFingerprint  string                `json:"searchFingerprint"`
	HasMorePending     bool                  `json:"hasMorePending"`
	PriceAlertsTop     PriceAlertsTopResults `json:"priceAlertsTopResults"`
	ItinerariesCount   int                   `json:"itinerariesCount"`
	MissingProviders   []interface{}         `json:"missingProviders"`
	StatusPerProvider  []StatusPerProvider   `json:"statusPerProvider"`
	TopFiveContainTHC  bool                  `json:"topFiveOriginalItinerariesContainTHC"`
	HasTier1Market     bool                  `json:"hasTier1MarketItineraries"`
	ContextFilters     []ContextFilter       `json:"contextuallyRecommendedFilters"`
	SharedItinerary    interface{}           `json:"sharedItinerary"`
}

type TopFiveResultsBaggageEligibility struct {
	PromptUserToSearchWithBags *int `json:"promptUserToSearchWithBags"`
	SearchIsLongTrip           *int `json:"searchIsLongTrip"`
	SearchIsFamilyTrip         *int `json:"searchIsFamilyTrip"`
	NumberOfBags               *int `json:"numberOfBags"`
}

type Carrier struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

/* Top results (best/cheapest/fastest...) */

type TopResults struct {
	Best                  ItinerarySummary `json:"best"`
	Cheapest              ItinerarySummary `json:"cheapest"`
	Fastest               ItinerarySummary `json:"fastest"`
	SourceTakeoffAsc      ItinerarySummary `json:"sourceTakeoffAsc"`
	DestinationLandingAsc ItinerarySummary `json:"destinationLandingAsc"`
}

type ItinerarySummary struct {
	Typename string      `json:"__typename"`
	Duration int         `json:"duration"`
	Price    PriceSimple `json:"price"`
	ID       string      `json:"id"`
}

type PriceSimple struct {
	Amount string `json:"amount"`
}

type PriceAlertsTopResults struct {
	Best                  TopPrice `json:"best"`
	Cheapest              TopPrice `json:"cheapest"`
	Fastest               TopPrice `json:"fastest"`
	SourceTakeoffAsc      TopPrice `json:"sourceTakeoffAsc"`
	DestinationLandingAsc TopPrice `json:"destinationLandingAsc"`
}

type TopPrice struct {
	Price PriceSimple `json:"price"`
}

type StatusPerProvider struct {
	Provider      ProviderRef `json:"provider"`
	ErrorHappened bool        `json:"errorHappened"`
	ErrorMessage  string      `json:"errorMessage"`
}

type ProviderRef struct {
	ID string `json:"id"`
}

/* Context filters – union-ish; we keep all possible fields nullable */

type ContextFilter struct {
	Typename      string `json:"__typename"`
	Count         *int   `json:"count,omitempty"`
	Ranks         []int  `json:"ranks,omitempty"`
	IsWeekendTrip *bool  `json:"isWeekendTrip,omitempty"`
}

/* ========== ITINERARIES ========== */

type Itinerary struct {
	Typename      string            `json:"__typename"`
	IsItinerary   string            `json:"__isItinerary"`
	ID            string            `json:"id"`
	ShareID       string            `json:"shareId"`
	Price         ItineraryPrice    `json:"price"`
	PriceEur      PriceSimple       `json:"priceEur"`
	Provider      ItineraryProvider `json:"provider"`
	BagsInfo      BagsInfo          `json:"bagsInfo"`
	BookingOpts   BookingOptions    `json:"bookingOptions"`
	TravelHack    TravelHack        `json:"travelHack"`
	LegacyID      string            `json:"legacyId"`
	Sector        Sector            `json:"sector"`
	LastAvailable LastAvailable     `json:"lastAvailable"`
}

type ItineraryPrice struct {
	Amount              string `json:"amount"`
	PriceBeforeDiscount string `json:"priceBeforeDiscount"`
}

type ItineraryProvider struct {
	Name                            string             `json:"name"`
	Code                            string             `json:"code"`
	HasHighProbabilityOfPriceChange bool               `json:"hasHighProbabilityOfPriceChange"`
	ContentProvider                 ContentProviderRef `json:"contentProvider"`
	ID                              string             `json:"id"`
}

type ContentProviderRef struct {
	Code string `json:"code"`
}

/* Bags */

type BagsInfo struct {
	IncludedCheckedBags   int       `json:"includedCheckedBags"`
	IncludedHandBags      int       `json:"includedHandBags"`
	HasNoBaggageSupported bool      `json:"hasNoBaggageSupported"`
	HasNoCheckedBaggage   bool      `json:"hasNoCheckedBaggage"`
	CheckedBagTiers       []BagTier `json:"checkedBagTiers"`
	HandBagTiers          []BagTier `json:"handBagTiers"`
}

type BagTier struct {
	TierPrice PriceSimple `json:"tierPrice"`
	Bags      []Bag       `json:"bags"`
}

type Bag struct {
	Weight Weight `json:"weight"`
}

type Weight struct {
	Value int `json:"value"`
}

/* Booking options */

type BookingOptions struct {
	Edges []BookingEdge `json:"edges"`
}

type BookingEdge struct {
	Node BookingNode `json:"node"`
}

type BookingNode struct {
	Token         string            `json:"token"`
	BookingURL    string            `json:"bookingUrl"`
	TrackingPixel string            `json:"trackingPixel"`
	Provider      ItineraryProvider `json:"itineraryProvider"`
	Price         PriceSimple       `json:"price"`
}

/* Travel hack */

type TravelHack struct {
	IsTrueHiddenCity     bool `json:"isTrueHiddenCity"`
	IsVirtualInterlining bool `json:"isVirtualInterlining"`
	IsThrowawayTicket    bool `json:"isThrowawayTicket"`
}

/* Sector / segments */

type Sector struct {
	ID             string          `json:"id"`
	SectorSegments []SectorSegment `json:"sectorSegments"`
	Duration       int             `json:"duration"`
}

type SectorSegment struct {
	Guarantee interface{} `json:"guarantee"`
	Segment   Segment     `json:"segment"`
	Layover   interface{} `json:"layover"`
}

type Segment struct {
	ID                   string          `json:"id"`
	Source               SegmentEndpoint `json:"source"`
	Destination          SegmentEndpoint `json:"destination"`
	Duration             int             `json:"duration"`
	Type                 string          `json:"type"`
	Code                 string          `json:"code"`
	Carrier              CarrierRef      `json:"carrier"`
	OperatingCarrier     CarrierRef      `json:"operatingCarrier"`
	CabinClass           string          `json:"cabinClass"`
	HiddenDestination    interface{}     `json:"hiddenDestination"`
	ThrowawayDestination interface{}     `json:"throwawayDestination"`
}

type SegmentEndpoint struct {
	LocalTime string  `json:"localTime"`
	UTCTime   string  `json:"utcTime"`
	Station   Station `json:"station"`
}

type Station struct {
	ID       string  `json:"id"`
	LegacyID string  `json:"legacyId"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Type     string  `json:"type"`
	GPS      GPS     `json:"gps"`
	City     City    `json:"city"`
	Country  Country `json:"country"`
}

type GPS struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type City struct {
	LegacyID string `json:"legacyId"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

type Country struct {
	Code string `json:"code"`
	ID   string `json:"id"`
}

type CarrierRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

/* Last availability */

type LastAvailable struct {
	SeatsLeft *int `json:"seatsLeft"`
}
