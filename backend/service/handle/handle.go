package handle

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"service/handle/Activity"
	"service/handle/Order"
)

func HandleClientStart(r *gin.Engine, db *sql.DB, redisC *redis.Client) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, client!")
	})
	Activity.CreateActivityApi(r, db, redisC)
	Order.OrderHandle(r, db, redisC)
}
