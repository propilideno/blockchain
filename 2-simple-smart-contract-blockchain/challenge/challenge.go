package main

import (
	"github.com/gofiber/fiber/v2"
)

var dataStore = make(map[string]string)

func main() {
	app := fiber.New()

	// Endpoint to submit data

	app.Post("/acme", func(c *fiber.Ctx) error {
		var data struct {
			Wallet string `json:"wallet"`
			Digest string `json:"digest"`
		}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		dataStore[data.Wallet] = data.Digest
		return c.Status(fiber.StatusOK).SendString("Data received")
	})


	// Endpoint to get data
	app.Get("/acme", func(c *fiber.Ctx) error {
		wallet := c.Query("wallet")
		if wallet == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing wallet ID")
		}

		value, exists := dataStore[wallet]
		if !exists {
			return c.Status(fiber.StatusNotFound).SendString("No data found")
		}

		return c.Status(fiber.StatusOK).SendString(value)
	})

	app.Listen(":3000")
}
