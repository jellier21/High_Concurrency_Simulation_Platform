package order

import (
	"client/Redis"
	"client/Resty"
	"client/model"
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
		var taskData Resty.TaskDB
		taskData, err := Resty.CreateOrderTask(mslNo, int(count), client)
		if err != nil {
			log.Printf("Failed to create order task: %v", err)
			c.JSON(500, gin.H{"message": "Failed to create order task"})
			return
		}
		c.JSON(200, gin.H{
			"message": "Order task created successfully",
			"data": gin.H{
				"uuid":   taskData.Uuid,
				"msl_no": taskData.MslNo,
				"status": taskData.Status,
				"total":  taskData.Total,
			},
		})
	})
	order.GET("/status/task", func(c *gin.Context) {
		var taskData []model.TaskDB
		taskData, err := Redis.GetTaskMessageInRedis(redisC)
		if err != nil {
			log.Println("无法从redis缓存中命中task,正在查找数据库")
			taskData, err = Resty.GetOrderTaskInformationBySQL(client)
			if err != nil {
				log.Printf("Failed to get order task information: %v", err)
				c.JSON(500, gin.H{"message": "获取任务列表失败"})
				return
			}
		}
		c.JSON(200, gin.H{"message": "获取任务列表成功", "data": taskData})
	})
	order.GET("/status/order", func(c *gin.Context) {
		MslNo := c.Query("uuid")
		log.Printf("接收到前端传来的查询活动预备模拟全部订单请求，活动ID为 %s", MslNo)
		var orderData []model.OrderDB
		orderData, err := Redis.GetOrderInformationByMslNo(MslNo, redisC)
		if err != nil {
			log.Printf("Failed to get order information: %v", err)
			c.JSON(500, gin.H{"message": "获取订单信息失败"})
			return
		}
		c.JSON(200, gin.H{"message": "获取订单信息成功", "data": orderData})
	})
}
