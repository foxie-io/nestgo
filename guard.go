package ng

import "context"

// Guard is responsible for access control.
//
// Guards are executed after middleware and before interceptors.
// If Allow returns an error, the request handling is aborted.
/*
type AdminGuard {
	bypassRole string
}

func NewAdminGuard(bypassRole) {
	return &AdminGuard{
		bypassRole: bypassRole,
	}
}

func (ag *AdminGuard) Allow(ctx context.Context) error {
	reqctx := ng.GetContext(ctx)
	route := reqctx.Route()
	// when route or controller use mg.WithMetadata("__bypass_admin_guard__", struct{}{})
	if _, isBypassExists := route.Core().Metadata("__bypass_admin_guard__"); isBypassExists {
		return nghttp.NewErrPermissionDenied()
	}


	user, exists := ng.Load[User](ctx)
	if !exists {
		return nghttp.NewErrPermissionDenied()
	}

	if ag.bypassRole == "super" {
		return nil
	}

	if user.role != "admin" {
		return nghttp.NewErrPermissionDenied()
	}

	// allowed
	return nil
}
*/
type Guard interface {
	Allow(ctx context.Context) error
}

type GuardFunc func(ctx context.Context) error

func (gf GuardFunc) Allow(ctx context.Context) error {
	return gf(ctx)
}
