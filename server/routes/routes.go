package routes

import (
	"basicCRUD/config"
	"basicCRUD/controllers"
	"basicCRUD/interfaces"

	"github.com/gofiber/fiber/v2"
)

func Serve(app *fiber.App) {
	// /health check
	{
		app.Get("/healthcheck", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})
	}

	// static
	{
		app.Static("/uploads", "./uploads")
	}

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// health check API V1
	{
		v1.Get("check", func(c *fiber.Ctx) error {
			return c.SendString("API V1 OK")
		})
	}

	/*
		type CreateUserForm struct {
			Name     string `json:"name" validate:"required"`
			Surname  string `json:"surname" validate:"required"`
			Password string `json:"password" validate:"required,min=4,max=15"`
		}
	*/

	// userGroup
	userGroup := v1.Group("/users")
	{
		var curl interfaces.CRUD = &controllers.UserController{DB: config.GetDB()}
		userGroup.Get("", curl.FindAll)
		userGroup.Get("/:id", curl.FindOne)
		userGroup.Post("", curl.Create)
		userGroup.Patch("/:id", curl.Update)
		userGroup.Delete("/:id", curl.Delete)
	}

	// roleGroup
	roleGroup := v1.Group("/roles")
	{
		var curl interfaces.CRUD = &controllers.RoleController{DB: config.GetDB()}
		roleGroup.Get("", curl.FindAll)
		roleGroup.Get("/:id", curl.FindOne)
		roleGroup.Post("", curl.Create)
		roleGroup.Patch("/:id", curl.Update)
		roleGroup.Delete("/:id", curl.Delete)
	}
}
