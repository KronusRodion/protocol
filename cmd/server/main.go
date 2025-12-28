package main

import (
	"context"
	"log"
	"os/signal"

	"github.com/KronusRodion/protocol/internal/server"
)

const (
	port = ":8080"
)

func main() {
	server := server.New(port)

	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()

	if err := server.Start(ctx); err != nil {
		log.Printf("Ошибка при работе сервера: %v, остановка...", err)
	}

}
