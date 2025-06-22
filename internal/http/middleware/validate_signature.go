package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
)

const signatureHeader = "X-Hub-Signature-256"

func ValidateHubSignature(secret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sigHeader := c.Get(signatureHeader)
		if sigHeader == "" {
			// TODO: move to response package
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing signature"})
		}

		body := c.Body()

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(expected), []byte(sigHeader)) {
			// TODO: move to response package
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid signature"})
		}

		return c.Next()
	}
}
