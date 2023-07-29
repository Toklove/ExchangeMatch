package engine

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 列表名
	Exchange  string // 交换机
	Key       string
	MQUrl     string
}

func NewRabbitMq(queueName, exchange, key string) *RabbitMQ {
	host := Conf.MqConf.Host
	port := Conf.MqConf.Port
	username := Conf.MqConf.Username
	password := Conf.MqConf.Password
	// MQ URL 格式 amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
	MQUrl := "amqp://" + username + ":" + password + "@" + host + ":" + port + "/"

	rabbitMq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MQUrl:     MQUrl,
	}
	var err error
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MQUrl)
	rabbitMq.failOnErr(err, "连接MQ错误")

	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.failOnErr(err, "获取channel失败")

	return rabbitMq
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic("MQ 连接错误")
	}
}

func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func NewSimpleRabbitMQ(queueName string) *RabbitMQ {
	return NewRabbitMq(queueName, "", "")
}

func (r *RabbitMQ) PublishNewOrder(message []byte) {
	//1. 申请队列，如果队列不存在会自动创建，如何存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		true,  // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性？
		false, // 是否阻塞
		nil,   // 其它属性
	)
	if err != nil {
		fmt.Println(err)
	}

	err = r.channel.PublishWithContext(context.Background(),
		r.Exchange,
		r.QueueName,
		false, // 如果为true, 会根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false, // 如果为true, 当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息发还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *RabbitMQ) ConsumeNewOrder() {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		true,  // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性？
		false, // 是否阻塞
		nil,   // 其它属性
	)
	if err != nil {
		fmt.Println(err)
	}

	messages, err := r.channel.Consume(
		r.QueueName,
		"",    // 用来区分多个消费者
		true,  // 是否自动应答
		false, // 是否具有排他性
		false, // 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, // 队列消费是否阻塞
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	if Debug {
		log.Printf("[*] Waiting for message, To exit press CTRL+C")
	}
	// 启用协程处理消息
	//go func() {
	for d := range messages {
		order := OrderNode{}
		err := json.Unmarshal(d.Body, &order)
		if err != nil {
			fmt.Println(err)
		}
		if Debug {
			fmt.Printf("来源数据----------:%#v\n", order)
		}
		DoOrder(order)
	}

	//}()
	if Debug {
		log.Printf("[*] Waiting for message, To exit press CTRL+C")
	}
	<-forever
}

func (r *RabbitMQ) ConsumeMatchOrder() {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		true,  // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性？
		false, // 是否阻塞
		nil,   // 其它属性
	)
	if err != nil {
		fmt.Println(err)
	}

	messages, err := r.channel.Consume(
		r.QueueName,
		"",    // 用来区分多个消费者
		true,  // 是否自动应答
		false, // 是否具有排他性
		false, // 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, // 队列消费是否阻塞
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	if Debug {
		log.Printf("[*] Waiting for message, To exit press CTRL+C")
	}
	// 启用协程处理消息
	//go func() {
	for d := range messages {
		order := MatchResult{}
		err := json.Unmarshal(d.Body, &order)
		if err != nil {
			fmt.Println(err)
		}
		if LogLevel == "debug" {
			log.Printf("撮合结果------：%#v\n", order)
		}
	}

	//}()
	if Debug {
		log.Printf("[*] Waiting for message, To exit press CTRL+C")
	}
	<-forever
}
