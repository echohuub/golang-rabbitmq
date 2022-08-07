package main

import "golang-rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("test_queue_name")
	rabbitmq.ConsumeSimple()
}
