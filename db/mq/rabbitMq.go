package mq

import (
	"feiyu.com/wx/srv/srvconfig"
	"github.com/lunny/log"
	"github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	"sync"
)

var (
	publisher rabbitmq.Publisher
)

var once = sync.Once{}

// 初始化 参数格式：amqp://用户名:密码@地址:端口号/host
func SetupRMQ() (err error) {
	//是否开启mq
	if !srvconfig.GlobalSetting.RabbitMq {
		return
	}
	publisher, _, err = rabbitmq.NewPublisher(srvconfig.GlobalSetting.RabbitMqUrl, amqp.Config{})
	if err != nil {
		return err
	}
	log.Println("connect rabbitMq  is  success")
	return nil
}

// 发布消息
func PublishRabbitMq(topic string, body []byte) (err error) {
	once.Do(func() {
		err = SetupRMQ()
		if err != nil {
			log.Fatal(err)
		}
	})
	err = publisher.Publish(body, []string{topic})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Debug("rabbitMq send success")
	return nil
}
