package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

var webhookSecret = os.Getenv("WEBHOOK_SECRET") // set in env

func main() {
	cfg := fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	app := fiber.New(cfg)

	app.Use(logger.New())
	app.Use(recover.New())

	app.Post("/webhook/issue", handle)

	addr := fmt.Sprintf(":3000")
	if err := app.Listen(addr); err != nil {
		panic(err.Error())
	}
}

func handle(c fiber.Ctx) error {
	body := c.Body()
	hs := c.GetReqHeaders()

	sigs, _ := hs["X-Hub-Signature-256"]
	if !validateSignature(body, sigs[0]) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "invalid signature",
		})
	}

	out, err := os.Create("out.json")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create out file",
		})
	}
	defer out.Close()
	_, err = out.Write(body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to write output",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"event": "X-GitHub-Event",
	})
}

func validateSignature(payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
