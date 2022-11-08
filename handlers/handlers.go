package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snykk/fiber-redis-shortener/cache"
	"github.com/snykk/fiber-redis-shortener/config"
	"github.com/snykk/fiber-redis-shortener/utils"
)

// dto
type request struct {
	LongURL string `json:"long_url" binding:"required"`
}

type Handler struct {
	redisCache cache.RedisCache
}

func NewHandler(redisCache cache.RedisCache) *Handler {
	return &Handler{
		redisCache: redisCache,
	}
}

func Root(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to fiber url shortener API",
	})
}

func (h Handler) ShortenURL(ctx *fiber.Ctx) error {
	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	shortenURL := utils.Randomize(10)
	if err := h.redisCache.Set(shortenURL, req.LongURL); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "shorten url created successfully",
		"shorten_url": config.AppConfig.Host + shortenURL,
	})

}

func (h Handler) ResolveURL(ctx *fiber.Ctx) error {
	shortenURL := ctx.Params("shortenURL")
	originURL, err := h.redisCache.Get(shortenURL)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Redirect(originURL, fiber.StatusTemporaryRedirect)
	return nil
}
