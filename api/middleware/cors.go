package middleware

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"pnm-todo-be/pkg"
	"strings"
	"sync"
	"time"
)

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		//AllowOrigins: "http://localhost:5173", //Deployment
		//AllowOrigins:     "*", //Deployment
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY, GO-ACCESS-SIGNATURE, GO-TIMESTAMP",
		AllowCredentials: true,
	})
}

var basicAuthCache sync.Map

var apiKeyCache sync.Map

func BasicAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if strings.Contains(c.OriginalURL(), "/product") {
			return c.Next()
		}

		const BASIC_SCHEMA = "Basic "

		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, BASIC_SCHEMA) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized basic auth"})
		}

		base64Credentials := strings.TrimPrefix(authHeader, BASIC_SCHEMA)

		if ok, found := basicAuthCache.Load(fmt.Sprintf("%s%s", base64Credentials, c.IP())); found && ok.(bool) {
			return c.Next()
		}

		decoded, err := base64.StdEncoding.DecodeString(base64Credentials)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized decode"})
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized part"})
		}

		username := parts[0]
		password := parts[1]

		validUsername := viper.GetString("SERVICE_USERNAME")
		validPassword := viper.GetString("SERVICE_PASSWORD")

		if username == validUsername && password == validPassword {
			basicAuthCache.Store(fmt.Sprintf("%s%s", base64Credentials, c.IP()), true)
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized not valid"})
	}
}

func APIKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-KEY")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API Key is missing",
			})
		}

		if ok, found := apiKeyCache.Load(fmt.Sprintf("%s%s", apiKey, c.IP())); found && ok.(bool) {
			return c.Next()
		}

		if apiKey != viper.GetString("SERVICE_API_KEY") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid API Key",
			})
		}

		apiKeyCache.Store(fmt.Sprintf("%s%s", apiKey, c.IP()), true)

		return c.Next()
	}
}

var TokenCache = cache.New(10*time.Minute, 15*time.Minute)

func AuthBearer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized header null"})
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized token null"})
		}

		if _, found := TokenCache.Get(tokenString); found {
			return c.Next()
		}

		_, err := pkg.ExtractTokenJWT(tokenString, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
		}
		TokenCache.Set(tokenString, true, cache.DefaultExpiration)
		return c.Next()
	}
}
