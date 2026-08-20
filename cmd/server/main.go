package main

import (
	"context"
	"github.com/wyw14/cry052/internal/application"
	"github.com/wyw14/cry052/internal/config"
	"github.com/wyw14/cry052/internal/repository/memory"
	httptransport "github.com/wyw14/cry052/internal/transport/http"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	log, _ := zap.NewProduction()
	defer log.Sync()
	svc := application.New(memory.New(), application.RealClock{}, &application.IDs{})
	r := httptransport.NewRouter(svc)
	go r.Run(cfg.HTTPAddr)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
