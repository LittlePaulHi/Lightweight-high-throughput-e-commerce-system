package api

import (
	"api-service/internal/kafka/sync"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type purchaseForm struct {
	AccountID int   `json:"accountID" binding:"required"`
	CartIDs   []int `json:"cartIDs" binding:"required"`
}

func PurchaseFromCarts(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	requestBody := purchaseForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		log.Printf("Bind gin context with specified struct occurs error: %v\n", err)
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	syncKafka := sync.Kafka{}
	syncKafka.Producer = sync.CrateNewSyncProducer()
	defer func() {
		if err := syncKafka.Close(); err != nil {
			log.Printf("Close kafka producer occurs error: %v\n", err)
		}
	}()

	payload, err := syncKafka.PublishBuyEvent(requestBody.AccountID, requestBody.CartIDs)
	if err != nil {
		log.Printf("Producer send message error %v\n", err)
		responseGin.Response(http.StatusBadRequest, nil)
	}

	responseGin.Response(http.StatusOK, payload)
}
