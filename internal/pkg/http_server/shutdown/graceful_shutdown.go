package httpserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NewContext() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return ctx
}
