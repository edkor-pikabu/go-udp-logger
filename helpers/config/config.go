package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
    AppPort string
	MysqlDsn string
}

func New() *Config {
    return &Config{
        AppPort: ":" + getEnv("APP_PORT", ""),
        MysqlDsn: getEnv("MYSQL_DSN", ""),
	}
}

func getEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
	return value
    }

    return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
    valueStr := getEnv(name, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
	return value
    }

    return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
    valStr := getEnv(name, "")
    if val, err := strconv.ParseBool(valStr); err == nil {
	return val
    }

    return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
    valStr := getEnv(name, "")

    if valStr == "" {
	return defaultVal
    }

    val := strings.Split(valStr, sep)

    return val
}