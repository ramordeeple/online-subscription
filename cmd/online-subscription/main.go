package main

import (
	"context"
	"net/http"
	_ "online-subscription/docs"
	"online-subscription/internal/app"
	"online-subscription/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title Online Subscriptions API service
// @version 1.0
// @description Агреграция данных об онлайн-подписках пользователей
// @BasePath /
func main() {
	application := app.Start()
	srv := application.Server

	go func() {
		logger.Info("Starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server crashed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", zap.Error(err))
	} else {
		logger.Info("Server stopped gracefully")
	}
}
