package configuration

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/estrategiahq/backend-libs/logger"
	"github.com/spf13/viper"
)

// Configuration wrapper
type Configuration struct {
	v   *viper.Viper
	env string
}

func LoadConfig() *Configuration {
	env := setEnvironment()
	log.Println("Using env:", env)
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

	for _, key := range v.AllKeys() {
		val := v.GetString(key)
		if val[0] == '$' {
			val = os.Getenv(val[1:])
		}
		v.Set(key, val)
	}

	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		logger.Fatal(context.TODO(), "não foi possível fazer o unmarshal da configuracao")
	}

	return &Configuration{v, env}
}

func setEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "development"
	}
	return env
}

// GetEnv returns the current env
func (config *Configuration) GetEnv() string {
	return string(config.env)
}

// GetEnvConfString returns a value and key on configuration
func (config *Configuration) GetEnvConfString(key string) string {
	return config.v.GetString(fmt.Sprintf("%s.%s", config.env, key))
}

// GetEnvConfBool returns a value of and key on configuration as boolean
func (config *Configuration) GetEnvConfBool(key string) bool {
	return config.v.GetBool(fmt.Sprintf("%s.%s", config.env, key))
}

// GetString return a value of a key as String
func (config *Configuration) GetString(key string) string {
	return config.v.GetString(key)
}

// GetEnvConfStringMap returns a value and key on configuration as map[string]
func (config *Configuration) GetEnvConfStringMap(key string) map[string]interface{} {
	return config.v.GetStringMap(fmt.Sprintf("%s.%s", config.env, key))
}

// GetEnvConfInteger returns a value of and key on configuration as boolean
func (config *Configuration) GetEnvConfInteger(key string) int {
	return config.v.GetInt(fmt.Sprintf("%s.%s", config.env, key))
}

// GetInteger return a value of a key as Integer
func (config *Configuration) GetInteger(key string) int {
	return config.v.GetInt(key)
}

// GetEnvConfDuration returns a value of a key as time.Duration
func (config *Configuration) GetEnvConfDuration(key string) time.Duration {
	return config.v.GetDuration(fmt.Sprintf("%s.%s", config.env, key))
}

type DatabaseConfig struct {
	UserName     string
	Password     string
	Host         string
	DatabaseName string
	DatabasePort int
}

func (config *Configuration) BuildDatabaseConfig() DatabaseConfig {
	host := config.GetEnvConfString("database.host")
	databasePort := config.GetEnvConfInteger("database.database_port")
	userName := config.GetEnvConfString("database.user")
	password := config.GetEnvConfString("database.password")
	databaseName := config.GetEnvConfString("database.database_name")

	return DatabaseConfig{
		Host:         host,
		DatabaseName: databaseName,
		DatabasePort: databasePort,
		UserName:     userName,
		Password:     password,
	}
}

type ServerConfig struct {
	Port                     int
	GracefullShutdownTimeout int
}

func (config *Configuration) BuildServerConfig() ServerConfig {
	port := config.GetEnvConfInteger("server.port")
	gracefullShutdownTimeout := config.GetEnvConfInteger("server.gracefull_shutdown_timeout")

	return ServerConfig{
		Port:                     port,
		GracefullShutdownTimeout: gracefullShutdownTimeout,
	}
}
