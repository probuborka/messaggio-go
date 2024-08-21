package config

import (
	"time"

	"github.com/probuborka/messaggio/internal/domain"
	"github.com/spf13/viper"
)

const (
	//HTTP
	defaultHTTPPort      = "8000"
	defaultHTTPRWTimeout = 10 * time.Second
)

type Config struct {
	path string
	file string
	//
	HTTP  domain.HTTPConfig
	DB    domain.DBConfig
	Kafka domain.KafkaConfig
}

func Init(cfgDir string, cfgFile string) (*Config, error) {
	//
	cfgPath, err := domain.Path(cfgDir)
	if err != nil {
		return nil, err
	}
	//
	cfg := &Config{
		path: cfgPath,
		file: cfgFile,
	}
	if err := cfg.defaults(); err != nil {
		return nil, err
	}
	if err := cfg.loading(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) defaults() error {
	c.HTTP.Port = defaultHTTPPort
	c.HTTP.ReadTimeout = defaultHTTPRWTimeout
	c.HTTP.WriteTimeout = defaultHTTPRWTimeout

	return nil
}

func (c *Config) loading() error {
	viper.AddConfigPath(c.path)
	viper.SetConfigName(c.file)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &c.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("db", &c.DB); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("kafka", &c.Kafka); err != nil {
		return err
	}

	return nil
}
