package api

import (
	"api-service/internal/kafka/sync"
	"api-service/metrics"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type PurchaseFromCartsRequestBody struct {
	AccountID int   `json:"accountID" binding:"required"`
	CartIDs   []int `json:"cartIDs" binding:"required"`
}

func PurchaseFromCarts(c *gin.Context) {
	responseGin := ResponseGin{Context: c}

	var httpStatus string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.PurchaseFromCartsLatency.WithLabelValues(httpStatus).Observe(v)
	}))

	requestBody := PurchaseFromCartsRequestBody{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
		logger.APILog.Warnln(err)
		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	syncKafka := sync.Kafka{}
	syncKafka.Producer = sync.CrateNewSyncProducer()
	defer func() {
		if err := syncKafka.Close(); err != nil {
			logger.KafkaProducer.Warnf("Close kafka producer occurs error: %v\n", err)
		}
		timer.ObserveDuration()
	}()

	payload, err := syncKafka.PublishBuyEvent(requestBody.AccountID, requestBody.CartIDs)
	if err != nil {
		logger.KafkaProducer.Warnf("Producer publish message error %v\n", err)

		// TODO: query database when occurs error on producer publish

		httpStatus = "BadRequest"
		responseGin.Response(http.StatusBadRequest, nil)
	}

	responseGin.Response(http.StatusOK, payload)
}
