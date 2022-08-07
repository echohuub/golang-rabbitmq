package main

import (
	"fmt"
	"golang-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("newProduct")
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishPub("订阅生产第" + strconv.Itoa(i) + "条数据")
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
