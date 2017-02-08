package main

import (
    // "encoding/json"
    "net/http"
    // "net/url"
    "github.com/hashicorp/go-cleanhttp"
)

type AlksClient struct {
    BaseURL  string
    Username string
    Password string
    Account  string
    Role     string

    Http *http.Client
}

func NewAlksClient(url string, username string, password string, account string, role string) (*AlksClient, error) {
    alksClient := AlksClient{
        BaseURL:  url,
        Username: username,
        Password: password,
        Account:  account,
        Role:     role,
        Http:     cleanhttp.DefaultClient(),
    }

    return &alksClient, nil
}
