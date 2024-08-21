package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/probuborka/messaggio/internal/config"
	"github.com/probuborka/messaggio/internal/infrastructure/producer"
	"github.com/probuborka/messaggio/internal/infrastructure/repository"
	"github.com/probuborka/messaggio/internal/server/broker"
	"github.com/probuborka/messaggio/internal/server/serverhttp"
	"github.com/probuborka/messaggio/internal/service"
	"github.com/probuborka/messaggio/internal/transport/handlerhttp"
	"github.com/probuborka/messaggio/internal/transport/handlerkafka"
	"github.com/probuborka/messaggio/pkg/database/postgres"
	"github.com/probuborka/messaggio/pkg/kafka/kafkago"
	"github.com/probuborka/messaggio/pkg/logger"
)

const (
	cfgDir  = "configs"
	cfgFile = "config_docker"
)

func Run() {
	//config
	cfg, err := config.Init(cfgDir, cfgFile)
	if err != nil {
		logger.Error(err)
		return
	}

	//db
	logger.Info("db")

	db, err := postgres.New(context.Background(), cfg.DB)
	if err != nil {
		logger.Error(err)
		return
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		logger.Error(err)
		return
	}

	repo := repository.New(db)

	//kafka
	//producer
	logger.Info("producer")

	//clientProducer, err := kafkasarama.NewProducer(cfg.Kafka)
	clientProducer, err := kafkago.NewProducer(cfg.Kafka)
	if err != nil {
		logger.Error(err)
		return
	}
	defer clientProducer.Close()

	producer := producer.New(*clientProducer)

	//services
	services := service.New(
		repo,
		producer,
	)

	//
	logger.Info("server started")

	//broker
	//handlers
	logger.Info("broker")

	kafkaHandler := handlerkafka.New(services)

	if err := broker.Run(kafkaHandler.Init(cfg.Kafka)); err != nil {
		logger.Errorf("broker: %v", err)
	}

	//HTTP
	//handlers
	logger.Info("http")

	httpHandlers := handlerhttp.New(services)

	//server
	httpServer := serverhttp.New(cfg.HTTP, httpHandlers.Init())

	go func() {
		if err := httpServer.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	//stop server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("server stop")

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := httpServer.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
