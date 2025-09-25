package config

import (
	"log/slog"
	"os"
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
