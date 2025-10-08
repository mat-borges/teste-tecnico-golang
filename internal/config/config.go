package config

import (
	"go-graphql-aggregator/internal/logger"
	"os"
	"time"
)

type Config struct {
	ServerPort   string
	UsersBaseURL string
	PostsBaseURL string
	HTTPTimeout  time.Duration
	AggTimeout   time.Duration
}

func LoadConfig() *Config {
	cfg := &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		UsersBaseURL: getEnv("USERS_BASE_URL", "https://jsonplaceholder.typicode.com/users"),
		PostsBaseURL: getEnv("POSTS_BASE_URL", "https://jsonplaceholder.typicode.com/posts"),
		HTTPTimeout:  getEnvAsDuration("HTTP_TIMEOUT", 5*time.Second),
		AggTimeout:   getEnvAsDuration("AGG_TIMEOUT", 5*time.Second),
	}
	logger.Log.Info("config loaded",
		"port", cfg.ServerPort,
		"usersURL", cfg.UsersBaseURL,
		"postsURL", cfg.PostsBaseURL,
		"httpTimeout", cfg.HTTPTimeout,
		"aggTimeout", cfg.AggTimeout,
	)
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		dur, err := time.ParseDuration(val)
		if err == nil {
			return dur
		}
		logger.Log.Info("invalid duration for env var, using default", "key", key, "val", val, "default", defaultVal)
	}
	return defaultVal
}