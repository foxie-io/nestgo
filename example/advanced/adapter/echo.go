package adapter

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"

	"github.com/labstack/echo/v4"
)

func ResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
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

func ToEchoHandler(scopeHandler func() ng.Handler) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx, rc := ng.NewContext(ectx.Request().Context())
		defer rc.Clear()
		ng.Store(ctx, ectx)
		return scopeHandler()(ctx)
	}
}

func RegisterRoutes(ng ng.App, echo *echo.Echo) {
	for _, route := range ng.Routes() {
		// fmt.Printf("Route: %s path=%s, %s\n", route.Method(), route.Path(), route.Name())
		echoHandler := ToEchoHandler(route.Handler)
		eroute := echo.Add(route.Method(), route.Path(), echoHandler)
		eroute.Name = route.Name()
	}

	for _, r := range echo.Routes() {
		fmt.Printf("Route: %s path=%s, name=%s\n", r.Method, r.Path, r.Name)
	}
}
