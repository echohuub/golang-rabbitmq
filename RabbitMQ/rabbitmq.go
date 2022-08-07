package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// MQURL amqp://账号:密码@服务器地址:端口号/vhost
const MQURL = "amqp:testuser:123456@127.0.0.1:5672/test"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// Key
	Key   string
	MQUrl string
}

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MQUrl: MQURL}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MQUrl)
	rabbitmq.failOnErr(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取Chanel失败")
	return rabbitmq
}

func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
		panic(fmt.Sprintf("%s: %s", message, err))
	}
}

// 创建简单模式下的RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	return rabbitmq
}

// 简单模式下发消息
func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列，如果队列不存在会自动创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排它性
		false, // 是否阻塞
		nil,   // 额外属性
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 简单模式下接收消息
func (r *RabbitMQ) ConsumeSimple() {
	// 1.申请队列，如果队列不存在会自动创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排它性
		false, // 是否阻塞
		nil,   // 额外属性
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2.接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}

// 创建订阅模式RabbitMQ实例
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchange, "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MQUrl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 订阅模式消费
func (r *RabbitMQ) ReceiveSub() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2.尝试创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}

// 创建路由模式RabbitMQ实例
func NewRabbitMQRouting(exchange string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchange, routingKey)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MQUrl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 路由模式接收消息
func (r *RabbitMQ) ReceiveRouting() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2.尝试创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
