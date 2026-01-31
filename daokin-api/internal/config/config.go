package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env             string
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		Env:             getenv("DAOKIN_ENV", "dev"),
		Addr:            getenv("DAOKIN_ADDR", ":8080"),
		ReadTimeout:     getDuration("DAOKIN_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDuration("DAOKIN_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:     getDuration("DAOKIN_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getDuration("DAOKIN_SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
