package endpoints

import (
	"github.com/gofiber/fiber/v2"
)

type ContactsControllers interface {
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RegisterContacts(ctrl ContactsControllers, router fiber.Router) {
	router = router.Group("/contacts")

	router.Get("/", ctrl.Get)
	router.Post("/", ctrl.Create)
	router.Patch("/:contactID", ctrl.Update)
	router.Delete("/:contactID", ctrl.Delete)
}
