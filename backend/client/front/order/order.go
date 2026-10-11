package order

import (
	"client/Resty"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

func ClientOrderhandle(r *gin.Engine, redisC *redis.Client, client *resty.Client) {
	order := r.Group("/api/order")
	order.POST("/create", func(c *gin.Context) {
		mslNo := c.PostForm("msl_no")
		count, _ := strconv.ParseInt(c.PostForm("count"), 10, 64)
		log.Printf("接收到前端传来的创建 %d 条订单请求", count)
		err := Resty.CreateOrderTask(mslNo, int(count), client)
		if err != nil {
			log.Printf("Failed to create order task: %v", err)
			c.JSON(500, gin.H{"message": "Failed to create order task"})
			return
		}
		c.JSON(200, gin.H{"message": "Order task created successfully"})
	})
}
