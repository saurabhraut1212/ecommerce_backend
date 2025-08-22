package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  mustEnv("MONGO_URI"),
		MongoDB:   getEnv("MONGO_DB", "ecommerce"),
		JWTSecret: mustEnv("JWT_SECRET"),
	}
}

func getEnv(k, d string) string {
	v := os.Getenv(k)
	if v != "" {
		return v
	}
	return d

}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v

}
