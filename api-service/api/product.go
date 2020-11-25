package api

import (
	"net/http"
	"time"

	"api-service/service"
	"github.com/gin-gonic/gin"
)

// GetAllProducts API
// @Router /api/product/getAll [GET]
// @Success 200
// @Failure 500
func GetAllProducts(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	// access cache first 
	products := redisProductCache.GetAllProducts()

	// cache miss
	if products == nil || len(products) == 0 {
		var err error
		products, err = service.GetAllProducts()
		if err != nil {
			responseGin.Response(http.StatusInternalServerError, nil)
			return
		}

		redisProductCache.SetAllProducts(products)
	}

	data := make(map[string]interface{})
	data["products"] = products
	data["timestamp"] = time.Now()

	responseGin.Response(http.StatusOK, data)
}
