package main

import (
	"fmt"
	"golang-rabbitmq/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("test_queue_name")
	rabbitmq.PublishSimple("hello")
	fmt.Println("发送消息成功")
}
