package main

import (
	"github.com/Cox-Automotive/alks-go"
	"log"
)

type Config struct {
	Url      string
	Username string
	Password string
	Account  string
	Role     string
}

func (c *Config) Client() (*alks.Client, error) {
	client, err := alks.NewClient(c.Url, c.Username, c.Password, c.Account, c.Role)

	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] ALKS Client configured")

	return client, nil
}
