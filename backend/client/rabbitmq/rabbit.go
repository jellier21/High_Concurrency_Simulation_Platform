package rabbitmq

import (
	"client/Redis"
	"client/model"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RabbitMQhandle(r *gin.Engine, client *redis.Client) {
	rabbit := r.Group("/api/rabbitmq")
	rabbit.POST("/start", func(c *gin.Context) {
		orderUid := c.PostForm("uuid")
		var OrderData model.OrderDB
		OrderData, err := Redis.GetRedisOrderByUuid(orderUid, client)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get order data from Redis"})
			return
		}
		c.JSON(200, gin.H{"message": "Order data retrieved successfully", "data": OrderData})
	})
}
