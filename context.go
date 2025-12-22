package ng

import (
	"context"
	"errors"
	"fmt"
	"sync"

	nghttp "github.com/foxie-io/ng/http"
)

const (
	responseKey PayloadKey = "response"
	routeKey    PayloadKey = "route"
)

// Context is a request context
type Context interface {
	// locals store

	// Store store value into context with given key
	Store(key PayloadKeyer, value any)

	// Load load value from context by given key
	Load(key PayloadKeyer) (value any, ok bool)

	// LoadOrStore load value from context by given key,
	// if not found, store the value into context
	LoadOrStore(key PayloadKeyer, value any) (actual any, loaded bool)

	// Delete delete value from context by given key
	Delete(key PayloadKeyer)

	// Clear clear all info stored in context
	Clear()

	// SetResponse
	SetResponse(resp nghttp.HttpResponse) Context

	// GetResponse
	GetResponse() nghttp.HttpResponse

	// response to endpoint
	Response() error

	// immutable data
	setOwner(route Route) Context

	// not available pre execute
	Route() RouteData

	// clone context for goroutine use
	Clone() Context
}

type RouteData interface {
	Core() Core
	Name() string
	Method() string
	Path() string
}

var _ Context = (*requestContext)(nil)

type requestContext struct {
	locals sync.Map
}

func newRequestContext() *requestContext {
	r := &requestContext{}
	return r
}

// Store store value into context with given key
func (r *requestContext) Store(key PayloadKeyer, value any) {
	r.locals.Store(key.PayloadKey(), value)
}

// Load load value from context by given key
func (r *requestContext) Load(key PayloadKeyer) (value any, ok bool) {
	return r.locals.Load(key.PayloadKey())
}

// Delete delete value from context by given key
func (r *requestContext) Delete(key PayloadKeyer) {
	r.locals.Delete(key.PayloadKey())
}

// LoadOrStore load value from context by given key,
// if not found, store the value into context
func (r *requestContext) LoadOrStore(key PayloadKeyer, value any) (actual any, loaded bool) {
	return r.locals.LoadOrStore(key.PayloadKey(), value)
}

// Clear clear all info stored in context
func (r *requestContext) Clear() {
	r.locals.Clear()
}

// SetResponse set request response
func (r *requestContext) SetResponse(resp nghttp.HttpResponse) Context {
	r.Store(responseKey, resp)
	return r
}

// Response get request response
func (r *requestContext) GetResponse() nghttp.HttpResponse {
	resp, ok := r.Load(responseKey)
	if ok {
		return resp.(nghttp.HttpResponse)
	}
	return nil
}

func (r *requestContext) Response() error {
	resp := r.GetResponse()
	if resp == nil {
		return errors.New("response not found")
	}

	ThrowResponse(resp)
	return nil
}

func (r *requestContext) GetRoute() Route {
	resp, ok := r.Load(routeKey)
	if ok {
		return resp.(Route)
	}
	return nil
}

func (r *requestContext) setOwner(route Route) Context {
	r.Store(routeKey, route)
	return r
}

func (r *requestContext) Route() RouteData {
	resp, ok := r.Load(routeKey)
	if ok {
		return resp.(Route)
	}
	return nil
}

// Clone create a clone of request context to use in goroutine after request end
func (r *requestContext) Clone() Context {
	clone := &requestContext{}
	r.locals.Range(func(key, value any) bool {
		clone.locals.Store(key, value)
		return true
	})
	return clone
}

func dynamicKey[T any](keys ...PayloadKeyer) PayloadKeyer {
	if len(keys) == 0 {
		return TypeKey[T]{}
	}
	return keys[0]
}

// Store store value into context with given key
func Store[T any](ctx context.Context, value T, keys ...PayloadKeyer) {
	key := dynamicKey[T](keys...)
	GetContext(ctx).Store(key, value)
}

// Load load value from context by given key
func Load[T any](ctx context.Context, keys ...PayloadKeyer) (value T, err error) {
	key := dynamicKey[T](keys...)
	val, loaded := GetContext(ctx).Load(key)
	if !loaded {
		var zero T
		return zero, fmt.Errorf("not found key: %s", key.PayloadKey())
	}

	expectedType, ok := val.(T)
	if !ok {
		return expectedType, fmt.Errorf("invalid type, expected %T, got %T", val, expectedType)
	}

	return expectedType, nil
}

func Delete[T any](ctx context.Context, keys ...PayloadKeyer) {
	key := dynamicKey[T](keys...)
	GetContext(ctx).Delete(key)
}

// LoadOrStore load value from context by given key,
// if not found, store the value into context
func LoadOrStore[T any](ctx context.Context, value T, keys ...PayloadKeyer) (actual T, loaded bool, err error) {
	key := dynamicKey[T](keys...)
	val, loaded := GetContext(ctx).LoadOrStore(key, value)
	expectedType, ok := val.(T)
	if !ok {
		return expectedType, loaded, fmt.Errorf("invalid type, expected %T, got %T", val, expectedType)
	}
	return expectedType, loaded, nil
}

// MustLoad load value from context by given key,
// panic if not found
func MustLoad[T any](ctx context.Context, keys ...PayloadKeyer) T {
	val, err := Load[T](ctx, keys...)
	if err != nil {
		panic(err)
	}
	return val
}

// MustLoadOrStore load value from context by given key,
// if not found, store the value into context, panic if not found
func MustLoadOrStore[T any](ctx context.Context, value T, keys ...PayloadKeyer) (val T, loaded bool) {
	val, loaded, err := LoadOrStore(ctx, value, keys...)
	if err != nil {
		panic(err)
	}
	return val, loaded
}

func withContext(ctx context.Context, rctx Context) context.Context {
	return context.WithValue(ctx, TypeKey[Context]{}, rctx)
}

// if return existed ctx, otherwise create new one
func NewContext(ctx context.Context) (context.Context, Context) {
	rc := GetContext(ctx)
	if rc != nil {
		return ctx, rc
	}

	rc = newRequestContext()
	return withContext(ctx, rc), rc
}

func acquireContext(ctx context.Context) (c context.Context, rc Context, new bool) {
	rc = GetContext(ctx)
	if rc != nil {
		return ctx, rc, false
	}

	rc = newRequestContext()
	return withContext(ctx, rc), rc, true
}

// GetContext get context from given context
func GetContext(ctx context.Context) Context {
	c, _ := ctx.Value(TypeKey[Context]{}).(Context)
	return c
}
