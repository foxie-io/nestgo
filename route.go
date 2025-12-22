package ng

import (
	"context"

	nghttp "github.com/foxie-io/ng/http"
)

var _ Route = (*route)(nil)

type (
	Route interface {
		Core() Core
		Name() string
		Method() string
		Path() string
		Handler() Handler
	}

	route struct {
		core    *core
		name    string
		method  string
		path    string
		handler Handler
	}
)

func (r *route) Core() Core     { return r.core }
func (r *route) Name() string   { return r.name }
func (r *route) Method() string { return r.method }
func (r *route) Path() string   { return r.path }
func (r *route) Handler() Handler {
	if r.handler == nil {
		panic("route has not built yet")
	}
	return r.handler
}

type HandlerOption = Option

// NewRoute create new route instance
func NewRoute(method string, path string, opt HandlerOption, opts ...Option) Route {
	route := &route{
		method: method,
		path:   normolizePath(path),
		core:   newCore(),
	}

	config := newConfig()
	config.bindRoute(route)
	config.bindCore(route.core)
	config.update(opt, Opitons(opts...))
	return route
}

func (r *route) addPreCore(preCores ...Core) Route {
	var (
		responseHander ResponseHandler
		preExcutes     = []PreHandler{}
		guards         = []Guard{}
		middlewares    = []Middleware{}
		prefix         string
	)

	allCores := []Core{}
	allCores = append(allCores, preCores...)
	allCores = append(allCores, r.core)

	for _, c := range allCores {
		core := c.(*core)
		prefix += c.Prefix()

		guards = append(guards, core.guards...)
		middlewares = append(middlewares, core.middlewares...)
		preExcutes = append(preExcutes, core.preExecutes...)

		// merge metadata
		core.metadata.Range(func(key, value any) bool {
			r.core.metadata.Store(key, value)
			return true
		})

		if core.responseHandler != nil {
			responseHander = core.responseHandler
		}
	}

	// final route info
	r.path = prefix + r.path
	r.core.responseHandler = responseHander
	r.core.preExecutes = preExcutes
	r.core.guards = guards
	r.core.middlewares = middlewares
	return r
}

func (r *route) build() {
	if r.core.built.Load() {
		panic("core already built")
	}

	r.handler = r.buildRequestFlow()
	r.core.built.Store(true)
}

func (r *route) buildResponseHandler() ResponseHandler {
	responseHandler := r.core.responseHandler
	if responseHandler == nil {
		panic("response handler is not defined, WithResponseHandler is required")
	}

	return responseHandler
}

func (r *route) buildHandler() Handler {
	handler := Handle(r.core.handlers...)
	return handler
}

func (r *route) withSavedResponseState(next Handler) Handler {
	return func(ctx context.Context) error {
		defer func() {
			if r := recover(); r != nil {
				setResponseAny(GetContext(ctx), r)
			}
		}()

		if err := next(ctx); err != nil {
			panic(err)
		}

		return nil
	}
}

// prexecute -> middleware -> guard -> interceptor -> route handler -> response handler
func (r *route) buildRequestFlow() Handler {

	// route handler with response capture
	routeHandler := r.withSavedResponseState(r.buildHandler())

	// interceptor around route handler
	interceptorChain := r.withSavedResponseState(r.core.buildInterceptorChain(routeHandler))

	// guard before interceptor
	guardChain := r.withSavedResponseState(r.core.buildGuardChain(interceptorChain))

	// middleware around guard
	middlewareChain := r.core.buildMiddlewareChain(guardChain)

	// last execution: response handler
	final := r.buildResponseHandler()

	return func(ctx context.Context) (err error) {
		ctx, rc, created := acquireContext(ctx)
		if created {
			defer rc.Clear()
		}

		defer func() {
			// final response handling
			httpResp := nghttp.WrapResponse(rc.GetResponse())
			err = final(ctx, httpResp)
		}()

		defer func() {
			// 3 capture panic to response
			if r := recover(); r != nil {
				setResponseAny(rc, r)
			}
		}()

		// 1 pre executes before everything
		if err := r.core.applyPreExecutes(ctx); err != nil {
			return err
		}

		// 2 middleware -> guard -> interceptor -> route handler
		if err := middlewareChain(ctx); err != nil {
			panic(err)
		}

		return nil
	}
}
