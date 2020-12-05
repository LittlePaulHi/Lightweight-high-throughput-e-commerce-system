package api

import (
	"api-service/metrics"
	"api-service/service"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
)

type cartForm struct {
	CartID    int `json:"cartID" binding:"required"`
	AccountID int `json:"accountID" binding:"required"`
	ProductID int `json:"productID" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

// GetAllCartsByAccountID API
// @Param {body: { cartID, accountID, productID }}
// @Router /api/cart/getAllByAccountID [GET]
// @Success 200
// @Failure 500
func GetAllCartsByAccountID(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllCartsByAccIDLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	accIDStr := c.Query("AccountID")
	if accIDStr == "" {
		logger.APILog.Warnln("AccountID shouldn't be empty in GetAllCartsByAccountID")
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	accID, err := strconv.Atoi(accIDStr)
	if err != nil {
		logger.APILog.Warnln("GetAllCartsByAccountID: ", err)
		httpStatus = "InternalServerError"
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	// access cache first
	carts := redisCartCache.GetAllCartsByAcctID(accID)
	// cache miss
	if carts == nil || len(carts) == 0 {
		carts, err = service.GetAllCartsByAccountID(accID)
		if err != nil {
			httpStatus = "InternalServerError"
			responseGin.Response(http.StatusInternalServerError, nil)
			return
		}

		redisCartCache.SetAllCartsByAcctID(accID, carts)
	}

	data := make(map[string]interface{})
	data["cart"] = carts
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)
}

// AddCart API
// @Param {body: { cartID, accountID, productID }}
// @Router /api/cart/addCart [POST]
// @Success 200
// @Failure 500
func AddCart(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.AddCartLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestBody := cartForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	cart := mariadb.Cart{}
	cart.Initialize(
		requestBody.AccountID, requestBody.ProductID, requestBody.Quantity,
	)

	savedCart, err := service.AddCart(&cart)
	if err != nil {
		httpStatus = "InternalServerError"
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["cart"] = savedCart
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)
}

// EditCart API
// @Router /api/cart/editCart [POST]
// @Success 200
// @Failure 500
func EditCart(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.EditCartLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestBody := cartForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	carts, err := service.EditCart(
		requestBody.CartID, requestBody.AccountID, requestBody.ProductID, requestBody.Quantity,
	)
	if err != nil {
		httpStatus = "InternalServerError"
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["carts"] = carts
	data["timestamp"] = time.Now()

	httpStatus = "OK"
	responseGin.Response(http.StatusOK, data)
}
