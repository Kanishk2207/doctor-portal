package configs

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DB_DSN      string
	GRPCAddress string
	HTTPAddress string
	JWTSecret   string
	JWTExpiry   int
}

func LoadConfig() *Config {
	// JWTExpiry, err := strconv.Atoi(os.Getenv("JWTEXPIRY"))
	JWTExpiry, err := strconv.Atoi("86400")
	if err != nil {
		log.Fatalf("Error occurred in config: %v", err)
	}

	var config Config = Config{
		DB_DSN:      os.Getenv("DSN"),
		GRPCAddress: ":" + os.Getenv("GRPCADDRESS"),
		HTTPAddress: ":" + os.Getenv("HTTPADDRESS"),
		JWTSecret:   os.Getenv("JWTSECRET"),
		JWTExpiry:   JWTExpiry,
	}

	return &config
}
