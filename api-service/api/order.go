package api

import (
	"net/http"
	"time"

	"api-service/metrics"
	"api-service/service"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// GetAllOrdersByAccountID API
// @Param {body: { accountID, orderID }}
// @Router /api/order/getAllByAccountID [GET]
// @Success 200
// @Failure 500

type GetAllOrdersByAccountIDRequestHeader struct {
	AccountID int `header:"accountID" binding:"required"`
}

func GetAllOrdersByAccountID(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllOrdersByAccIDLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestHeader := GetAllOrdersByAccountIDRequestHeader{}
	err := c.ShouldBindHeader(&requestHeader)
	if err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	accountID := requestHeader.AccountID
	orders, err := service.GetAllOrdersByAccountID(accountID)
	if err != nil {
		httpStatus = "InternalServerErro"
		responseGin.Response(http.StatusInternalServerError, nil)
		return
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

type GetAllOrderItemsByOrderIDRequestHeader struct {
	OrderID int `header:"orderID" binding:"required"`
}

func GetAllOrderItemsByOrderID(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllOrderItemsByOrderIDLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestHeader := GetAllOrderItemsByOrderIDRequestHeader{}
	err := c.ShouldBindHeader(&requestHeader)
	if err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	orderID := requestHeader.OrderID
	orderItems, err := service.GetAllOrderItemsByOrderID(orderID)
	if err != nil {
		httpStatus = "InternalServerError"
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["orderItems"] = orderItems
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)
}
