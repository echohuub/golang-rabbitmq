package main

import "golang-rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQRouting("routingExchange", "routing_key_two")
	rabbitmq.ReceiveRouting()
}
