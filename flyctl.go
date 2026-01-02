package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func oldmain() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "search-flight",
				Usage: "Search for a specific flight",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "s",
						Usage: "Source Iata code",
					},
					&cli.StringFlag{
						Name:  "d",
						Usage: "Destination Iata code",
					},
					&cli.StringFlag{
						Name:  "dd",
						Usage: "Date of Departure",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Flight Searching....")

					// departDate, err := time.Parse("2006-01-02", cmd.String("dd"))

					// if err != nil {
					// 	log.Fatal(err)
					// }
					// _ := domain.SearchRequest{
					// 	Origin:      cmd.String("s"),
					// 	Destination: cmd.String("d"),
					// 	DepartDate:  departDate,
					// 	Adults:      1,
					// 	MaxResults:  5,
					// }

					// offers, err := rapidgoogleflights.SearchFlights(ctx, req)
					// // offers, err := amadeus.SearchFlights(ctx, req)
					// if err != nil {
					// 	log.Fatal(err)
					// }
					// RenderFlightSearchTable(offers)
					return nil
				},
			},
			{
				Name:  "ap-code",
				Usage: "Search for a airport and get the airport code back",
				Action: func(context.Context, *cli.Command) error {
					fmt.Println("Getting Airport Code....")
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
