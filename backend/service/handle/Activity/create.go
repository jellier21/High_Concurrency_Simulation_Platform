package Activity

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"service/Tools"
	"service/redisDB"
	"service/sqldb"
	"strconv"
)

func CreateActivityApi(r *gin.Engine, db *sql.DB, redisC *redis.Client) {
	activity := r.Group("/api/activity")
	activity.POST("/create", func(c *gin.Context) {
		StartTime, _ := strconv.ParseInt(c.PostForm("start"), 10, 64)
		EndTime, _ := strconv.ParseInt(c.PostForm("end"), 10, 64)
		start := Tools.CutTimeToMin(StartTime)
		end := Tools.CutTimeToMin(EndTime)
		log.Printf("收到创建活动请求，开始时间: %v, 结束时间: %v", start, end)
		uuid := Tools.CreateOnlyUuid()
		err := sqldb.CreateActivityDB(uuid, start, end, db)
		if err != nil {
			log.Printf("创建活动失败: %v", err)
			c.JSON(500, gin.H{"message": "创建失败"})
			return
		}
		c.JSON(200, gin.H{
			"message": "创建成功",
			"data": gin.H{
				"uid":        uuid,
				"start_time": start,
				"end_time":   end,
			},
		})
		err = redisDB.InsertActivityRedisMessage(uuid, start, end, redisC)
		if err != nil {
			log.Printf("插入活动到Redis失败: %v,不影响后续操作", err)
		}
	})
	activity.GET("/get", func(c *gin.Context) {
		var activity []sqldb.ActivityDB
		activity, err := sqldb.GetActivityDB(db)
		if err != nil {
			log.Printf("获取活动失败: %v", err)
			c.JSON(500, gin.H{"message": "获取失败"})
			return
		}
		c.JSON(200, gin.H{
			"message": "获取成功",
			"data":    activity,
		})
		for _, item := range activity {
			err := redisDB.InsertActivityRedisMessage(item.MslUid, item.StartTime, item.EndTime, redisC)
			if err != nil {
				log.Printf("插入活动到Redis失败: %v,不影响后续操作", err)
			}
		}
	})
}
