package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/snykk/fiber-redis-shortener/cache"
	"github.com/snykk/fiber-redis-shortener/config"
	"github.com/snykk/fiber-redis-shortener/handlers"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New()

	// Cache as datasources
	redisCache := cache.NewRedisCache(config.AppConfig.REDISHost, config.AppConfig.REDISDbno, config.AppConfig.REDISPassword, time.Duration(config.AppConfig.REDISExpired))

	// Handler
	handler := handlers.NewHandler(redisCache)

	// Routes
	app.Get("/", handlers.Root)
	app.Post("/generate-shorten-url", handler.ShortenURL)
	app.Get("/:shortenURL", handler.ResolveURL)

	err := app.Listen(fmt.Sprintf(":%d", config.AppConfig.Port))
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}

}
