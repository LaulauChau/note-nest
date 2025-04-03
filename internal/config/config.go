package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE struct {
		URL string
	}

	Server struct {
		Port int
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{}

	config.DATABASE.URL, err = parseURL("DATABASE_URL")
	if err != nil {
		return nil, fmt.Errorf("error parsing DATABASE_URL: %w", err)
	}

	config.Server.Port, err = parseInt("SERVER_PORT")
	if err != nil {
		return nil, fmt.Errorf("error parsing SERVER_PORT: %w", err)
	}

	return config, nil
}

func parseURL(envName string) (string, error) {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return "", fmt.Errorf("%s is not set", envName)
	}

	_, err := url.Parse(os.Getenv(envName))
	if err != nil {
		return "", fmt.Errorf("invalid %s: %w", envName, err)
	}

	return envValue, nil
}

func parseInt(envName string) (int, error) {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return 0, fmt.Errorf("%s is not set", envName)
	}

	port, err := strconv.Atoi(envValue)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", envName, err)
	}

	return port, nil
}
