package client

import (
	"context"
	"log"
	"os/signal"
	"protocol/internal/server"
	"testing"
	"time"
)

const (
	port = ":8080"
)

func TestGetSend(t *testing.T) {
	server := server.New(port)

	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()

	go func() {
		if err := server.Start(ctx); err != nil {
			log.Printf("Ошибка при работе сервера: %v, остановка...", err)
		}
	}()
	time.Sleep(500 * time.Millisecond) // Задержка, чтобы сервер успел запуститься
	// Создаем клиента
	client, err := NewClient(port)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	key := []byte("test_key")
	value := []byte("test_value")

	err = client.Send(key, value)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}

	// Проверяем метод Get
	value, err = client.Get([]byte("test_key"))
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	expectedValue := []byte("test_value")
	if string(value) != string(expectedValue) {
		t.Errorf("Expected value %s, got %s", expectedValue, value)
	}

}
