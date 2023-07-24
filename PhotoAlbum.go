package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"photoAlbum/models"
)

func main() {

	err := models.InitPhotoAlbum("")
	if err != nil {
		fmt.Println(err)
		return
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
