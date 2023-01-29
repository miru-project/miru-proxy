package main

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/imroc/req/v3"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	client := req.C().SetTimeout(5 * time.Second)
	app.Group("/request/*", func(ctx *fiber.Ctx) error {
		method := ctx.Method()
		header := ctx.GetReqHeaders()
		url := header["Miru-Url"]
		if ctx.Query("miru-iframe") != "" {
			b, err := base64.StdEncoding.DecodeString(ctx.Query("miru-iframe"))
			if err != nil {
				return fiber.ErrBadGateway
			}
			url = string(b)
			ctx.Context().SetContentType("text/html; charset=utf-8")
		}
		ua := header["Miru-Ua"]
		if ua == "" {
			ua = header["User-Agent"]
		}
		referer := header["Miru-Referer"]
		if referer == "" {
			referer = header["Referer"]
		}
		url = url + strings.Replace(ctx.OriginalURL(), "/request", "", -1)
		proxy := client.Get(url)
		if method == "POST" {
			proxy = client.Post(url)
		}
		p := proxy.SetHeaders(map[string]string{
			"Referer":      referer,
			"User-Agent":   ua,
			"Content-Type": header["Content-Type"],
		}).SetBody(ctx.Body()).Do()
		if p.Err != nil {
			return fiber.ErrBadGateway
		}
		println(url)
		return ctx.Status(p.StatusCode).SendString(p.String())
	})
	app.Get("/", func(ctx *fiber.Ctx) error {
		println(ctx.Request().Header.String())
		return ctx.JSON(map[string]string{
			"status":  "ok",
			"version": "v1.0.1",
			"miru":    "見る",
		})
	})
	app.Listen(":8080")
}
