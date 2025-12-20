package main

import (
	"context"
	"example/chi/adapter"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/go-chi/chi/v5"
)

type HelloController struct {
	ng.DefaultControllerInitializer
}

func (c *HelloController) GetHello() ng.Route {
	return ng.NewRoute("GET", "/hello",
		ng.WithHandler(
			func(ctx context.Context) error {
				return ng.Respond(ctx, nghttp.NewResponse("Hello, Chi!"))
			},
		),
	)
}

func main() {
	app := ng.NewApp(
		ng.WithResponseHandler(adapter.ChiResponseHandler),
	)

	app.AddController(&HelloController{})

	app.Build()

	r := chi.NewRouter()

	adapter.ChiRegisterRoutes(app, r)

	http.ListenAndServe(":8080", r)
}
