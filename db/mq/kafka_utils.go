package mq

import (
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/lunny/log"
	"strings"
)

// 构建一个结构体，用来实例化单例
type ProducerBeam struct {
	Producer sarama.SyncProducer
}

// 声明一个私有变量，作为单例
var producerModel *ProducerBeam

// init函数将在包初始化时执行，实例化单例
func InitKafKa() {
	//是否开启kafka
	if !srvconfig.GlobalSetting.Kafka {
		return
	}
	producerModel = new(ProducerBeam)
	ProducerCreate, _ := createKafka()
	producerModel.Producer = ProducerCreate
}

// 获取单例
func GetProducer() *ProducerBeam {
	return producerModel
}

// 初始化Kafka
func createKafka() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true
	//设备认证方式
	if srvconfig.GlobalSetting.KafkaUsername != "" && srvconfig.GlobalSetting.KafkaPassword != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = srvconfig.GlobalSetting.KafkaUsername
		config.Net.SASL.Password = srvconfig.GlobalSetting.KafkaPassword
	}
	// 使用给定代理地址和配置创建一个同步生产者
	url := srvconfig.GlobalSetting.KafkaUrl
	if url == "" {
		log.Error("kafka url 不能为空!")
		panic("kafka url 不能为空!")
	}
	v := strings.Split(url, ",")
	Producer, err := sarama.NewSyncProducer(v, config)
	if err != nil {
		panic(fmt.Sprintf("producer mq start err:%v", err))
		return nil, err
	}
	log.Info("KafKa init success")
	return Producer, nil
}

// 发送KafKa对象
func SendKafKaMsg(topic string, val []byte) {
	//是否开启kafka
	if !srvconfig.GlobalSetting.Kafka {
		return
	}
	//创建kafka
	Producer := GetProducer().Producer
	//构建发送的消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(val) // sarama.StringEncoder("this is a test log 1111")
	// 发送消息
	_, _, err := Producer.SendMessage(msg)
	if err != nil {
		log.Error("send msg failed, err::%v", err.Error())
		return
	}
	//log.Printf("pid:%v offset:%v\n", pid, offset)
	//defer Producer.Close()
}
