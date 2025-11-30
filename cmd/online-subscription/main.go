package main

import (
	"log"
	"net/http"
	_ "online-subscription/docs"
	"online-subscription/internal/app"
)

// @title Online Subscriptions API service
// @version 1.0
// @description агреграция данных об онлайн-подписках пользователей
// @BasePath /
func main() {
	app := app.Start()

	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
