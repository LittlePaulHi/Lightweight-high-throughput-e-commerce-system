package api

import (
	"net/http"
	"time"

	"api-service/metrics"
	"api-service/service"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type orderForm struct {
	AccountID int `json:"accountID" binding:"required"`
}

type orderItemForm struct {
	OrderID int `json:"orderID" binding:"required"`
}

// GetAllOrdersByAccountID API
// @Param {body: { accountID, orderID }}
// @Router /api/order/getAllByAccountID [GET]
// @Success 200
// @Failure 500
func GetAllOrdersByAccountID(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllOrdersByAccIDLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestBody := orderForm{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	accID := requestBody.AccountID

	//access cache first
	orders := redisOrderCache.GetAllOrdersByAcctID(accID)

	//cache miss
	if orders == nil || len(orders) == 0 {
		orders, err = service.GetAllOrdersByAccountID(accID)
		if err != nil {
			httpStatus = "InternalServerErro"
			responseGin.Response(http.StatusInternalServerError, nil)
			return
		}

		redisOrderCache.SetAllOrdersByAcctID(accID, orders)
	}

	data := make(map[string]interface{})
	data["orders"] = orders
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)
}

// GetAllOrderItemsByOrderID API
// @Param {body: { accountID, orderID }}
// @Router /api/order/getAllOrderItems [GET]
// @Success 200
// @Failure 500
func GetAllOrderItemsByOrderID(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllOrderItemsByOrderIDLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestBody := orderItemForm{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	orderID := requestBody.OrderID
	// access cache first
	orderItems := redisOrderCache.GetAllOrderItemsByOrderID(orderID)

	//cache miss
	if orderItems == nil || len(orderItems) == 0 {
		orderItems, err = service.GetAllOrderItemsByOrderID(orderID)
		if err != nil {
			httpStatus = "InternalServerError"
			responseGin.Response(http.StatusInternalServerError, nil)
			return
		}

		redisOrderCache.SetAllOrderItemsByOrderID(orderID, orderItems)
	}

	data := make(map[string]interface{})
	data["orderItems"] = orderItems
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)
}
