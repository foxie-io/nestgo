package router

import (
	"example/fx/adapter"

	"github.com/foxie-io/ng"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Param struct {
	fx.In
	GlobalController []ng.ControllerInitializer `group:"global_controller_initializers"`
}

type Router struct {
	globalControllers []ng.ControllerInitializer
}

func NewRouter(p Param) *Router {
	return &Router{
		globalControllers: p.GlobalController,
	}
}

func (r *Router) Register(app ng.App, fiberApp *fiber.App) {
	app.AddController(r.globalControllers...)

	app.Build()

	adapter.FiberRegisterRoutes(app, fiberApp)
}
