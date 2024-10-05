package internal

import (
	"fmt"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ServerConfig struct {
	Port  int  `env:"SERVER_PORT,required"`
	Debug bool `env:"SERVER_DEBUG,required"`
	DB    struct {
		Host  string `env:"DB_HOST,required"`
		Name  string `env:"DB_NAME,required"`
		Port  int    `env:"DB_PORT,required"`
		User  string `env:"DB_USER,required"`
		Pass  string `env:"DB_PASS,required"`
		Debug bool   `env:"DB_DEBUG,required"`
	}
}

func NewServerConfig() *ServerConfig {
	// Find .env file
	goEnv := os.Getenv("GO_ENV")
	envFileName := fmt.Sprintf("%s.env", goEnv)
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	var cfg ServerConfig
	err = envdecode.StrictDecode(&cfg)
	//err = envdecode.Decode(&cfg)
	if err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &cfg
}
