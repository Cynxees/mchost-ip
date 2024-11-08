package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Listen address is an array of IP addresses and port combinations.
	// Listen address is an array so that this service can listen to many interfaces at once.
	// You can use this value for example: []string{"192.168.1.12:80", "25.49.25.73:80"} to listen to
	// listen to interfaces with IP address of 192.168.1.12 and 25.49.25.73, both on port 80.
	ListenAddress string `config:"LISTEN_ADDRESS"`

	CorsAllowedHeaders []string `config:"CORS_ALLOWED_HEADERS"`
	CorsAllowedMethods []string `config:"CORS_ALLOWED_METHODS"`
	CorsAllowedOrigins []string `config:"CORS_ALLOWED_ORIGINS"`

	AppName string `config:"APP_NAME"`
	AppKey  string `config:"APP_KEY"`

	AppPort          string `config:"APP_PORT"`
	MicroservicePort string `config:"MICROSERVICE_PORT"`

	DbPort     string `config:"DB_PORT"`
	DbHost     string `config:"DB_HOST"`
	DbName     string `config:"DB_NAME"`
	DbUsername string `config:"DB_USERNAME"`
	DbPassword string `config:"DB_PASSWORD"`

	AwsAccessKeyId     string `config:"AWS_ACCESS_KEY_ID"`
	AwsAccessKeySecret string `config:"AWS_ACCESS_KEY_SECRET"`
}

func InitConfig(path string) *Config {
	// Todo: add env checker

	if path == "" {
		godotenv.Load(".env")
	} else {
		godotenv.Load(path)
	}

	appName := getEnv("APP_NAME")
	appKey := getEnv("APP_KEY")
	appPort := getEnv("APP_PORT")
	microservicePort := getEnv("MICROSERVICE_PORT")
	dbPort := getEnv("DB_PORT")
	dbHost := getEnv("DB_HOST")
	dbName := getEnv("DB_NAME")
	dbUsername := getEnv("DB_USERNAME")
	dbPassword := getEnv("DB_PASSWORD")
	awsAccessKeyId := getEnv("AWS_ACCESS_KEY_ID")
	awsAccessKeySecret := getEnv("AWS_ACCESS_KEY_SECRET")

	return &Config{
		ListenAddress: fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		CorsAllowedHeaders: []string{
			"Connection", "User-Agent", "Referer",
			"Accept", "Accept-Language", "Content-Type",
			"Content-Language", "Content-Disposition", "Origin",
			"Content-Length", "Authorization", "ResponseType",
			"X-Requested-With", "X-Forwarded-For",
		},
		CorsAllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT"},
		CorsAllowedOrigins: []string{"*"},
		AppName:            appName,
		AppKey:             appKey,
		AppPort:            appPort,
		MicroservicePort:   microservicePort,
		DbPort:             dbPort,
		DbHost:             dbHost,
		DbName:             dbName,
		DbUsername:         dbUsername,
		DbPassword:         dbPassword,
		AwsAccessKeyId:     awsAccessKeyId,
		AwsAccessKeySecret: awsAccessKeySecret,
	}
}

func (c *Config) AsString() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("Environment variable %s is not set", key))

}
