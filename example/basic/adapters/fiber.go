package adapters

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gofiber/fiber/v2"
)

func FiberResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	fctx := ng.MustLoad[*fiber.Ctx](ctx)

	switch val := info.(type) {
	case *nghttp.RawResponse:
		return fctx.Status(val.StatusCode()).SendString(string(val.Value()))
	case *nghttp.Response:
		return fctx.Status(val.StatusCode()).JSON(val.Response())
	case *nghttp.PanicError:
		fmt.Println("Unknown response type:", fmt.Sprintf("%T, value: %v", info, val.Value()), string(debug.Stack()))
		fctx.Status(info.StatusCode()).JSON(info.Response())
		return nil
	}

	fmt.Println("Unknown response type:", fmt.Sprintf("%T", info), string(debug.Stack()))
	return fctx.Status(info.StatusCode()).JSON(info.Response())
}

func FiberHandler(scopeHandler func() ng.Handler) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx, rc := ng.NewContext(fctx.Context())
		defer rc.Clear()

		// store fiber context
		ng.Store(ctx, fctx)

		ip := fctx.IP()
		ng.Store(ctx, ClientIp(ip))

		// get fiber context from ng ctx
		// fctx := ng.MustLoad[*fiber.Ctx](ctx)
		return scopeHandler()(ctx)
	}
}

func FiberRegisterRoutes(ng ng.App, app *fiber.App) {
	for _, route := range ng.Routes() {
		fiberHandler := FiberHandler(route.Handler)
		r := app.Add(route.Method(), route.Path(), fiberHandler)
		r.Name(route.Name())
	}
}
