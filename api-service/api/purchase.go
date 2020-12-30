package api

import (
	"api-service/internal/kafka/sync"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PurchaseFromCartsRequestBody struct {
	AccountID int   `json:"accountID" binding:"required"`
	CartIDs   []int `json:"cartIDs" binding:"required"`
}

func PurchaseFromCarts(c *gin.Context) {
	if MaxGoRoutines != 0 {
		<-GoRoutineSemaPhore
	}
	responseGin := ResponseGin{Context: c}

	requestBody := PurchaseFromCartsRequestBody{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
		logger.APILog.Warnln(err)
		responseGin.Response(http.StatusBadRequest, nil)
		return
	}

	var syncKafka sync.Kafka
	syncKafka.Producer = sync.CrateNewSyncProducer()
	defer func() {
		if err := syncKafka.Close(); err != nil {
			logger.KafkaProducer.Warnf("Close syncProducer occurs error: %v\n", err)
		}
	}()

	payload, err := syncKafka.PublishBuyEvent(requestBody.AccountID, requestBody.CartIDs)
	if err != nil {
		logger.KafkaProducer.Warnf("Producer publish message error %v\n", err)

		// TODO: query database when occurs error on producer publish

		responseGin.Response(http.StatusInternalServerError, nil)
	} else {
		responseGin.Response(http.StatusOK, payload)
	}

	if MaxGoRoutines != 0 {
		GoRoutineSemaPhore <- 1
	}
}
