package interfaces

import "github.com/gofiber/fiber/v2"

type CRUD interface {
	Create(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}
