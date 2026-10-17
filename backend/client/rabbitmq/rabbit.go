package rabbitmq

import (
	"client/Redis"
	"client/Resty"
	"client/model"
	"client/rabbitmq/Rabbit"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"log"
)

func RabbitMQhandle(r *gin.Engine, client *redis.Client, restyClient *resty.Client) {
	rabbit := r.Group("/api/rabbitmq")
	rabbit.POST("/start", func(c *gin.Context) {
		orderUid := c.PostForm("uuid")
		var OrderData model.OrderDB
		OrderData, err := Redis.GetRedisOrderByUuid(orderUid, client)
		if err != nil {
			log.Printf("Redis未命中Uuid为 %s 的订单信息，正在从数据库中获取订单信息", orderUid)
			OrderData, err = Resty.GetOrderInformationByUuid(orderUid, restyClient)
		}
		Err := Rabbit.SendOrderToQueueAndProcessOrder(OrderData)
		if Err != nil {
			log.Printf("Uuid为 %s 的订单处理失败: %v", orderUid, Err)
			c.JSON(500, gin.H{"error": Err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "订单处理成功",
			"data":    OrderData,
		})
	})
}
