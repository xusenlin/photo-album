package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"photoAlbum/config"
	"photoAlbum/global"
	"photoAlbum/service"
)

func main() {
	var err error
	global.Config, err = config.New("./config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	global.PhotoAlbumList, global.PhotoAlbumMap, err = service.InitPhotoAlbum(global.Config.PhotoAlbumAbsolutePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(global.PhotoAlbumList)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err = app.Listen(":" + global.Config.ListenPort)

	if err != nil {
		fmt.Println(err)
		return
	}

}
