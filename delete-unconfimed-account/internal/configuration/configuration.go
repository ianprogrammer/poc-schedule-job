package configuration

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/estrategiahq/backend-libs/logger"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	UserName               string `yaml:"user"`
	Password               string `yaml:"password"`
	Host                   string `yaml:"host"`
	DatabaseName           string `yaml:"database_name"`
	DatabasePort           int    `yaml:"database_port"`
	MaxIdleConns           int    `yaml:"max_idle_connections"`
	MaxOpenConns           int    `yaml:"max_open_connections"`
	ConnMaxLifetimeSeconds int    `yaml:"max_lifetime_connections_seconds"`
}
type Configuration struct {
	Database    DatabaseConfig
	Environment string
}

func LoadConfig() Configuration {
	env := getEnvironment()
	fmt.Println("Using env:", env)
	v := viper.New()
	v.SetConfigName("config")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(env)
	viper.SetConfigType("yaml")
	v.AddConfigPath(os.Getenv("CONFIG_DIR"))
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		logger.Fatal(context.TODO(), "não consegui ler o arquivo de configuracao:", err)
	}

	// Fill config with ENV variables
	for _, key := range v.AllKeys() {
		val := v.GetString(key)
		if val[0] == '$' {
			val = os.Getenv(val[1:])
		}
		v.Set(key, val)
	}
	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		logger.Fatal(context.TODO(), "não foi possível fazer o unmarshal da configuracao")
	}
	config.Environment = env

	return config
}

func getEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "development"
	}
	return env
}
