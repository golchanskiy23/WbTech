package config

import (
	viper "github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	App           App        `mapstructure:"app"`
	Server        HttpServer `mapstructure:"server"`
	Database      DB         `mapstructure:"database"`
	NatsStreaming NATS       `mapstructure:"nats_streaming"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type HttpServer struct {
	ReadTimeout     *time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    *time.Duration `mapstructure:"write_timeout"`
	Addr            string         `mapstructure:"addr"`
	ShutdownTimeout *time.Duration `mapstructure:"shutdown_timeout"`
}

type DB struct {
	Name        string `mapstructure:"name"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Schema      string `mapstructure:"schema"`
	MaxPoolSize int    `mapstructure:"max_pool_size"`
	user        string `mapstructure:"user"`
	password    string `mapstructure:"password"`
}

type NATS struct {
	URL string `mapstructure:"url"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		// заменить обёртку ошибок
		log.Fatal("Не удается найти файл .env : ", err)
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		// заменить обёртку ошибок
		log.Fatal("Не удается загрузить среду: ", err)
		return nil, err
	}
	return cfg, nil
}
