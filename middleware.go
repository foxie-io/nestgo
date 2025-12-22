package ng

import "context"

// Middleware executes before guards and interceptors.
//
// It is typically used for cross-cutting concerns such as, due to its position
//   - Logging
//   - Request mutation
//   - Tracing
//   - Attaching values to context
//
// Middleware must call next to continue the request flow.
/*
type TokenParser {
}

func UseTokenParser() {
	return &TokenParser{}
}

func (ag *TokenParser) Use(ctx context.Context, next Handler) {
	echoCtx := ng.MustLoad[echo.Context](ctx)
	token := getToken(echoCtx)

	if token := "" {
		user := getUser(token)
		ng.Store[User](ctx,user)
	}

	// parse token
	next(ctx)
}
*/
type Middleware interface {
	Use(ctx context.Context, next Handler)
}

type MiddlewareFunc func(ctx context.Context, next Handler)

func (mf MiddlewareFunc) Use(ctx context.Context, next Handler) {
	mf(ctx, next)
}
