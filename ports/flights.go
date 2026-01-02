package ports

import (
	"context"

	"github.com/justinm35/flyctl/domain"
)

type FlightSearcher interface {
	SearchFlights(ctx context.Context, req domain.SearchRequest) ([]domain.FlightOffer, error)
}
