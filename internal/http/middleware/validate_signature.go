package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"thingify/internal/http/response"

	"github.com/gofiber/fiber/v2"
)

const signatureHeader = "X-Hub-Signature-256"

// TODO: FOR DEBUG
var debug bool

func init() {
	if f, _ := strconv.ParseBool(os.Getenv("DEBUG")); f {
		slog.Info("DEBUG MODE ENABLED: Signature validation will be skipped")
		debug = true
	}
}

func ValidateHubSignature(secret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sigHeader := c.Get(signatureHeader, "")
		if sigHeader == "" {
			return response.BadRequest(c,
				fmt.Errorf("missing %s header", signatureHeader),
				"Missing signature header",
			)
		}

		body := c.Body()
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if !debug && !hmac.Equal([]byte(expected), []byte(sigHeader)) {
			return response.Unauthorized(c,
				fmt.Errorf("signature mismatch: expected %s, got %s", expected, sigHeader),
				"Invalid github signature",
			)
		}

		return c.Next()
	}
}
