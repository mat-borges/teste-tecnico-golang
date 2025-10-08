package config

import (
	"log"
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
	log.Printf("[config] loaded (port=%s, usersURL=%s, postsURL=%s, timeouts=%s/%s)",
		cfg.ServerPort, cfg.UsersBaseURL, cfg.PostsBaseURL, cfg.HTTPTimeout, cfg.AggTimeout)
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
		log.Printf("[config] invalid duration for %s: %s, using default %v", key, val, defaultVal)
	}
	return defaultVal
}