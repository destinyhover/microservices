package config

import (
	"log/slog"
	"os"
	"strconv"
)

func getEnvValue(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("environment variable is missing.", "key", key)
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
		slog.Error("port is invalid", "port", str)

	}
	return port
}

// Tracing variabales
func GetOLTPEndpoint() string {
	if endp, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT"); ok {
		return endp
	}
	return ""
}
func IsOTLPInsecure() bool {
	v := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE")
	if v == "" {
		return true
	}

	ins, err := strconv.ParseBool(v)
	if err != nil {
		slog.Error("Cannot parse isotlpinsecure()", "err", err)
		return false
	}
	return ins
}
