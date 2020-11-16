package api

import (
	"log"
	"net/http"
	"time"

	"api-service/service"
	"github.com/gin-gonic/gin"
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

	requestBody := orderForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		log.Fatal(err)
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	orders, err := service.GetAllOrdersByAccountID(
		requestBody.AccountID,
	)
	if err != nil {
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["orders"] = orders
	data["timestamp"] = time.Now()

	responseGin.Response(http.StatusOK, data)
}

// GetAllOrderItemsByOrderID API
// @Param {body: { accountID, orderID }}
// @Router /api/order/getAllOrderItems [GET]
// @Success 200
// @Failure 500
func GetAllOrderItemsByOrderID(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	requestBody := orderItemForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		log.Fatal(err)
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	orderItems, err := service.GetAllOrderItemsByOrderID(
		requestBody.OrderID,
	)
	if err != nil {
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["orderItems"] = orderItems
	data["timestamp"] = time.Now()

	responseGin.Response(http.StatusOK, data)
}
