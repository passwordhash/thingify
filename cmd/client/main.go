package main

import (
	"context"
	"os"
	"os/signal"
	"thingify/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.MustLoadClient()

	<-ctx.Done()
}
