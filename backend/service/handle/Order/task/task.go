package task

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"service/handle/Order/order"
	"service/redisDB"
	"service/sqldb"
	"strconv"
)

func OrderTaskService(r *gin.Engine, db *sql.DB, redisC *redis.Client) {
	task := r.Group("/api/task")
	task.POST("/create", func(c *gin.Context) {
		mslNo := c.PostForm("msl_no")
		count, _ := strconv.ParseInt(c.PostForm("count"), 10, 64)
		var taskData sqldb.TaskDB
		taskData, err := sqldb.CreateTaskDB(mslNo, int(count), db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"task": taskData})
		err = redisDB.InsertTaskRedisMessage(taskData, redisC)
		if err != nil {
			log.Println("插入Redis失败，此情况不影响后续程序")
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("订单生成协程 panic: %v", r)
				}
			}()
			order.StartOrderTasking(mslNo, taskData.Uuid, int(count), db, redisC)
		}()

	})
}
