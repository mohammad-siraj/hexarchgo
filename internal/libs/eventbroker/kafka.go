package eventbroker

import (
	"fmt"

	"github.com/IBM/sarama"
)

type IEventBroker interface {
	SyncSendMessageToTopic(msg string, topic string, isAsync bool) (int32, int64, error)
	ConsumeTopic(topic string, offset int64) (<-chan string, error)
}

type eventbroker struct {
	consumer      sarama.Consumer
	asyncProducer sarama.AsyncProducer
	syncProducer  sarama.SyncProducer
}

func NewEventBroker(BrokerConfig []string) (IEventBroker, error) {
	brokers := BrokerConfig
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &eventbroker{
		syncProducer:  syncProducer,
		asyncProducer: asyncProducer,
		consumer:      consumer,
	}, nil
}

func (eb *eventbroker) SyncSendMessageToTopic(msg string, topic string, isAsync bool) (int32, int64, error) {
	if isAsync {
		return 0, 0, nil
	}
	partition, offset, err := eb.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		return 0, 0, err
	}
	return partition, offset, nil
}

func (eb *eventbroker) ConsumeTopic(topic string, offset int64) (<-chan string, error) {
	partitions, err := eb.consumer.Partitions(topic)

	if err != nil {
		return nil, err
	}
	ch := make(chan string)
	for _, partition := range partitions {
		fmt.Println(topic+" partitions are :", partition)
		go func(partition int32) {
			partitionConsumer, _ := eb.consumer.ConsumePartition(topic, partition, offset)
			for message := range partitionConsumer.Messages() {
				ch <- string(message.Value)
			}
		}(partition)
	}
	return ch, nil
}
