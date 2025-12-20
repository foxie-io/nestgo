package main

import (
	"context"
	"example/gin/adapter"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
	ng.DefaultControllerInitializer
}

func (c *HelloController) GetHello() ng.Route {
	return ng.NewRoute("GET", "/hello",
		ng.WithHandler(
			func(ctx context.Context) error {
				return ng.Respond(ctx, nghttp.NewResponse("Hello, Gin!"))
			},
		),
	)
}

// This is the entry point for the Gin-based example server.
// It demonstrates how to use the NG framework with the Gin adapter.
func main() {
	app := ng.NewApp(
		ng.WithResponseHandler(adapter.GinResponseHandler),
	)

	app.AddController(&HelloController{})

	app.Build()

	r := gin.Default()

	adapter.GinRegisterRoutes(app, r)

	r.Run(":8080")
	// curl http://localhost:8080/hello
	// out => {"code":"OK","data":"Hello, Gin!"}
}
