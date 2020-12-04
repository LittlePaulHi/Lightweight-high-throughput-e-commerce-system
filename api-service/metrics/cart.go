package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var GetAllCartsByAccIDLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "get_all_carts_by_accID_duration_seconds",
		Help:    "Latency of get_all_carts_by_accID request in second.",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"status"},
)

var AddCartLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "add_cart_duration_seconds", // metric name
		Help: "Latency of add_cart request in second",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"status"}, // labels
)

var EditCartLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "edit_cart_duration_seconds", // metric name
		Help: "Latency of edit_cart request in second.",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"status"}, // labels
)

func init() {
	prometheus.MustRegister(GetAllCartsByAccIDLatency)
	prometheus.MustRegister(AddCartLatency)
	prometheus.MustRegister(EditCartLatency)
}
