package main

import (
	"context"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gofiber/fiber/v2"
)

type HelloController struct {
	ng.DefaultControllerInitializer
}

func (c *HelloController) GetHello() ng.Route {
	return ng.NewRoute("GET", "/hello",
		ng.WithHandler(
			func(ctx context.Context) error {
				return ng.Respond(ctx, nghttp.NewResponse("Hello, Fiber!"))
			},
		),
	)
}

func main() {
	app := ng.NewApp(
		ng.WithResponseHandler(func(ctx context.Context, info *ng.ResponseInfo) error {
			fctx := ng.MustLoad[*fiber.Ctx](ctx)
			if info.HttpResponse != nil {
				return fctx.Status(info.HttpResponse.StatusCode()).JSON(info.HttpResponse.Response())
			}
			return fctx.Status(500).SendString("Internal Server Error")
		}),
	)

	app.AddController(&HelloController{})

	app.Build()

	fiberApp := fiber.New()
	for _, route := range app.Routes() {
		fiberApp.Add(route.Method(), route.Path(), func(c *fiber.Ctx) error {
			ngCtx := ng.NewContext()
			ctx := ng.WithContext(c.Context(), ngCtx)
			ng.Store(ctx, c)
			return route.Handler()(ctx)
		})
	}

	fiberApp.Listen(":8080")
}
