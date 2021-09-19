package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"imserver/conf"
	"log"
	"time"
)

var (
	MQ *RabbitMQ
)

type Consumer func(d *amqp.Delivery)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// Key
	Key string
	// 连接信息
	Mqurl string
	//
	Consumer Consumer
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName, exchange, key string, consumer Consumer) *RabbitMQ {
	config := conf.Configs.RabbitMQConfig
	mqUrl := fmt.Sprintf("amqp://%s:%s@%s/", config.User, config.Pwd, config.Address)
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     mqUrl,
		Consumer:  consumer,
	}
	return rabbitmq
}

func (r *RabbitMQ) connectionMQ() {
	var err error
	// 创建rabbitmq连接
	r.conn, err = amqp.Dial(r.Mqurl)
	r.failOnErr(err, "创建连接错误！")

	r.channel, err = r.conn.Channel()
	r.failOnErr(err, "获取channel失败！")
}

func (r *RabbitMQ) Start() {
	for {
		r.connectionMQ()
		r.createExchange()
		r.createQueue()
		r.bindQueue()
		r.registerConsumer()
		//10秒重连
		time.Sleep(time.Second * 10)
	}
}

func (r *RabbitMQ) createExchange() {
	// 声明 exchange
	if err := r.channel.ExchangeDeclare(
		r.Exchange, //name
		"topic",    //exchangeType
		true,       //durable
		false,      //auto-deleted
		false,      //internal
		false,      //noWait
		nil,        //arguments
	); err != nil {
		r.failOnErr(err, "Failed to declare a exchange:")
	}
}
func (r *RabbitMQ) createQueue() {
	// 声明一个queue
	if _, err := r.channel.QueueDeclare(
		r.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	); err != nil {
		r.failOnErr(err, "Failed to declare a queue:")
	}
}

func (r *RabbitMQ) bindQueue() {
	err := r.channel.QueueBind(r.QueueName, r.Key, r.Exchange, false, nil)
	if err != nil {
		r.failOnErr(err, "bind queue error :")
	}
}

func (r *RabbitMQ) SendMessage(data []byte, key string) (err error) {
	// 发送
	if err = r.channel.Publish(
		r.Exchange, // exchange
		key,        // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            data,
			//Expiration:      "60000", // 消息过期时间
		},
	); err != nil {
		log.Println("Failed to publish a message:", err.Error())
		return err
	}
	return nil
}

func (r *RabbitMQ) registerConsumer() error {
	// 注册消费者
	msgs, err := r.channel.Consume(
		r.QueueName, // queue
		"project",   // 标签
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Println("Failed to register a consumer:", err.Error())
		return err
	}
	for d := range msgs {
		r.Consumer(&d)
	}
	return nil
}

// failOnErr 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
