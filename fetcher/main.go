package fetcher

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/akerl/stock-exporter/config"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/akerl/metrics/metrics"
	"github.com/akerl/timber/v2/log"
)

var logger = log.NewLogger("stock-exporter.fetcher")

// Fetcher defines the ticker fetching engine
type Fetcher struct {
	Interval int
	Tickers  []string
	Cache    *config.Cache
	Token    string
	client   *finnhub.DefaultApiService
}

// NewFetcher creates a new syslog engine from the given config
func NewFetcher(conf config.Config, cache *config.Cache) *Fetcher {
	return &Fetcher{
		Interval: conf.Interval,
		Tickers:  conf.Tickers,
		Token:    conf.Token,
		Cache:    cache,
	}
}

// RunAsync launches the fetcher engine in the background
func (f *Fetcher) RunAsync() {
	go f.Run()
}

// Run launches the fetcher engine in the foreground
func (f *Fetcher) Run() {
	for {
		logger.DebugMsg("running fetcher loop")
		ms := metrics.MetricSet{
			metrics.Metric{
				Name:  "last_updated",
				Type:  "gauge",
				Value: fmt.Sprintf("%d", time.Now().Unix()),
			},
		}
		for _, symbol := range f.Tickers {
			m, err := f.fetchMetric(symbol)
			if err != nil {
				logger.InfoMsgf("failed fetching %s: %s", symbol, err)
				continue
			}
			ms = append(ms, m)
		}
		f.Cache.MetricSet = ms
		time.Sleep(time.Duration(f.Interval) * time.Second)
	}
}

func (f *Fetcher) newClient() *finnhub.DefaultApiService {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", f.Token)
	return finnhub.NewAPIClient(cfg).DefaultApi
}

func (f *Fetcher) fetchMetric(symbol string) (metrics.Metric, error) {
	if f.client == nil {
		f.client = f.newClient()
	}

	res, _, err := f.client.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return metrics.Metric{}, err
	}

	m := metrics.Metric{
		Name: "ticker",
		Type: "gauge",
		Tags: map[string]string{
			"symbol": symbol,
		},
		Value: strconv.FormatFloat(float64(*res.C), 'f', 2, 32),
	}
	return m, nil
}
