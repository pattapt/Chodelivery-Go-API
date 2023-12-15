package config

import (
	"time"

	"github.com/google/uuid"
)

func GenerateRefreshToken() string {
	return uuid.New().String()
}

func GetCurrentDateTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

func GetMainAPIURL() string {
	return "https://mobile-api-gateway.patta.dev"
}
