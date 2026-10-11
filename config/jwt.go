package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type JWTConfig struct {
	SecretKey       string
	ExpirationHours time.Duration
}

func LoadJWTConfig() (*JWTConfig, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not configured")
	}
	hours := 24
	if v := os.Getenv("JWT_EXPIRATION_HOURS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return nil, errors.New("JWT_EXPIRATION_HOURS must be a positive integer")
		}
		hours = n
	}
	return &JWTConfig{SecretKey: secret, ExpirationHours: time.Duration(hours) * time.Hour}, nil
}
