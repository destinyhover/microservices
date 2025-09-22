package config

import (
	"log"
	"os"
	"strconv"
)

func getEnvValue(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}

	return v
}
func GetEnv() string {
	return getEnvValue("ENV")
}

func GetDataSourceURL() string {
	return getEnvValue("DATA_SOURCE_URL")
}
func GetPaymentSourceURL() string {
	return getEnvValue("PAYMENT_SERVICE_URL")
}

func GetApplicationPort() int {
	str := getEnvValue("APPLICATION_PORT")
	port, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("port: %s is invalid", str)
	}
	return port
}
