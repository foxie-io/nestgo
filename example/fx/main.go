package main

import (
	"context"
	"example/fx/adapter"
	"example/fx/components/orders"
	"example/fx/components/users"
	"example/fx/router"

	"github.com/foxie-io/ng"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Starter struct {
	router *router.Router
	fiber  *fiber.App
	app    ng.App
}

func NewStarter(router *router.Router, app ng.App, fiber *fiber.App) *Starter {
	return &Starter{
		router: router,
		app:    app,
		fiber:  fiber,
	}
}

func (s *Starter) OnStart(ctx context.Context) error {
	// Application startup logic can be added here
	// Connect db, migrate, etc.

	s.router.Register(s.app, s.fiber)

	go s.fiber.Listen(":8080")

	return nil
}

func (s *Starter) OnStop(ctx context.Context) error {
	// Application shutdown logic can be added here
	// Close db, flush cache, etc.
	return nil
}

func NewApp() ng.App {
	return ng.NewApp(
		ng.WithResponseHandler(
			adapter.FiberResponseHandler,
		),
	)
}

func NewFiber() *fiber.App {
	return fiber.New(
		fiber.Config{
			EnablePrintRoutes: true,
		},
	)
}

func RunStarter(starter *Starter, lf fx.Lifecycle) {
	lf.Append(fx.StartStopHook(starter.OnStart, starter.OnStop))
}

func main() {
	fx.New(
		fx.Provide(
			NewApp,
			NewFiber,
			NewStarter,
			router.NewRouter,
		),

		users.Module,
		orders.Module,

		fx.Invoke(RunStarter),
	).Run()
}
