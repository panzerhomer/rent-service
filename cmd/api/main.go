package main

import (
	"avito/internal/app"
	"avito/internal/config"
	"log"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
