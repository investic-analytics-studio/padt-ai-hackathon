package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

func AuthCRMMiddleware(userRepo port.AuthRepo, jwtService port.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawToken := c.Get("Authorization")
		if rawToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		idToken := extractBearerToken(rawToken)
		if idToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		token, err := jwtService.VerifyCRMToken(idToken)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		uid, err := json.Marshal(token.UID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = userRepo.CheckCRMUserByUid(token.UID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}

		uidString := strings.ReplaceAll(string(uid), "\"", "")
		userInfo, err := userRepo.GetCRMUserInfo(uidString)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user info"})
		}
		c.Locals("uid", uidString)
		c.Locals("username", userInfo.Username)
		return c.Next()
	}
}

func AuthMiddleware(userRepo port.AuthRepo, jwtService port.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawToken := c.Get("Authorization")
		if rawToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		idToken := extractBearerToken(rawToken)
		if idToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		token, err := jwtService.VerifyToken(idToken)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		uid, err := json.Marshal(token.UID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = userRepo.CheckUserByUid(string(uid))
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}

		uidString := strings.ReplaceAll(string(uid), "\"", "")

		userInfo, err := userRepo.GetUserInfo(uidString)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Locals("uid", uidString)
		c.Locals("email", userInfo.Email)
		c.Locals("twitter_uid", userInfo.TwitterUID)
		c.Locals("twitter_name", userInfo.TwitterName)
		return c.Next()
	}
}

func extractBearerToken(rawToken string) string {
	return strings.TrimSpace(strings.Replace(rawToken, "Bearer", "", 1))
}
