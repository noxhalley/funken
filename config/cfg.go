package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

const cfgFilePath = ".env"

type (
	Config struct {
		App   app
		Mongo mongo
	}

	app struct {
		Port     string `env:"APP_PORT"  env-required:"true"`
		LogLevel string `env:"LOG_LEVEL" env-default:"DEBUG"`
	}

	mongo struct {
		Hostname     string `env:"MONGO_HOST"          env-default:"localhost"`
		Port         string `env:"MONGO_PORT"          env-default:"27017"`
		Database     string `env:"MONGO_DATABASE"      env-required:"true"`
		AuthSource   string `env:"MONGO_AUTH_DB"       env-required:"true"`
		Timeout      int    `env:"MONGO_TIMEOUT"       env-default:"30000"`
		ConnTimeout  int    `env:"MONGO_CONN_TIMEOUT"  env-default:"30000"`
		PoolSize     int    `env:"MONGO_POOL_SIZE"     env-default:"10"`
		MaxIdleTime  int    `env:"MONGO_MAX_IDLE_TIME" env-default:"300000"`
		ConnAttempts int    `env:"MONGO_CONN_ATTEMPTS" env-default:"3"`
		Username     string `env:"MONGO_USERNAME"      env-required:"true"`
		Password     string `env:"MONGO_PASSWORD"      env-required:"true"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	root := projectRoot()
	configFilePath := root + cfgFilePath

	err := loadCfg(configFilePath, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadCfg(cfgFilePath string, cfg *Config) error {
	envFileExists := checkFileExists(cfgFilePath)
	if envFileExists {
		err := cleanenv.ReadConfig(cfgFilePath, cfg)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	} else {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			if _, statErr := os.Stat(cfgFilePath); statErr != nil {
				return fmt.Errorf("missing environment variable: %w", err)
			}
			return err
		}
	}
	return nil
}

func checkFileExists(fileName string) bool {
	exist := false
	if _, err := os.Stat(fileName); err == nil {
		exist = true
	}
	return exist
}

func projectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(b)
	return cwd + "/../"
}
