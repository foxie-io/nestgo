package adapter

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/labstack/echo/v4"
)

func EchoResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	ectx := ng.MustLoad[echo.Context](ctx)

	switch val := info.(type) {
	case *nghttp.RawResponse:
		return ectx.String(val.StatusCode(), string(val.Value()))

	case *nghttp.Response:
		return ectx.JSON(val.StatusCode(), val.Response())

	case *nghttp.PanicError:
		fmt.Println("Unknown response type:", fmt.Sprintf("%T, value: %v", info, val.Value()), string(debug.Stack()))
		return ectx.JSON(val.StatusCode(), val.Response())
	}

	fmt.Println("Unknown response type:", fmt.Sprintf("%T", info), string(debug.Stack()))
	resp := nghttp.NewErrUnknown()
	return ectx.JSON(resp.StatusCode(), resp.Response())
}

func EchoHandler(scopeHandler func() ng.Handler) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		ctx, rc := ng.NewContext(echoCtx.Request().Context())
		defer rc.Clear()

		// store echo context
		ng.Store(ctx, echoCtx)

		// get echo context from ng ctx
		// echoCtx := ng.MustLoad[echo.Context](ctx)
		return scopeHandler()(ctx)
	}
}

func EchoRegisterRoutes(ng ng.App, echo *echo.Echo) {
	for _, route := range ng.Routes() {
		echoHandler := EchoHandler(route.Handler)
		eroute := echo.Add(route.Method(), route.Path(), echoHandler)
		eroute.Name = route.Name()
	}
}
