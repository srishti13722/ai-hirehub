package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func RoleBasedMiddleware(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get the JWT token from the request header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		// parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// extract role from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["role"] != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: You do not have permission"})
		}

		return c.Next()
	}
}
