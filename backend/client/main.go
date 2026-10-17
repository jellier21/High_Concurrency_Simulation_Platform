package main

import (
	"client/Redis"
	"client/front"
	"client/rabbitmq"
	"client/rabbitmq/Rabbit"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
)

func main() {
	print("client")
	client := resty.New()

	serve := gin.Default()
	InitRedis := Redis.InitRedis()

	// 启动时只建一次 RabbitMQ 连接
	if err := Rabbit.InitProducer(); err != nil {
		log.Fatalf("RabbitMQ 生产者初始化失败: %v", err)
	}

	front.HandleClientStart(serve, InitRedis, client)
	rabbitmq.RabbitMQhandle(serve, InitRedis, client)
	log.Printf("Starting client on %s:8080", "http://localhost")
	err := serve.Run(":8080")
	if err != nil {
		panic(err)
		return
	}

}
