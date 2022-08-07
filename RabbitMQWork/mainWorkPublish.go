package main

import (
	"fmt"
	"golang-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("test_queue_name")
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("hello" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
