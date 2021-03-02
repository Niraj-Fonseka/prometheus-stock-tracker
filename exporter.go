package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ResponseQuote struct {
	QuoteResponse struct {
		Result []struct {
			Language                          string  `json:"language"`
			Region                            string  `json:"region"`
			QuoteType                         string  `json:"quoteType"`
			QuoteSourceName                   string  `json:"quoteSourceName"`
			Triggerable                       bool    `json:"triggerable"`
			Currency                          string  `json:"currency"`
			Exchange                          string  `json:"exchange"`
			ShortName                         string  `json:"shortName"`
			LongName                          string  `json:"longName"`
			MessageBoardID                    string  `json:"messageBoardId"`
			ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
			ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
			GmtOffSetMilliseconds             int     `json:"gmtOffSetMilliseconds"`
			Market                            string  `json:"market"`
			EsgPopulated                      bool    `json:"esgPopulated"`
			Tradeable                         bool    `json:"tradeable"`
			EarningsTimestamp                 int     `json:"earningsTimestamp"`
			EarningsTimestampStart            int     `json:"earningsTimestampStart"`
			EarningsTimestampEnd              int     `json:"earningsTimestampEnd"`
			TrailingPE                        float64 `json:"trailingPE"`
			EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths"`
			EpsForward                        float64 `json:"epsForward"`
			EpsCurrentYear                    float64 `json:"epsCurrentYear"`
			PriceEpsCurrentYear               float64 `json:"priceEpsCurrentYear"`
			SharesOutstanding                 int     `json:"sharesOutstanding"`
			BookValue                         float64 `json:"bookValue"`
			FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
			FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
			FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
			TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
			TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
			TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
			MarketCap                         int64   `json:"marketCap"`
			ForwardPE                         float64 `json:"forwardPE"`
			PriceToBook                       float64 `json:"priceToBook"`
			SourceInterval                    int     `json:"sourceInterval"`
			ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy"`
			FirstTradeDateMilliseconds        int64   `json:"firstTradeDateMilliseconds"`
			PriceHint                         int     `json:"priceHint"`
			PostMarketChangePercent           float64 `json:"postMarketChangePercent"`
			PostMarketTime                    int     `json:"postMarketTime"`
			PostMarketPrice                   float64 `json:"postMarketPrice"`
			PostMarketChange                  float64 `json:"postMarketChange"`
			RegularMarketChange               float64 `json:"regularMarketChange"`
			RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
			RegularMarketTime                 int     `json:"regularMarketTime"`
			RegularMarketPrice                float64 `json:"regularMarketPrice"`
			RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
			RegularMarketDayRange             string  `json:"regularMarketDayRange"`
			RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
			RegularMarketVolume               int     `json:"regularMarketVolume"`
			RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
			Bid                               float64 `json:"bid"`
			Ask                               float64 `json:"ask"`
			BidSize                           int     `json:"bidSize"`
			AskSize                           int     `json:"askSize"`
			FullExchangeName                  string  `json:"fullExchangeName"`
			FinancialCurrency                 string  `json:"financialCurrency"`
			RegularMarketOpen                 float64 `json:"regularMarketOpen"`
			AverageDailyVolume3Month          int     `json:"averageDailyVolume3Month"`
			AverageDailyVolume10Day           int     `json:"averageDailyVolume10Day"`
			FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
			FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
			FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
			FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
			FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
			FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
			FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
			MarketState                       string  `json:"marketState"`
			DisplayName                       string  `json:"displayName"`
			Symbol                            string  `json:"symbol"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteResponse"`
}

func recordMetrics() {
	go func() {
		for {
			regular, post := fetchStockPrice()
			if post == 0 {
				g.Set(regular)
			} else {
				g.Set(post)
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

func fetchStockPrice() (float64, float64) {
	resp, err := http.Get("https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&symbols=TSLA")
	if err != nil {
		log.Println(err)
		return 0.0, 0.0
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return 0.0, 0.0
	}

	var qresp ResponseQuote
	err = json.Unmarshal(body, &qresp)

	if err != nil {
		log.Println(err)
		return 0.0, 0.0
	}

	return qresp.QuoteResponse.Result[0].RegularMarketPrice, qresp.QuoteResponse.Result[0].PostMarketPrice
}

var (
	g = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "regular_market",
		Help: "regular market value",
	})
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
