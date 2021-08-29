package mq

import (
	"fmt"
	"gindemo/conf"
	"github.com/streadway/amqp"
	"log"
)

var (
	SendMsgChan = make(chan string, 100)
	RabbitMq    = New(&QueueExchange{
		"test2.rabbit",
		"rabbit.key",
		"amq.direct",
		"direct",
	})
)

func Init() {
	RabbitMq.RegisterReceiver(Consumer)
	RabbitMq.Start()
}

func Close() {
	RabbitMq.mqClose()
}

// 定义RabbitMQ对象
type RabbitMQ struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	queueName    string // 队列名称
	routingKey   string // key名称
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
	receiver     func([]byte) error
}

// 定义队列交换机对象
type QueueExchange struct {
	QuName string // 队列名称
	RtKey  string // key值
	ExName string // 交换机名称
	ExType string // 交换机类型
}

// 定义全局变量,指针类型
var mqConn *amqp.Connection
var mqChan *amqp.Channel

// 链接rabbitMQ
func (r *RabbitMQ) mqConnect() {
	var err error
	config := conf.Configs.RabbitMQConfig
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s/", config.User, config.Pwd, config.Address)
	mqConn, err = amqp.Dial(RabbitUrl)
	r.connection = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		panic(fmt.Sprintf("MQ打开链接失败:%s \n", err))
	}
	mqChan, err = mqConn.Channel()
	r.channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		panic(fmt.Sprintf("MQ打开管道失败:%s \n", err))
	}
}

// 关闭RabbitMQ连接
func (r *RabbitMQ) mqClose() {
	// 先关闭管道,再关闭链接
	err := r.channel.Close()
	if err != nil {
		fmt.Printf("MQ管道关闭失败:%s \n", err)
	}
	err = r.connection.Close()
	if err != nil {
		fmt.Printf("MQ链接关闭失败:%s \n", err)
	}
}

func New(q *QueueExchange) *RabbitMQ {
	return &RabbitMQ{
		queueName:    q.QuName,
		routingKey:   q.RtKey,
		exchangeName: q.ExName,
		exchangeType: q.ExType,
	}
}

func (r *RabbitMQ) Start() {
	r.mqConnect()
	go r.listenProducer()
	go r.listenReceiver(r.receiver)
}

// 发送任务
func (r *RabbitMQ) listenProducer() {
	for Msg := range SendMsgChan {
		// 发送任务消息
		err := r.channel.Publish(r.exchangeName, r.routingKey, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(Msg),
		})
		if err != nil {
			fmt.Printf("MQ任务发送失败:%s \n", err)
		}
	}
}

// 注册接收指定队列指定路由的数据接收者
func (r *RabbitMQ) RegisterReceiver(receiver func([]byte) error) {
	r.receiver = receiver
}

// 监听接收者接收任务
func (r *RabbitMQ) listenReceiver(receiver func([]byte) error) {
	// 处理结束关闭链接
	defer r.mqClose()
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err := r.channel.QueueDeclarePassive(r.queueName, true, false, true, false, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.channel.QueueDeclare(r.queueName, true, false, true, false, nil)
		if err != nil {
			log.Printf("MQ注册队列失败:%s \n", err)
		}
	}
	// 绑定任务
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, false, nil)
	if err != nil {
		log.Printf("绑定队列失败:%s \n", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.channel.Qos(1, 0, true)
	msgList, err := r.channel.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("获取消费通道异常:%s \n", err)
		return
	}
	for msg := range msgList {
		// 处理数据
		err := receiver(msg.Body)
		if err != nil {
			err = msg.Ack(true)
			if err != nil {
				log.Printf("确认消息未完成异常:%s \n", err)
			}
		} else {
			// 确认消息,必须为false
			err = msg.Ack(false)
			if err != nil {
				log.Printf("确认消息完成异常:%s \n", err)
			}
		}
	}
}
