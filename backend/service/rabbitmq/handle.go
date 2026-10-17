package rabbitmq

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"log"
	"service/rabbitmq/business"
)

func RabbitMQBusinessHandle(db *sql.DB, client *redis.Client) {
	go func() {
		err := business.RunConsumer(db, client)
		if err != nil {
			log.Printf("消费者启动失败: %v", err)
			return
		}
	}()
	go func() {
		err := business.RunDeadConsumer()
		if err != nil {
			log.Printf("死信消费者启动失败: %v", err)
			return
		}
	}()
	select {}
}
