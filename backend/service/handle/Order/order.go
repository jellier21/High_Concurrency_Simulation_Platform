package Order

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"service/handle/Order/task"
)

func OrderHandle(r *gin.Engine, db *sql.DB, redisC *redis.Client) {
	task.OrderTaskService(r, db, redisC)
}
