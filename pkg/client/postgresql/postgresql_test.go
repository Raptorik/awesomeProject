package postgresql

import (
	"awesomeProject/internal/config"
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 3

	sc := config.StorageConfig{
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		Username: "postgres",
		Password: "postgres",
	}

	t.Run("successful connection", func(t *testing.T) {
		client, err := NewClient(ctx, maxAttempts, sc)
		if err != nil {
			t.Fatalf("NewClient() failed: %v", err)
		}
		defer client.Close()
	})

	scInvalid := config.StorageConfig{
		Host:     "localhost",
		Port:     "5432",
		Database: "invalid_db",
		Username: "invalid_user",
		Password: "invalid_password",
	}

	t.Run("failed connection", func(t *testing.T) {
		client, err := NewClient(ctx, maxAttempts, scInvalid)
		if err == nil {
			if client != nil {
				client.Close()
			}
			t.Error("NewClient() should return an error when connections is failed")
		} else {
			t.Logf("Expected error: %v", err)
		}
	})
}
