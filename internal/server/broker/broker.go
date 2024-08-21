package broker

import (
	"context"

	"github.com/probuborka/messaggio/pkg/kafka/kafkago"
)

func Run(consumers []kafkago.Config) error {

	for _, v := range consumers {
		go func(consumerMessage kafkago.Consumer) {
			consumerMessage.Read(context.Background())
		}(kafkago.NewConsumer(v.ConsumerConfig, v.Processes, v.HandlerFn))
	}

	return nil
}
