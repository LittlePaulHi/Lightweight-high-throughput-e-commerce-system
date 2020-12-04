package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PurchaseFromCartsLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "purchase_from_carts_duration_seconds",
		Help:    "Latency of purchase_from_carts request in second.",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"status"},
)


func init() {
	prometheus.MustRegister(PurchaseFromCartsLatency)
}