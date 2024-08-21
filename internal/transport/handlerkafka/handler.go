package handlerkafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/probuborka/messaggio/internal/domain"
	"github.com/probuborka/messaggio/internal/service"
	"github.com/probuborka/messaggio/pkg/kafka/kafkago"
	"github.com/probuborka/messaggio/pkg/logger"
)

type Handler struct {
	message service.Message
}

func New(services *service.Services) *Handler {
	return &Handler{
		message: services.Message,
	}
}

func (h Handler) Init(kafkaConfig domain.KafkaConfig) []kafkago.Config {
	consumers := make([]kafkago.Config, 0)

	// MessageProcessing
	consumers = append(consumers, kafkago.Config{
		ConsumerConfig: kafkago.ConsumerConfig{
			KafkaURL: fmt.Sprintf("%s:%s", kafkaConfig.Host, kafkaConfig.Port),
			Topic:    "message",
			GroupID:  "test",
		},
		HandlerFn: h.MessageProcessing, // run func
		Processes: kafkago.Processes{
			Ch:  make(chan struct{}, 3), // count processes
			Any: false,                  // any count of processes
		},
	})

	return consumers
}

// MessageProcessing
func (h Handler) MessageProcessing(msg []byte, processes kafkago.Processes) {
	//json
	message := &domain.Message{}
	err := json.Unmarshal(msg, message)
	if err != nil {
		logger.Error(err)
	}

	//
	err = h.message.Process(context.Background(), *message)
	if err != nil {
		logger.Error(err)
		return
	}

	//
	if !processes.Any {
		<-processes.Ch
	}
}
