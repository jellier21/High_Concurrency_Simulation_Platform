package business

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"service/rabbitmq/InitMQ"
)

func RunDeadConsumer() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := InitMQ.InitRabbitMQ(ch); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		"dead_queue", "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	log.Println("死信队列消费者启动，等待最终失败的消息……")

	for msg := range msgs {
		log.Printf("收到死信: %s", string(msg.Body))
		if err := msg.Ack(false); err != nil {
			return err
		}
	}

	return errors.New("死信队列消息通道已关闭")
}
