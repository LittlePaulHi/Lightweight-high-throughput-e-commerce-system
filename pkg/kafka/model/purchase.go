package model

type PurchaseMessage struct {
	RedisChannel string `json:"redisChannel"`
	AccountID    int    `json:"accountID"`
	CartIDs      []int  `json:"cartIDs"`
}
