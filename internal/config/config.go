package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT  string
	ENV   string
	DBURL string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is requires")
	}

	dburl := os.Getenv("DBURL")
	if dburl == "" {
		panic("data_url is required")
	}
	return Config{
		PORT:  port,
		ENV:   env,
		DBURL: dburl,
	}
}
