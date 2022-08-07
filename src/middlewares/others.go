package middlewares

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvLocal() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Env local")
	}
}
