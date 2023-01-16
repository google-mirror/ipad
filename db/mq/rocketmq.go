package mq

import (
	"context"
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/lunny/log"
)

type MqConf struct {
	NameServers []string `mapstructure:"nameServers"`
}

var (
	MqProducer            rocketmq.Producer
	MqPushConsumerSuccess rocketmq.PushConsumer
	MqPushConsumerFail    rocketmq.PushConsumer
	MqPushConsumerDelay   rocketmq.PushConsumer
)

const (
	MqRetryTimes = 3
)

func InitMq() {
	//是否开启mq
	if !srvconfig.GlobalSetting.RocketMq {
		return
	}
	mqConf := &MqConf{
		NameServers: []string{srvconfig.GlobalSetting.RocketMqHost},
	}
	if mqConf == nil {
		panic("mq config is nil")
		return
	}
	var err error
	MqProducer, err = rocketmq.NewProducer(
		producer.WithGroupName("fx-group"),
		producer.WithNameServer(mqConf.NameServers),
		producer.WithRetry(MqRetryTimes),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: srvconfig.GlobalSetting.RocketAccessKey,
			SecretKey: srvconfig.GlobalSetting.RocketSecretKey,
		}),
	)
	if err != nil {
		panic(fmt.Sprintf("init rocket mq producer err:%v", err))
		return
	}

	err = MqProducer.Start()
	if err != nil {
		panic(fmt.Sprintf("producer mq start err:%v", err))
		return
	}
	log.Println("connect rocketMq  is  success")
}

func PushRocketMq(topic string, body []byte) {
	//是否开启mq
	if !srvconfig.GlobalSetting.RocketMq {
		return
	}
	msg := primitive.NewMessage(topic, body)
	MqProducer.SendSync(context.Background(), msg)
}

func ShutDownMq() {
	_ = MqProducer.Shutdown()
	_ = MqPushConsumerSuccess.Shutdown()
	_ = MqPushConsumerFail.Shutdown()
	_ = MqPushConsumerDelay.Shutdown()
}
