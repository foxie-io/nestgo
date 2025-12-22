package ngadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

func ServeMuxResponseHandler(ctx context.Context, info nghttp.HttpResponse) error {
	w := ng.MustLoad[http.ResponseWriter](ctx)

	if res, ok := info.(*nghttp.Response); ok {
		if res.Code == nghttp.CodeUnknown {
			raw, _ := res.GetMetadata("raw")
			res.Update(nghttp.Meta("error", fmt.Sprintf("%v", raw)))
		}
	}

	w.WriteHeader(info.StatusCode())
	bytes, _ := json.Marshal(info.Response())
	_, _ = w.Write(bytes)
	return nil
}

func ServeMuxHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := ng.NewContext(r.Context())

		// store in context
		ng.Store(ctx, w)
		ng.Store(ctx, r)

		// can extract from ctx if needed
		// w := ng.MustLoad[http.ResponseWriter](ctx)
		// r := ng.MustLoad[*http.Request](ctx)

		// invoke the handler
		scopeHandler()(ctx)
	}
}

func ServeMuxRegisterRoutes(ng ng.App, mux *http.ServeMux) {
	for _, route := range ng.Routes() {
		mux.HandleFunc(route.Path(), ServeMuxHandler(route.Handler))
	}
}
