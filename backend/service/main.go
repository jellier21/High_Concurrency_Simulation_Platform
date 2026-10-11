package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"service/handle"
	"service/redisDB"
	"service/sqldb"
)

func main() {
	serve := gin.Default()
	db := sqldb.InitMySQL()
	redis := redisDB.InitRedis()
	handle.HandleClientStart(serve, db, redis)
	log.Printf("服务端启动在： %s:8081", "http://localhost")
	err := serve.Run(":8081")
	if err != nil {
		panic(err)
		return
	}
}
