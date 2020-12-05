package api

import (
	"api-service/internal/kafka/sync"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type purchaseForm struct {
	AccountID int   `json:"accountID" binding:"required"`
	CartIDs   []int `json:"cartIDs" binding:"required"`
}

func PurchaseFromCarts(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.PurchaseFromCartsLatency.WithLabelValues(httpStatus).Observe(v)
	}))

	requestBody := purchaseForm{}
	if err := c.ShouldBind(&requestBody); err != nil {
		log.Printf("Bind gin context with specified struct occurs error: %v\n", err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	syncKafka := sync.Kafka{}
	syncKafka.Producer = sync.CrateNewSyncProducer()
	defer func() {
		if err := syncKafka.Close(); err != nil {
			log.Printf("Close kafka producer occurs error: %v\n", err)
		}
		timer.ObserveDuration()
	}()

	payload, err := syncKafka.PublishBuyEvent(requestBody.AccountID, requestBody.CartIDs)
	if err != nil {
		log.Printf("Producer send message error %v\n", err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
	}

	responseGin.Response(http.StatusOK, payload)
}
