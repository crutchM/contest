package main

import (
	"Contest/internal/app"
	. "Contest/internal/domain"
	"context"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	cfg := MustLoadConfig()
	application := app.New(cfg)
	application.MustRun()
}

func MustLoadConfig() *Config {
	cfg := &Config{}

	err := envconfig.Process(context.Background(), cfg)
	if err != nil {
		panic("Load Config error: " + err.Error())
	}

	return cfg
}
