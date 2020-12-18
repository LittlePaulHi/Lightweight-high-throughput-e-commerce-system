package api

import (
	"api-service/metrics"
	"api-service/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// GetAllProducts API
// @Router /api/product/getAll [GET]
// @Success 200
// @Failure 500
func GetAllProducts(c *gin.Context) {
	if MaxGoRoutines != 0 {
		<-GoRoutineSemaPhore
	}

	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllProductsLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	products, err := service.GetAllProducts()
	if err != nil {
		httpStatus = "InternalServerError"
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["products"] = products
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)

	if MaxGoRoutines != 0 {
		GoRoutineSemaPhore<-1
	}
}
