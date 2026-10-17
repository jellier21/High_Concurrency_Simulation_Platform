package business

import (
	"database/sql"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"log"
	"service/rabbitmq/InitMQ"
)

func RunConsumer(db *sql.DB, client *redis.Client) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	// 消费者通道：负责接收消息和 ACK。
	consumerCh, err := conn.Channel()
	if err != nil {
		return err
	}
	defer consumerCh.Close()

	if err := InitMQ.InitRabbitMQ(consumerCh); err != nil {
		return err
	}

	// 发布通道：专门负责重试/死信消息的发布确认。
	pubCh, err := conn.Channel()
	if err != nil {
		return err
	}
	defer pubCh.Close()

	if err := pubCh.Confirm(false); err != nil {
		return err
	}
	confirms := pubCh.NotifyPublish(make(chan amqp.Confirmation, 1))

	// autoAck=false：必须由业务代码手动 ACK。
	msgs, err := consumerCh.Consume(
		"order_queue", "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	log.Println("订单消费者启动，正在等待消息……")

	for msg := range msgs {
		if err := handleOrderMessage(pubCh, confirms, msg, db, client); err != nil {
			// 退出函数会关闭连接，使未 ACK 的消息可以重新投递。
			return fmt.Errorf("消费流程异常: %w", err)
		}
	}

	return errors.New("订单消费者消息通道已关闭")
}
