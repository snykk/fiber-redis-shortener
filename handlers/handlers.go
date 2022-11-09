package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snykk/fiber-redis-shortener/cache"
	"github.com/snykk/fiber-redis-shortener/config"
	"github.com/snykk/fiber-redis-shortener/utils"
)

// dto
type request struct {
	LongURL   string `json:"long_url" binding:"required"`
	CustomURL string `json:"custom_url"`
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
		"message":     "Welcome to fiber url shortener API",
		"maintainer":  "Moh. Najib Fikri aka snykk",
		"repository":  "https://github.com/snykk/fiber-redis-shortener",
		"another api": "https://golib-backend.herokuapp.com/",
	})
}

func (h Handler) ShortenURL(ctx *fiber.Ctx) error {
	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"error":  err.Error(),
		})
	}

	var shortenURL string
	if req.CustomURL == "" {
		shortenURL = utils.SecureRandomString(10)
		for redisURL, _ := h.redisCache.Get(shortenURL); redisURL != ""; {
			shortenURL = utils.SecureRandomString(10)
		}
	} else {
		redisURL, _ := h.redisCache.Get(req.CustomURL)
		if redisURL != "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "shorten url already exists",
			})
		}

		shortenURL = req.CustomURL
	}

	// set shorten url to redis
	if err := h.redisCache.Set(shortenURL, req.LongURL); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"error":  err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      true,
		"message":     "shorten url created successfully",
		"shorten_url": config.AppConfig.Host + shortenURL,
	})

}

func (h Handler) ResolveURL(ctx *fiber.Ctx) error {
	shortenURL := ctx.Params("shortenURL")
	originURL, err := h.redisCache.Get(shortenURL)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"error":  err.Error(),
		})
	}

	ctx.Redirect(originURL, fiber.StatusTemporaryRedirect)
	return nil
}
