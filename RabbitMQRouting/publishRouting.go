package main

import (
	"fmt"
	"golang-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := RabbitMQ.NewRabbitMQRouting("routingExchange", "routing_key_one")
	rabbitmqTwo := RabbitMQ.NewRabbitMQRouting("routingExchange", "routing_key_two")
	for i := 0; i <= 100; i++ {
		rabbitmqOne.PublishRouting("hello one" + strconv.Itoa(i))
		rabbitmqTwo.PublishRouting("hello two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
