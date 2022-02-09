package main

import (
	"basicCRUD/config"
	"basicCRUD/infrastructure"
	"basicCRUD/migrations"
	"basicCRUD/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatal(err)
	}

	// Init App such as connect DB etc.
	config.AppInit()
	defer config.CloseDB()

	migrations.Migrate()

	app := fiber.New()

	uploadDirs := [...]string{"users"}
	for _, dir := range uploadDirs {
		// permission-calculator.org
		os.MkdirAll("uploads/"+dir, 0755)
	}
	routes.Serve(app)

	app.Listen(":" + infrastructure.Get().AppPort)
}
