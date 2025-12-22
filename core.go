package ng

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"

	nghttp "github.com/foxie-io/ng/http"
)

type EventStage string

const (
	StageAfterGuards       EventStage = "AFTER_GUARDS"
	StageAfterHandler      EventStage = "AFTER_HANDLER"
	StageAfterInterceptors EventStage = "AFTER_INTERCEPTORS"
	StageAfterMiddlewares  EventStage = "AFTER_MIDDLEWARES"
)

type ResponseHandler func(ctx context.Context, resp nghttp.HttpResponse) error

type (
	Core interface {
		Prefix() string
		Metadata(key any) (value any, found bool)
	}

	/*core
	1- Execute Guards → abort if error

	2- Run Middlewares in chain

	3- Execute Handler

	4- Run Interceptors after hooks
	*/
	core struct {
		prefix string

		// root execution
		preExecutes []PreHandler

		middlewares []Middleware

		guards []Guard

		interceptors []Interceptor

		// handlers
		handlers []Handler

		// metadata
		metadata sync.Map

		// built checker
		built atomic.Bool

		responseHandler ResponseHandler
	}
)

func newCore() *core {
	return &core{}
}

func (c *core) Metadata(key any) (value any, found bool) {
	return c.metadata.Load(key)
}

func (c *core) Prefix() string {
	return c.prefix
}

func (c *core) applyPreExecutes(ctx context.Context) error {
	for _, pre := range c.preExecutes {
		pre(ctx)
	}

	return nil
}

func (c *core) buildGuardChain(handler Handler) Handler {
	return func(ctx context.Context) (err error) {
		if len(c.guards) == 0 {
			return handler(ctx)
		}

		skipIds := getSkipperIds(ctx)
		hasSkipAllGuards := slices.Contains(skipIds, allGuard)

		if hasSkipAllGuards {
			return handler(ctx)
		}

		for _, guard := range c.guards {
			if canSkip(guard, skipIds) {
				continue
			}

			if err := guard.Allow(ctx); err != nil {
				return err
			}
		}

		return handler(ctx)
	}
}

func (c *core) buildMiddlewareChain(routeHandler Handler) Handler {
	next := routeHandler

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		middleware := c.middlewares[i]
		n := next

		next = func(ctx context.Context) (err error) {
			if canSkip(middleware, getSkipperIds(ctx)) { // runtime evaluation
				return n(ctx)
			}

			middleware.Use(ctx, n)
			return
		}
	}

	return next
}

func (c *core) buildInterceptorChain(routeHandler Handler) Handler {
	next := routeHandler

	for i := len(c.interceptors) - 1; i >= 0; i-- {
		interceptor := c.interceptors[i]
		n := next

		next = func(ctx context.Context) (err error) {
			if canSkip(interceptor, getSkipperIds(ctx)) { // runtime evaluation
				return n(ctx)
			}

			interceptor.Intercept(ctx, n)
			return
		}
	}

	return next
}

// WithPrefix sets a prefix for all routes in the core
func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.core.prefix = normolizePath(prefix)
	}
}

/*
can use IgnoreGuard to skip
*/
func WithGuards(guards ...Guard) Option {
	return func(c *config) {
		c.core.guards = append(c.core.guards, guards...)
	}
}

// WithMiddleware adds middlewares to the core
func WithMiddleware(middlewares ...Middleware) Option {
	return func(c *config) {
		c.core.middlewares = append(c.core.middlewares, middlewares...)
	}
}

// WithInterceptor adds interceptors to the core
func WithInterceptor(interceptors ...Interceptor) Option {
	return func(c *config) {
		c.core.interceptors = append(c.core.interceptors, interceptors...)
	}
}

// WithMetadata requires key-value pairs
func WithMetadata(pairs ...any) Option {
	if len(pairs)%2 != 0 {
		panic("WithMetadata requires key-value pairs")
	}

	return func(c *config) {
		for i := 0; i < len(pairs); i += 2 {
			k, v := pairs[i], pairs[i+1]
			c.core.metadata.Store(k, v)
		}
	}
}

// WithResponseHandler sets a custom response handler for the core
func WithResponseHandler(handler ResponseHandler) Option {
	return func(c *config) {
		c.core.responseHandler = handler
	}
}

// root level execution: before guard,middleware,etc...
func WithPreExecute(pre PreHandler) Option {
	return func(c *config) {
		c.core.preExecutes = append(c.core.preExecutes, pre)
	}
}
