package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/snykk/fiber-redis-shortener/backend/cache"
	"github.com/snykk/fiber-redis-shortener/backend/config"
	"github.com/snykk/fiber-redis-shortener/backend/handlers"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	engine := html.New("./frontend/html", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Cache as datasources
	redisCache := cache.NewRedisCache(config.AppConfig.REDISHost, config.AppConfig.REDISDbno, config.AppConfig.REDISPassword, time.Duration(config.AppConfig.REDISExpired))

	// Handler
	handler := handlers.NewHandler(redisCache)

	// CORS middleware
	app.Use(cors.New())

	// Routes
	app.Get("/", handlers.Root)
	app.Get("/shorten", handlers.Shorten)
	app.Post("/generate-shorten-url", handler.ShortenURL)
	app.Get("/:shortenURL", handler.ResolveURL)
	app.Static("/static", "./frontend")

	err := app.Listen(fmt.Sprintf(":%d", config.AppConfig.Port))
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}

}
