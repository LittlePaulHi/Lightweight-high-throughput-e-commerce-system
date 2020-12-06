package api

import (
	"api-service/metrics"
	"api-service/service"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
)

// GetAllCartsByAccountID API
// @Param {body: { cartID, accountID, productID }}
// @Router /api/cart/getAllByAccountID [GET]
// @Success 200
// @Failure 500

type GetAllCartsByAccountIDRequestHeader struct {
	AccountID int `header:"accountID" binding:"required"`
}

func GetAllCartsByAccountID(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.GetAllCartsByAccIDLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestHeader := GetAllCartsByAccountIDRequestHeader{}
	err := c.ShouldBindHeader(&requestHeader)
	if err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	accID := requestHeader.AccountID

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

type AddCartRequestBody struct {
	AccountID int `header:"accountID" binding:"required"`
	ProductID int `header:"productID" binding:"required"`
	Quantity  int `header:"quantity" binding:"required"`
}

func AddCart(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.AddCartLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestBody := AddCartRequestBody{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
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

type EditCartRequestBody struct {
	CartID    int `header:"cartID" binding:"required"`
	AccountID int `header:"accountID" binding:"required"`
	ProductID int `header:"productID" binding:"required"`
	Quantity  int `header:"quantity" binding:"required"`
}

func EditCart(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.EditCartLatency.WithLabelValues(httpStatus).Observe(v)
	}))
	defer timer.ObserveDuration()

	requestBody := EditCartRequestBody{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
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
