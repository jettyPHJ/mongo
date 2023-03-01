package kafka

import (
	"ToMongoAndRedis/utils/config"
	"ToMongoAndRedis/utils/log"

	"github.com/Shopify/sarama"
)

var Consumer sarama.Consumer

func init() {
	var err error
	Consumer, err = sarama.NewConsumer([]string{config.Param.Kafka}, nil)
	if err != nil {
		log.ErrorLogger.Println("consumer connect error:", err)
		return
	}
}
