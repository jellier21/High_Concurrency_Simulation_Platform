package main

import (
	"client/Redis"
	"client/front"
	"client/rabbitmq"
	"client/rabbitmq/MqTools"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
)

func main() {
	print("client")
	client := resty.New()

	serve := gin.Default()
	InitRedis := Redis.InitRedis()
	rabbmqChannel, MqErr := MqTools.InitMQChannel()
	if MqErr != nil {
		log.Printf("初始化RabbitMQ交换机和初始化通道失败，%v", MqErr)
		return
	}
	front.HandleClientStart(serve, InitRedis, client)
	rabbitmq.RabbitMQhandle(serve, InitRedis, client, rabbmqChannel)
	log.Printf("Starting client on %s:8080", "http://localhost")
	err := serve.Run(":8080")
	if err != nil {
		panic(err)
		return
	}

}
