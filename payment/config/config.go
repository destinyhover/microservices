package config

import (
	"log/slog"
	"os"
	"strconv"
)

func GetVar(str string) string {
	vr := os.Getenv(str)
	if vr == "" {
		slog.Error("cannot get env var", "str", str, "var", vr)
	}
	return vr
}

func GetDataSourceURL() string {
	return GetVar("DATA_SOURCE_URL")
}

func GetEnv() string {
	return GetVar("ENV")
}
func GetPaymentPort() string {
	return GetVar("PAYMENT_PORT")
}

func GetOTLPEndpoint() string {
	if endpoint, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT"); ok {
		return endpoint
	}

	return ""
}

func IsOTLPInsecure() bool {
	v := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE")
	if v == "" {
		return true
	}

	insecure, err := strconv.ParseBool(v)
	if err != nil {
		slog.Error("invalid OTEL_EXPORTER_OTLP_INSECURE value", "value", v)
		return true
	}

	return insecure
}

// func GetPaymentSourceURL() string {
// 	return GetVar("PAYMENT_SERVICE_URL")
// }
// func GetApplicationPort() int {
// 	v := GetVar("APPLICATION_PORT")
// 	i, err := strconv.Atoi(v)
// 	if err != nil {
// 		slog.Error("cannot conv", "port", v, "err", err)
// 	}
// 	return i
// }
