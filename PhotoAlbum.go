package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"photoAlbum/models"
)

func main() {

	path, err := filepath.Abs("./PhotoAlbum")

	if err != nil {
		fmt.Println(err)
		return
	}

	err = models.InitPhotoAlbum(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err = app.Listen(":3000")

	if err != nil {
		fmt.Println(err)
		return
	}

}
