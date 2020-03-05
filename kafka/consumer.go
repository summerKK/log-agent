package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func InitConsumer(addrs []string, topic string) {
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		panic(err)
	}
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		panic(err)
	}
	fmt.Printf("topci:%v,获取分区:%v\n", topic, partitions)

	for _, partition := range partitions {
		consumePartition, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(fmt.Sprintf("failed to start consumer for partition %d,err:%v\n", partition, err))
		}
		defer consumePartition.AsyncClose()
		// 异步从每个分区消费信息
		go func(cp sarama.PartitionConsumer) {
			for msg := range cp.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(consumePartition)
	}

	for {
		time.Sleep(time.Millisecond)
	}
}
