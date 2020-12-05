package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var GetAllOrdersByAccIDLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "get_all_orders_by_accID_duration_seconds",
		Help:    "Latency of get_all_orders_by_accID request in second.",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"status"},
)

var GetAllOrderItemsByOrderIDLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "get_all_orders_items_by_orderID_seconds",
		Help:    "Latency of get_all_orders_items_by_orderID request in second.",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"status"},
)

func init() {
	prometheus.MustRegister(GetAllOrdersByAccIDLatency)
	prometheus.MustRegister(GetAllOrderItemsByOrderID)
}
