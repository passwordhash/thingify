package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"thingify/internal/http/response"

	"github.com/gofiber/fiber/v2"
)

const signatureHeader = "X-Hub-Signature-256"

func ValidateHubSignature(secret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sigHeader := c.Get(signatureHeader)
		if sigHeader == "" {
			return response.BadRequest(c, nil, "Missing signature header")
		}

		body := c.Body()

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(expected), []byte(sigHeader)) {
			return response.Unauthorized(c, nil, "Invalid github signature")
		}

		return c.Next()
	}
}
