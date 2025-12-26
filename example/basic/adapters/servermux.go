package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

func ServeMuxResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	w := ng.MustLoad[http.ResponseWriter](ctx)
	w.WriteHeader(info.StatusCode())

	switch val := info.(type) {
	case *nghttp.RawResponse:
		_, err := w.Write(val.Value())
		return err

	case *nghttp.Response:
		jsonbytes, _ := json.Marshal(info.Response())
		_, err := w.Write(jsonbytes)
		return err

	case *nghttp.PanicError:
		fmt.Println("Unknown response type:", fmt.Sprintf("%T, value: %v", info, val.Value()), string(debug.Stack()))
		jsonbytes, _ := json.Marshal(info.Response())
		_, err := w.Write(jsonbytes)
		return err
	}

	fmt.Println("Unknown response type:", fmt.Sprintf("%T, info", info), string(debug.Stack()))
	jsonbytes, _ := json.Marshal(info.Response())
	_, _ = w.Write(jsonbytes)
	return nil
}

func ServeMuxHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, rc := ng.NewContext(r.Context())
		defer rc.Clear()

		// store http.ResponseWriter in context
		ng.Store(ctx, w)

		ip := r.RemoteAddr
		ng.Store(ctx, ClientIp(ip))

		// invoke the handler
		scopeHandler()(ctx)
	}
}

func ServeMuxRegisterRoutes(ng ng.App, mux *http.ServeMux) {
	for _, route := range ng.Routes() {
		// GET /path format
		muxPath := fmt.Sprintf("%s %s", route.Method(), route.Path())
		mux.HandleFunc(muxPath, ServeMuxHandler(route.Handler))
	}
}
