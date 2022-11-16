package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Debug       bool   `env:"DEBUG" envDefault:"false"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"fi-fibonacci-fee-rate-activation"`

	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	LogOutput   string `env:"LOG_OUTPUT" envDefault:"stdout"`
	LogEncoding string `env:"LOG_ENCODING" envDefault:"json"`

	ServerHost         string        `env:"HTTP_SERVER_HOST" envDefault:"0.0.0.0"`
	ServerPort         string        `env:"HTTP_SERVER_PORT" envDefault:"80"`
	ServerReadTimeout  time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" envDefault:"1200s"`
	ServerWriteTimeout time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" envDefault:"1200s"`

	DBPgURI         string `env:"DB_POSTGRESQL_URI" envDefault:""`
	MaxIdleConn     int    `env:"DB_MAX_IDLE_CONN" envDefault:"3"`
	MaxOpenConn     int    `env:"DB_MAX_OPEN_CONN" envDefault:"6"`
	MaxLifetimeConn int64  `env:"DB_MAX_LIFETIME_CONN" envDefault:"1"`

	ModelMicroserviceURI             string `env:"MODEL_MICROSERVICE_URI" envDefault:"http://0.0.0.0:80"`
	ModelMicroservicePredictEndpoint string `env:"MODEL_MICROSERVICE_PREDICT_ENDPOINT" envDefault:"/combined_model/predict_multiple"`

	LocalStorageFolder string `env:"LOCAL_STORAGE_FOLDER" envDefault:"data"`

	PlotMicroserviceURI                string `env:"PLOT_MICROSERVICE_URI" envDefault:"http://0.0.0.0:80"`
	PlotMicroserviceEfficiencyEndpoint string `env:"PLOT_MICROSERVICE_EFFICIENCY_ENDPOINT" envDefault:"/plot_efficiency"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
