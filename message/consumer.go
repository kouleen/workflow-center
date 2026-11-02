package message

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kouleen/common/message"
)

func init() {
	go func() {
		if err := message.InitConsumers(message.ClientOptions{
			URL: os.Getenv("RABBITMQ_URL"),
		}); err != nil {
			log.Fatal(err)
		}
		defer func() { _ = message.CloseConsumers() }()
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		select {
		case <-ctx.Done():
			log.Printf("receive interrupt signal")
		}
	}()
}
