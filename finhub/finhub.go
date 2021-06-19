package main

import (
	"context"
	"fmt"
	"log"
	"os"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

func main() {
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("FINHUBAPIKEY"), // Replace this
	})

	stockSymbols, _, err := finnhubClient.StockSymbols(auth, "US")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total Size  %v\n", len(stockSymbols))

	for _, s := range stockSymbols {
		fmt.Print(s.Symbol + ",")
	}

}
