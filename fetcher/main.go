package fetcher

import (
	"fmt"
	"time"

	"github.com/akerl/stock-exporter/config"

	"github.com/akerl/metrics/metrics"
	"github.com/akerl/timber/v2/log"
)

var logger = log.NewLogger("stock-exporter.fetcher")

// Fetcher defines the ticker fetching engine
type Fetcher struct {
	Interval   int
	Token      string
	Tickers    []string
	MetricFile *metrics.MetricFile
}

// NewFetcher creates a new syslog engine from the given config
func NewFetcher(conf config.Config, mf *metrics.MetricFile) *Fetcher {
	return &Fetcher{
		Interval:   conf.Interval,
		Token:      conf.Token,
		Tickers:    conf.Tickers,
		MetricFile: mf,
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
		mf := metrics.MetricFile{
			metrics.Metric{
				Name:  "last_updated",
				Type:  "gauge",
				Value: fmt.Sprintf("%d", time.Now().Unix()),
			},
		}
		for _, symbol := range f.Tickers {
			mf := append(mf, f.fetchMetric(symbol))
		}
		time.Sleep(time.Duration(f.Interval) * time.Second)
	}
}

func (f *Fetcher) fetchMetric(symbol strin) {

}
