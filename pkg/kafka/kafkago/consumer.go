package kafkago

import (
	"context"
	"strings"

	"github.com/probuborka/messaggio/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader    *kafka.Reader
	handlerFn HandlerFunc
	processes Processes
}

type Config struct {
	ConsumerConfig ConsumerConfig
	HandlerFn      HandlerFunc
	Processes      Processes
}

type ConsumerConfig struct {
	KafkaURL string
	Topic    string
	GroupID  string
}

type HandlerFunc func([]byte, Processes)

type Processes struct {
	Ch  chan struct{} // количество очереди в канале = количеству процессов одновременно запущенных для обработки сообщений в consumer
	Any bool          // если true любое количество горутин на обработку сообщений в consumer
}

func NewConsumer(cfg ConsumerConfig, processes Processes, handlerFn HandlerFunc) Consumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  strings.Split(cfg.KafkaURL, ","),
			GroupID:  cfg.GroupID,
			Topic:    cfg.Topic,
			MaxBytes: 10e6, // 10MB
		}),
		handlerFn: handlerFn,
		processes: processes,
	}
}

type Consumer interface {
	Read(ctx context.Context)
}

func (k *KafkaConsumer) Read(ctx context.Context) {
	defer k.reader.Close()

	for {
		if !k.processes.Any {
			k.processes.Ch <- struct{}{}
		}

		m, err := k.reader.ReadMessage(ctx)
		if err != nil {
			logger.Error(err)
			continue
		}

		go k.handlerFn(m.Value, k.processes)
	}
}
