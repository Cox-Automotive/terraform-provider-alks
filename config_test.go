package main

import (
	"os"
	"testing"
)

func TestGetPluginVersion_WithVersion(t *testing.T) {
	original := versionNumber
	versionNumber = "1.2.3"
	defer func() { versionNumber = original }()

	if got := getPluginVersion(); got != "1.2.3" {
		t.Fatalf("expected 1.2.3, got %s", got)
	}
}

func TestGetPluginVersion_WithoutVersion(t *testing.T) {
	original := versionNumber
	versionNumber = ""
	defer func() { versionNumber = original }()

	if got := getPluginVersion(); got != "unknown" {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestGetCredentials_StaticTakesPriority(t *testing.T) {
	c := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Token:     "test-session-token",
	}
	creds := getCredentials(c)
	cp, err := creds.Get()
	if err != nil {
		t.Fatalf("unexpected error getting credentials: %s", err)
	}
	if cp.AccessKeyID != "test-access-key" {
		t.Fatalf("expected AccessKeyID=test-access-key, got %s", cp.AccessKeyID)
	}
	if cp.SecretAccessKey != "test-secret-key" {
		t.Fatalf("expected SecretAccessKey=test-secret-key, got %s", cp.SecretAccessKey)
	}
	if cp.SessionToken != "test-session-token" {
		t.Fatalf("expected SessionToken=test-session-token, got %s", cp.SessionToken)
	}
}

func TestGetCredentials_EnvFallback(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "env-access-key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "env-secret-key")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}()

	c := &Config{}
	creds := getCredentials(c)
	cp, err := creds.Get()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if cp.AccessKeyID != "env-access-key" {
		t.Fatalf("expected env-access-key, got %s", cp.AccessKeyID)
	}
}

func TestGetCredentials_ReturnsChain(t *testing.T) {
	c := &Config{}
	creds := getCredentials(c)
	if creds == nil {
		t.Fatal("expected non-nil credential chain")
	}
}
