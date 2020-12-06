package router

import (
	"github.com/gin-gonic/gin"

	"api-service/api"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Initialize the gin router
func Initialize() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	productAPI := r.Group("/api/product")
	{
		productAPI.GET("/getAll", api.GetAllProducts)
	}

	cartAPI := r.Group("/api/cart")
	{
		cartAPI.GET("/getAllByAccountID", api.GetAllCartsByAccountID)
		cartAPI.POST("/addCart", api.AddCart)
		cartAPI.POST("/editCart", api.EditCart)
	}

	orderAPI := r.Group("/api/order")
	{
		orderAPI.GET("/getAllByAccountID", api.GetAllOrdersByAccountID)
		orderAPI.GET("/getAllItemsByOrderID", api.GetAllOrderItemsByOrderID)
	}

	purchaseAPI := r.Group("/api/purchase")
	{
		purchaseAPI.POST("/sync", api.PurchaseFromCarts)
	}

	//for promethus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
