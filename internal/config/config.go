package config

import (
	"errors"
	"os"
)

type Config struct {
	Address string
}

func New() Config {
	return Config{
		Address: os.Getenv("ADDRESS"),
	}
}

func (c Config) Validate() error {
	if c.Address == "" {
		return errors.New("ADDRESS should not be empty")
	}

	return nil
}
