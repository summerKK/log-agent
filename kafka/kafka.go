package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

var client sarama.SyncProducer
var clientChan chan *lodData

type lodData struct {
	topic string
	msg   string
}

func Init(addrs []string, chanSize int) {
	config := sarama.NewConfig()
	// 发送文件需要leader和follower都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机选择分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("kafka connected")
	client = syncProducer
	clientChan = make(chan *lodData, chanSize)
	go sendToKafka()
}

func SendToChan(topic, msg string) {
	fmt.Printf("got item:%v,%v\n", topic, msg)
	clientChan <- &lodData{
		topic: topic,
		msg:   msg,
	}
}

func sendToKafka() {
	for {
		select {
		case msg := <-clientChan:
			// 构造一个message
			message := &sarama.ProducerMessage{}
			message.Topic = msg.topic
			message.Value = sarama.StringEncoder(msg.msg)
			partition, offset, err := client.SendMessage(message)
			if err != nil {
				log.Printf("syncProducer.SendMessage error:%v", err)
				return
			}

			fmt.Printf("pid:%v offset:%v\n", partition, offset)
		}
	}
}
