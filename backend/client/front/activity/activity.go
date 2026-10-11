package activity

import (
	"client/Redis"
	"client/Resty"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"log"
)

func ClientActivityhandle(r *gin.Engine, redisC *redis.Client, client *resty.Client) {
	activity := r.Group("/api/activity")
	activity.POST("/create", func(c *gin.Context) {
		start := c.PostForm("start")
		end := c.PostForm("end")
		log.Printf("Received start: %s, end: %s", start, end)
		err := Resty.CreateActivityMessageInService(start, end, client)
		if err != nil {
			log.Printf("创建秒杀活动失败: %v", err)
			c.JSON(500, gin.H{"message": "创建失败"})
			return
		}
		c.JSON(200, gin.H{"message": "创建成功"})
	})
	activity.GET("/get", func(c *gin.Context) {
		activity, err := Redis.GetAllActivity(redisC)
		if err != nil || len(activity) == 0 {
			log.Printf("获取缓存中全部秒杀活动信息失败或者为空: %v， 正在从数据库获取", err)
			activity, err = Resty.GetActivityMessageInService(client)
			if err != nil {
				log.Printf("获取数据库中全部秒杀活动信息失败: %v", err)
				c.JSON(500, gin.H{"message": "获取失败"})
				return
			}
			// 回源成功 → 全量刷新缓存，修复缓存缺失
			if rerr := Redis.RefreshActivity(activity, redisC); rerr != nil {
				log.Printf("刷新Redis缓存失败: %v", rerr)
			}
		}
		c.JSON(200, gin.H{"message": "获取成功", "data": activity})
	})

}
