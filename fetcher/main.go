package fetcher

import (
	"fmt"
	"time"

	"github.com/akerl/stock-exporter/config"

	"github.com/akerl/metrics/metrics"
	"github.com/akerl/timber/v2/log"
	yfa "github.com/oscarli916/yahoo-finance-api"
)

var logger = log.NewLogger("stock-exporter.fetcher")

// Fetcher defines the ticker fetching engine
type Fetcher struct {
	Interval int
	Token    string
	Tickers  []string
	Cache    *config.Cache
}

// NewFetcher creates a new syslog engine from the given config
func NewFetcher(conf config.Config, cache *config.Cache) *Fetcher {
	return &Fetcher{
		Interval: conf.Interval,
		Token:    conf.Token,
		Tickers:  conf.Tickers,
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

func (f *Fetcher) fetchMetric(symbol string) (metrics.Metric, error) {
	t := yfa.NewTicker(symbol)
	info, err := t.Info()
	if err != nil {
		return metrics.Metric{}, err
	}
	m := metrics.Metric{
		Name: "ticker",
		Type: "gauge",
		Tags: map[string]string{
			"symbol": symbol,
		},
		Value: info.RegularMarketPrice.Fmt,
	}
	return m, nil
}
