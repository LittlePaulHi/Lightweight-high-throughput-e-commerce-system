package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var GetAllProductsLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "get_all_products_duration_seconds",
		Help:    "Latency of get_all_products request in second.",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"status"},
)


func init() {
	prometheus.MustRegister(GetAllProductsLatency)
}