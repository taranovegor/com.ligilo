package config

import "os"

const (
	AmqpDsn     = "AMQP_DSN"
	AppUrl      = "APP_URL"
	DatabaseDsn = "DATABASE_DSN"
	FallbackUrl = "FALLBACK_URL"
	HttpPort    = "HTTP_PORT"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}
