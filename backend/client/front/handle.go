package front

import (
	"client/front/activity"
	"client/front/order"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
)

func HandleClientStart(r *gin.Engine, redisC *redis.Client, client *resty.Client) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, client!")
	})
	activity.ClientActivityhandle(r, redisC, client)
	order.ClientOrderhandle(r, redisC, client)
}
