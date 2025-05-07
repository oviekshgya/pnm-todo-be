package middleware

import (
	"github.com/gofiber/fiber/v2"
	"sync"
	"time"
)

type rateLimitEntry struct {
	count     int
	timestamp time.Time
}

var (
	rateLimitMap = make(map[string]*rateLimitEntry)
	mu           sync.Mutex
)

func RateLimitMiddleware(limit int, period time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		mu.Lock()
		defer mu.Unlock()

		entry, exists := rateLimitMap[clientIP]
		now := time.Now()

		if !exists || now.Sub(entry.timestamp) > period {
			rateLimitMap[clientIP] = &rateLimitEntry{
				count:     1,
				timestamp: now,
			}
		} else {
			entry.count++
			if entry.count > limit {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"message": "Too many requests. Please try again later.",
				})
			}
		}

		return c.Next()
	}
}
