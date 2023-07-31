package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"photoAlbum/config"
	"photoAlbum/global"
	"photoAlbum/service"
	"strconv"
)

func main() {
	var err error
	global.Config, err = config.New("./config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	global.PhotoAlbumList, err = service.InitPhotoAlbum(global.Config.PhotoAlbumAbsolutePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/public", "./public")

	app.Get("/:pageNum?/:pageSize?", func(c *fiber.Ctx) error {

		pageNum, err := strconv.Atoi(c.Params("pageNum"))
		if err != nil {
			pageNum = 1
		}
		pageSize, err := strconv.Atoi(c.Params("pageSize"))
		if err != nil {
			pageSize = 2
		}

		photoAlbumList, p := global.PhotoAlbumList.Pagination(pageNum, pageSize)
		return c.Render("index", fiber.Map{
			"Config":         global.Config,
			"PhotoAlbumList": photoAlbumList,
			"Pagination":     p,
		})
	})

	err = app.Listen(":" + global.Config.ListenPort)

	if err != nil {
		fmt.Println(err)
		return
	}

}
