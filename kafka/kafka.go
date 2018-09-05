package kafka

import (
	"github.com/kumparan/go-lib/logger"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

var kafkaProducer sarama.SyncProducer
var kafkaConsumer sarama.Consumer

type Message struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

func Init(producer *sarama.SyncProducer) {
	kafkaProducer = *producer
	logger.Info("kafka producer init")
}

func InitCons(consumer *sarama.Consumer) {
	kafkaConsumer = *consumer
	logger.Info("kafka consumer init")
}

func Publish(msg Message) error {
	logger.Infof("Message receive: %v", msg)
	_, _, err := kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.StringEncoder(msg.Content),
	})
	if err != nil {
		logger.Infof("Error receive: %v", err)
	}
	return err
}

func Consume(topic string, callback func([]byte)) (offset int64, err error) {

	logger.Infof("Start consuming topic: %v", topic)
	partitionConsumer, err := kafkaConsumer.ConsumePartition(topic, 0, -1)
	if err != nil {
		logger.Infof("Error receive: %v", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			logger.Infof("Error receive: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			callback(msg.Value)
			offset = msg.Offset
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	return offset, err
}
