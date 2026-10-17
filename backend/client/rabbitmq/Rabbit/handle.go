package Rabbit

import (
	"client/model"
	"client/rabbitmq/MqTools"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

// 全局单例：整个进程共用一条连接、一个 channel
var (
	pubChan  *amqp091.Channel
	confirms <-chan amqp091.Confirmation
)

// InitProducer ：服务启动时调用一次，建好全局连接和 channel
func InitProducer() error {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("连接RabbitMQ失败: %w", err)
	}

	pubChan, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("创建通道失败: %w", err)
	}

	if err = MqTools.InitRabbitMQ(pubChan); err != nil {
		return fmt.Errorf("初始化RabbitMQ失败: %w", err)
	}

	if err := pubChan.Confirm(false); err != nil {
		return fmt.Errorf("启动发布前确认失败: %w", err)
	}
	confirms = pubChan.NotifyPublish(make(chan amqp091.Confirmation, 1))

	log.Println("RabbitMQ 生产者初始化完成")
	return nil
}

// SendOrderToQueueAndProcessOrder 处理订单：直接用全局 channel 发布，不再每次 Dial
func SendOrderToQueueAndProcessOrder(data model.OrderDB) error {
	if pubChan == nil {
		return fmt.Errorf("RabbitMQ 生产者未初始化")
	}

	body, err := json.Marshal(map[string]interface{}{
		"data":        data,
		"retry_count": 0,
	})
	if err != nil {
		return err
	}

	return MqTools.PublishConfirmed(pubChan, confirms, "order_exchange", "order", body)
}
