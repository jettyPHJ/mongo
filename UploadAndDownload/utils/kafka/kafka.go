package kafka

import (
	cfg "UploadAndDownload/utils/config"
	"fmt"

	"github.com/Shopify/sarama"
)

var SyncProducer sarama.SyncProducer

func init() {
	// 新建一个arama配置实例
	config := sarama.NewConfig()

	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll

	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Return.Successes = true

	// 新建一个同步生产者
	var err error
	SyncProducer, err = sarama.NewSyncProducer([]string{cfg.Param.Kafka}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
}
