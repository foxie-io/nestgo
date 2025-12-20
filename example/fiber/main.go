package main

import (
	"context"
	"example/fiber/adapter"

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
		ng.WithResponseHandler(adapter.FiberResponseHandler),
	)

	app.AddController(&HelloController{})

	app.Build()

	fiberApp := fiber.New()

	adapter.FiberRegisterRoutes(app, fiberApp)

	fiberApp.Listen(":8080")
}
