package main

import (
    "log"
)

type Config struct {
    Url      string
    Username string
    Password string
    Account  string
    Role     string
}

func (c *Config) Client() (*AlksClient, error) {
    client, err := NewAlksClient(c.Url, c.Username, c.Password, c.Account, c.Role)

    if err != nil {
        return nil, err
    }

    log.Printf("[INFO] ALKS Client configured")

    return client, nil
}