package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	Environment      string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string
	LogLevel         string
	RPCPort          string
	LockPort         string
	LockHost         string
	AdminServiceHost string
	AdminServicePort string
}

func Load() Config {
	c := Config{}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DB", "lockdb"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "compos1995"))

	c.AdminServiceHost = cast.ToString(getOrReturnDefault("ADMIN_HOST", "localhost"))
	c.AdminServicePort = cast.ToString(getOrReturnDefault("ADMIN_PORT", "8088"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.LockPort = cast.ToString(getOrReturnDefault("LOCK_PORT", "9679"))
	c.LockHost = cast.ToString(getOrReturnDefault("LOCK_HOST", "143.42.61.34"))
	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":8820"))
	return c
}

func getOrReturnDefault(key string, defaulValue interface{}) interface{} {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return defaulValue
	}
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaulValue
}
