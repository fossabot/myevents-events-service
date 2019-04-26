package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

func NewMetricsServer() http.Handler {
	server := http.NewServeMux()
	server.Handle("/metrics", prometheus.Handler())
	return server
}
