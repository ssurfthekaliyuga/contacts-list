package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func extractContactID(c *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(c.Params("contactID"))
}
