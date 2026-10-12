package main

import (
	"client/Redis"
	"client/front"
	"client/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
)

func main() {
	print("client")
	client := resty.New()

	serve := gin.Default()
	InitRedis := Redis.InitRedis()
	front.HandleClientStart(serve, InitRedis, client)
	rabbitmq.RabbitMQhandle(serve, InitRedis, client)
	log.Printf("Starting client on %s:8080", "http://localhost")
	err := serve.Run(":8080")
	if err != nil {
		panic(err)
		return
	}

}
