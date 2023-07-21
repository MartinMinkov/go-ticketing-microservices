package config

import (
	"os"
	"strconv"
)

type Config struct {
	ApplicationConfig *ApplicationConfig
	DatabaseConfig    *DatabaseConfig
}

type ApplicationConfig struct {
	Environment string
	Host        string
	Port        string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
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

func BuildDatabaseConfig() *DatabaseConfig {
	port, err := strconv.Atoi(os.Getenv("MONGO_PORT"))
	if err != nil {
		port = 27017
	}

	host := os.Getenv("MONGO_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("MONGO_USER")
	if user == "" {
		user = "admin"
	}

	password := os.Getenv("MONGO_PASSWORD")
	if password == "" {
		password = "pass"
	}

	database := os.Getenv("MONGO_DB_NAME")
	if database == "" {
		database = "tickets"
	}

	return &DatabaseConfig{
		Host:     host,
		Port:     port,
		Username: user,
		Password: password,
		Database: database,
	}
}

func BuildConfig() *Config {
	var config Config
	applicationConfig := BuildApplicationConfig()
	config.ApplicationConfig = applicationConfig

	databaseConfig := BuildDatabaseConfig()
	config.DatabaseConfig = databaseConfig

	return &config
}

func (dbConfig *DatabaseConfig) GetConnectionString() string {
	return "mongodb://" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + "/" + dbConfig.Database
}

func (dbConfig *DatabaseConfig) GetConnectionStringWithUser() string {
	return "mongodb://" + dbConfig.Username + ":" + dbConfig.Password + "@" +
		dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + "/" + dbConfig.Database
}

func (config *ApplicationConfig) GetAddress() string {
	return "http://" + config.Host + ":" + config.Port
}

func (config *ApplicationConfig) IsProduction() bool {
	return config.Environment == "production"
}

func GetJWTSecret() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", nil
	}
	return secret, nil
}
