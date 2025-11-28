package main

import (
	"log"
	"net/http"
	"online-subscription/internal/app"
)

func main() {
	app := app.Start()

	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
