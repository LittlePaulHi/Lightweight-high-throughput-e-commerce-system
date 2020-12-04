package api

import (
	"net/http"
	"time"
	"api-service/metrics"
	"api-service/service"
	"github.com/gin-gonic/gin"
)

// GetAllProducts API
// @Router /api/product/getAll [GET]
// @Success 200
// @Failure 500
func GetAllProducts(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string 
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllProductsLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	// access cache first
	products := redisProductCache.GetAllProducts()

	// cache miss
	if products == nil || len(products) == 0 {
		var err error
		products, err = service.GetAllProducts()
		if err != nil {
			httpStatus = "InternalServerError"
			responseGin.Response(http.StatusInternalServerError, nil)
			return
		}

		redisProductCache.SetAllProducts(products)
	}

	data := make(map[string]interface{})
	data["products"] = products
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)
}
