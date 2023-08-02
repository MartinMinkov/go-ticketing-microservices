package config

import (
	"os"
	"strconv"
)

type Config struct {
	ApplicationConfig *ApplicationConfig
	RedisConfig       *RedisConfig
	NatsConfig        *NatsConfig
}

type ApplicationConfig struct {
	Environment string
	Host        string
	Port        string
}

func (config *ApplicationConfig) GetAddress() string {
	return "http://" + config.Host + ":" + config.Port
}

type RedisConfig struct {
	Host string
	Port int
}

func (redisConfig *RedisConfig) GetAddress() string {
	return redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port)
}

type NatsConfig struct {
	Host string
	Port int
}

func (natsConfig *NatsConfig) GetAddress() string {
	return natsConfig.Host + ":" + strconv.Itoa(natsConfig.Port)
}

func BuildConfig() *Config {
	var config Config
	applicationConfig := BuildApplicationConfig()
	config.ApplicationConfig = applicationConfig

	natsConfig := BuildNatsConfig()
	config.NatsConfig = natsConfig

	redisConfig := BuildRedisConfig()
	config.RedisConfig = redisConfig

	return &config
}

func BuildApplicationConfig() *ApplicationConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	return &ApplicationConfig{
		Environment: env,
		Host:        host,
		Port:        port,
	}
}

func BuildRedisConfig() *RedisConfig {
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		port = 6379
	}

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	return &RedisConfig{
		Host: host,
		Port: port,
	}
}

func BuildNatsConfig() *NatsConfig {
	port, err := strconv.Atoi(os.Getenv("NATS_PORT"))
	if err != nil {
		port = 4222
	}

	host := os.Getenv("NATS_HOST")
	if host == "" {
		host = "localhost"
	}

	return &NatsConfig{
		Host: host,
		Port: port,
	}
}

func (config *ApplicationConfig) IsProduction() bool {
	return config.Environment == "production"
}
