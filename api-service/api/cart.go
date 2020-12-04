package api

import (
	"log"
	"net/http"
	"time"

	"api-service/service"

	"github.com/gin-gonic/gin"

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

	requestBody := cartForm{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("Bind with cart body occurs error: %v", err)
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	accID := requestBody.AccountID

	// access cache first
	carts := redisCartCache.GetAllCartsByAcctID(accID)

	// cache miss
	if carts == nil || len(carts) == 0 {
		carts, err = service.GetAllCartsByAccountID(accID)
		if err != nil {
			responseGin.Response(http.StatusInternalServerError, nil)
			return
		}

		redisCartCache.SetAllCartsByAcctID(accID, carts)
	}

	data := make(map[string]interface{})
	data["cart"] = carts
	data["timestamp"] = time.Now()

	responseGin.Response(http.StatusOK, data)
}

// AddCart API
// @Param {body: { cartID, accountID, productID }}
// @Router /api/cart/addCart [POST]
// @Success 200
// @Failure 500
func AddCart(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	requestBody := cartForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		log.Printf("Bind with cart body occurs error: %v", err)
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	cart := mariadb.Cart{}
	cart.Initialize(
		requestBody.AccountID, requestBody.ProductID, requestBody.Quantity,
	)

	savedCart, err := service.AddCart(&cart)
	if err != nil {
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["cart"] = savedCart
	data["timestamp"] = time.Now()

	responseGin.Response(http.StatusOK, data)
}

// EditCart API
// @Router /api/cart/editCart [POST]
// @Success 200
// @Failure 500
func EditCart(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	requestBody := cartForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		log.Printf("Bind with cart body occurs error: %v", err)
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	carts, err := service.EditCart(
		requestBody.CartID, requestBody.AccountID, requestBody.ProductID, requestBody.Quantity,
	)
	if err != nil {
		responseGin.Response(http.StatusInternalServerError, nil)
		return
	}

	data := make(map[string]interface{})
	data["carts"] = carts
	data["timestamp"] = time.Now()

	responseGin.Response(http.StatusOK, data)
}
