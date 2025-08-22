package main

import (
	"log"

	"github.com/saurabhraut1212/ecommerce_backend/internal/config"
	"github.com/saurabhraut1212/ecommerce_backend/internal/db"
	"github.com/saurabhraut1212/ecommerce_backend/internal/router"
)

func main() {
	cfg := config.Load()

	client, err := db.New(cfg.MongoURI)

	if err != nil {
		log.Fatal(err)
	}
	app := router.New(cfg, client)
	log.Fatal(app.Listen(":" + cfg.Port))
}
