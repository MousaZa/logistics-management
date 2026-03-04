package kafka_adapter

import (
	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
)

func NewSubscriber(brokers []string, consumerGroup string, logger watermill.LoggerAdapter) (*kafka.Subscriber, error) {
	saramaConfig := kafka.DefaultSaramaSubscriberConfig()

	// Read from the beginning if we haven't seen this consumer group before
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaConfig,
			ConsumerGroup:         consumerGroup, // e.g., "inventory_service"
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}
