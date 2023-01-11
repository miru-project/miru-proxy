package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/imroc/req/v3"
	"time"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	client := req.C().SetTimeout(5 * time.Second)
	app.Group("/request", func(ctx *fiber.Ctx) error {
		method := ctx.Method()
		header := ctx.GetReqHeaders()
		url := header["Miru-Url"]
		ua := header["Miru-Ua"]
		if ua == "" {
			ua = header["User-Agent"]
		}
		referer := header["Miru-Referer"]
		if referer == "" {
			referer = header["Referer"]
		}

		proxy := client.Get(url)
		if method == "POST" {
			proxy = client.Post(url)
		}
		println(ctx.Request().Header.String())
		p := proxy.SetHeaders(map[string]string{
			"Referer":    referer,
			"User-Agent": ua,
		}).Do()
		if p.Err != nil {
			return fiber.ErrBadGateway
		}

		return ctx.Status(p.StatusCode).SendString(p.String())
	})
	app.Get("/", func(ctx *fiber.Ctx) error {
		println(ctx.Request().Header.String())
		return ctx.JSON(map[string]string{
			"status":  "ok",
			"version": "v0.0.1",
			"miru":    "見る",
		})
	})
	app.Listen(":8080")
}
