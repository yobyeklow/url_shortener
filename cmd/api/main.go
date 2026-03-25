package main

import (
	"url_shortener/internal/app"
	"url_shortener/internal/config"
)

func main() {
	cfg := config.NewConfig()
	app := app.NewApplication(cfg)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
