package kfk_moudle

import (
	kfk "ToMongoAndRedis/utils/kafka"
	"ToMongoAndRedis/utils/log"

	"github.com/Shopify/sarama"
)

type Partitions struct {
	topic string
	parts []int32
}

type ConsumeAct func(partConsumer sarama.PartitionConsumer)

func NewPartitions(topic string) *Partitions {
	Partitions := Partitions{}
	Partitions.topic = topic
	var err error
	Partitions.parts, err = kfk.Consumer.Partitions(topic)
	if err != nil {
		log.ErrorLogger.Println("geet partitions failed, err:", err)
	}
	return &Partitions
}

//采用统一方式act消费所有partition,num为每个partition消费的协程数
func (p *Partitions) ConsumeAllParts(act ConsumeAct, num int) {
	for _, part := range p.parts {
		partConsumer, err := kfk.Consumer.ConsumePartition(p.topic, part, sarama.OffsetNewest)
		if err != nil {
			log.ErrorLogger.Println("partitionConsumer err:", err)
			continue
		}
		for i := 0; i < num; i++ {
			go func() {
				act(partConsumer)
			}()
		}
	}
}
