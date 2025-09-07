package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/akerl/stock-exporter/config"

	"github.com/akerl/metrics/metrics"
	"github.com/akerl/timber/v2/log"
)

var logger = log.NewLogger("stock-exporter.server")

// Server defines a Prometheus-compatible metrics engine
type Server struct {
	Port       int
	MetricFile *metrics.MetricFile
}

// NewServer creates a new Server object
func NewServer(conf config.Config, mf *metrics.MetricFile) *Server {
	return &Server{
		Port:       conf.Port,
		MetricFile: mf,
	}
}

// Run starts the Server object in the foreground
func (s *Server) Run() error {
	bindStr := fmt.Sprintf(":%d", s.Port)
	logger.InfoMsgf("binding metrics server to %s", bindStr)
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", s.handleMetrics)
	return http.ListenAndServe(bindStr, mux)
}

func (s *Server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	if !s.MetricFile.Validate() {
		logger.InfoMsg("invalid metrics file requested")
		http.Error(w, "invalid metrics file", http.StatusInternalServerError)
	} else {
		logger.InfoMsg("successful metrics request")
		io.WriteString(w, s.MetricFile.String())
	}
}
