package adapter

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gin-gonic/gin"
)

func GinResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	ginctx := ng.MustLoad[*gin.Context](ctx)

	switch val := info.(type) {
	case *nghttp.RawResponse:
		ginctx.Writer.WriteHeader(val.StatusCode())
		ginctx.Writer.Write(val.Value())
		return nil
	case *nghttp.Response:
		ginctx.JSON(info.StatusCode(), info.Response())
		return nil
	case *nghttp.PanicError:
		fmt.Println("Unknown response type:", fmt.Sprintf("%T, value: %v", info, val.Value()), string(debug.Stack()))
		ginctx.JSON(info.StatusCode(), info.Response())
		return nil
	}

	fmt.Println("Unknown response type:", fmt.Sprintf("%T", info), string(debug.Stack()))
	ginctx.JSON(info.StatusCode(), info.Response())
	return nil
}

func GinHandler(scopeHandler func() ng.Handler) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		ctx, _ := ng.NewContext(gctx.Request.Context())

		// Store Gin context in NG context
		ng.Store(ctx, gctx)

		// Invoke the handler
		scopeHandler()(ctx)
	}
}

func GinRegisterRoutes(ngApp ng.App, router *gin.Engine) {
	for _, route := range ngApp.Routes() {
		router.Handle(route.Method(), route.Path(), GinHandler(route.Handler))
	}
}
