package sync

import (
	"encoding/json"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/kafka/model"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/redis"
	"log"

	"github.com/Shopify/sarama"
)

var (
	redisClient redis.Redis
)

type Consumer struct {
	Ready chan bool
}

func init() {
	redisClient.Initialize()
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s",
			string(message.Value), message.Timestamp, message.Topic)

		purchaseMsg := &model.PurchaseMessage{}
		err := json.Unmarshal(message.Value, purchaseMsg)
		if err != nil {
			log.Printf("Convert purchase message struct to bytes occurs error: %v\n", err)
			continue
		}

		// query `Cart` table and update the `Product` table
		carts, err := mariadb.FindAllCartsByCardIDs(purchaseMsg.CartIDs)
		if err != nil {
			log.Printf("Get carts in consumer occurs error: %v\n", err)
			continue
		}

		newOrder := mariadb.Order{}
		newOrder.Initialize(purchaseMsg.AccountID, 0)
		if _, err = newOrder.SaveOrder(); err != nil {
			log.Printf("Create order occurs error: %v\n", err)
			continue
		}

		status, err := updateTables(carts, newOrder)
		if err != nil {
			log.Printf("Update purchase related tables occurs error: %v\n", err)
			continue
		}

		if err = publishToRedis(purchaseMsg.RedisChannel, status); err != nil {
			log.Printf("Publish the results to Redis occurs error: %v\n", err)
			continue
		}

		session.MarkMessage(message, "")
	}

	return nil
}

func updateTables(carts []*mariadb.Cart, newOrder mariadb.Order) (string, error) {
	var amounts = 0
	for _, cart := range carts {
		product, err := mariadb.FindProductByID(cart.ProductID)
		if err != nil {
			log.Panicf("Get product in consumer occurs error: %v\n", err)
			return "", err
		}

		if product.Quantity >= cart.Quantity {
			product.Quantity -= cart.Quantity
			if _, err := product.UpdateProduct(); err != nil {
				log.Printf("Purchase the prouct %v occurs error: %v", product.Name, err)
				continue
			}

			orderItem := mariadb.OrderItem{}
			orderItem.Initialize(newOrder.ID, cart.ProductID, cart.Quantity)
			if _, err = orderItem.SaveOrderItem(); err != nil {
				log.Printf("Save the orderItem with product %v occurs error: %v", product.Name, err)
				continue
			}

			amounts += product.Price * cart.Quantity
		} else {
			log.Printf("Product %v already sold out", product.Name)
		}
	}

	if amounts == 0 {
		newOrder.Status = "fail"
	} else {
		newOrder.Amount = amounts
	}

	if _, err := newOrder.UpdateOrder(); err != nil {
		log.Panicf("Update order %v occurs error: %v\n", newOrder.ID, err)
		return "", err
	}

	return newOrder.Status, nil
}

func publishToRedis(channel string, message string) error {
	if err := redisClient.Publish(channel, message); err != nil {
		return err
	}

	return nil
}
