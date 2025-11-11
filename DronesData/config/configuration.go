package config

import (
	"log"

	"github.com/joho/godotenv"
)

// ************************************************************
// use to load the properties from properties.env file

func LoadLocal() {

	envPath := "./.env"

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading %s: %v", envPath, err)
	}
}
